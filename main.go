package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/alecthomas/kong"
	"github.com/atotto/clipboard"
	"github.com/briandowns/spinner"
	"github.com/fatih/color"
	"github.com/gofireflyio/aiac/v4/libaiac"
	"github.com/gofireflyio/aiac/v4/libaiac/bedrock"
	"github.com/gofireflyio/aiac/v4/libaiac/openai"
	"github.com/gofireflyio/aiac/v4/libaiac/types"
	"github.com/manifoldco/promptui"
	"github.com/rodaine/table"
)

type flags struct {
	Backend    libaiac.BackendName `help:"Backend to use (openai, bedrock)" enum:"openai,bedrock" default:"openai" short:"b" env:"AIAC_BACKEND"`
	ListModels struct {
		Type types.ModelType `arg:"" help:"List models of specific type" optional:""`
	} `cmd:"" help:"List supported models"`
	Get struct {
		// OpenAI flags
		APIKey     string `help:"OpenAI API key" env:"OPENAI_API_KEY"`
		URL        string `help:"OpenAI API url. Can be Azure Open AI service" default:"https://api.openai.com/v1" env:"OPENAI_API_URL"`
		APIVersion string `help:"OpenAI API version" default:"" env:"OPENAI_API_VERSION"`

		// Amazon Bedrock flags
		AWSProfile string `help:"AWS profile" default:"default" env:"AWS_PROFILE"`
		AWSRegion  string `help:"AWS region" default:"us-east-1" env:"AWS_REGION"`

		// Generic Flags
		OutputFile string   `help:"Output file to push resulting code to" optional:"" type:"path" short:"o"`         //nolint: lll
		ReadmeFile string   `help:"Readme file to push entire Markdown output to" optional:"" type:"path" short:"r"` //nolint: lll
		Quiet      bool     `help:"Non-interactive mode, print/save output and exit" default:"false" short:"q"`      //nolint: lll
		Full       bool     `help:"Print full Markdown output to stdout" default:"false" short:"f"`                  //nolint: lll
		Model      string   `help:"Model to use, default to \"gpt-3.5-turbo\""`
		What       []string `arg:"" help:"Which IaC template to generate"`
		Clipboard  bool     `help:"Copy generated code to clipboard (in --quiet mode)"`
	} `cmd:"" help:"Generate IaC code" aliases:"generate"`
	Version struct{} `cmd:"" help:"Print aiac version and exit"`
}

func main() {
	if len(os.Args) < 2 { //nolint: gomnd
		os.Args = append(os.Args, "--help")
	}

	var cli flags
	parser := kong.Must(
		&cli,
		kong.Name("aiac"),
		kong.Description("Artificial Intelligence Infrastructure-as-Code Generator."),
		kong.ConfigureHelp(kong.HelpOptions{
			FlagsLast: true,
		}),
	)

	ctx, err := parser.Parse(os.Args[1:])
	if err != nil {
		if err.Error() == "missing flags: --api-key=STRING" {
			fmt.Fprintln(os.Stderr, `You must provide an OpenAI API key via the --api-key flag, or
the OPENAI_API_KEY environment variable.

Get your API key from https://platform.openai.com/account/api-keys.`)
		} else {
			fmt.Fprintf(os.Stderr, "%v\n", err)
		}

		os.Exit(1)
	}

	switch ctx.Command() {
	case "version":
		fmt.Fprintf(os.Stdout, "aiac version %s\n", libaiac.Version)
		os.Exit(0)
	case "list-models", "list-models <type>":
		printModels(cli)
		os.Exit(0)
	case "get <what>":
		err := generateCode(cli)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %s\n", err)
			os.Exit(1)
		}
		os.Exit(0)
	default:
		fmt.Fprintln(os.Stderr, "Unknown command")
		os.Exit(1)
	}
}

func printModels(cli flags) {
	headerFmt := color.New(color.FgGreen, color.Underline).SprintfFunc()
	columnFmt := color.New(color.FgYellow).SprintfFunc()

	tbl := table.New("Name", "Type", "Maximum Tokens")
	tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)

	var client types.Backend

	switch cli.Backend {
	case libaiac.BackendOpenAI:
		client = &openai.Client{}
	case libaiac.BackendBedrock:
		client = &bedrock.Client{}
	default:
		fmt.Fprintf(os.Stderr, "Unknown backend %s\n", cli.Backend)
		os.Exit(1)
	}

	models := client.ListModels()

	for _, model := range models {
		if cli.ListModels.Type == "" || cli.ListModels.Type == model.Type {
			tbl.AddRow(model.Name, model.Type, model.MaxTokens)
		}
	}

	tbl.Print()
}

var errInvalidInput = errors.New("invalid input, please try again")

func generateCode(cli flags) error { //nolint: funlen, cyclop
	client := libaiac.NewClient(&libaiac.NewClientOptions{
		Backend:    cli.Backend,
		ApiKey:     cli.Get.APIKey,
		URL:        cli.Get.URL,
		APIVersion: cli.Get.APIVersion,
		AWSProfile: cli.Get.AWSProfile,
		AWSRegion:  cli.Get.AWSRegion,
	})

	var model types.Model
	if cli.Get.Model == "" {
		model = client.DefaultModel()
	} else {
		for _, supported := range client.ListModels() {
			if supported.Name == cli.Get.Model {
				model = supported
				break
			}
		}

		if model.Name == "" {
			return fmt.Errorf("%w %q", types.ErrUnsupportedModel, cli.Get.Model)
		}
	}

	spin := spinner.New(
		spinner.CharSets[11],
		100*time.Millisecond, //nolint: gomnd
		spinner.WithWriter(color.Error),
		spinner.WithSuffix("\tGenerating code ..."))

	defer func() {
		if spin.Active() {
			spin.Stop()
		}
	}()

	// NOTE: we are prepending the string "generate sample code for a..."
	// to the prompt, this is meant to ensure that the language model
	// actually generates code.
	prompt := fmt.Sprintf("Generate sample code for a %s", strings.Join(cli.Get.What, " "))

	if cli.Get.ReadmeFile != "" || cli.Get.Full {
		prompt = fmt.Sprintf(
			"Generate sample code for a %s. Include explanations.",
			strings.Join(cli.Get.What, " "),
		)
	}

	var res types.Response
	var err error

	var conversation types.Conversation
	if model.Type == types.ModelTypeChat {
		conversation = client.Chat(model)
	}

	ctx := context.TODO()

ATTEMPTS:
	for {
		spin.Start()

		if conversation != nil {
			res, err = conversation.Send(ctx, prompt)
		} else {
			res, err = client.Complete(ctx, model, prompt)
		}

		options := [][2]string{
			{"r", "retry same prompt"},
			{"y", "copy to clipboard"},
			{"q", "quit"},
		}

		if err != nil {
			spin.Stop()
			fmt.Fprintf(os.Stderr, "Failed generating code: %s\n", err)
		} else {
			spin.Stop()

			stdoutOutput := res.Code
			if cli.Get.Full {
				stdoutOutput = res.FullOutput
			}

			fmt.Fprintln(os.Stdout, stdoutOutput)

			if cli.Get.Quiet {
				if cli.Get.Clipboard {
					clipboard.WriteAll(stdoutOutput)
				}
				break ATTEMPTS
			}

			if conversation != nil {
				options = append(
					[][2]string{
						{"s", "save and exit"},
						{"w", "save and chat"},
						{"c", "continue chatting"},
					},
					options...,
				)
			} else {
				options = append(
					[][2]string{
						{"s", "save and exit"},
					},
					options...,
				)
			}
		}

	PROMPT:
		for {
			fmt.Println()
			for _, opt := range options {
				fmt.Printf(
					"[%s/%s]: %s\n",
					strings.ToUpper(opt[0]), opt[0], opt[1],
				)
			}

			input := promptui.Prompt{
				Label: "Choice",
				Validate: func(s string) error {
					key := strings.ToLower(s)
					for _, opt := range options {
						if opt[0] == key {
							return nil
						}
					}

					return errInvalidInput
				},
			}

			result, err := input.Run()
			if err != nil {
				if errors.Is(err, io.EOF) {
					return nil
				}

				return fmt.Errorf("prompt failed: %w", err)
			}

			choice := strings.ToLower(result)

			switch choice {
			case "r":
				continue ATTEMPTS
			case "q":
				// finish without saving
				return nil
			case "y":
				// copy code to clipboard
				clipboard.WriteAll(res.Code)
				fmt.Fprintf(os.Stderr, "Generated code copied to clipboard.\n")
				continue PROMPT
			case "c":
				// continue chatting
				prompt = newMessage()
				continue ATTEMPTS
			case "s", "w":
				err = saveOutput(cli, res)
				if err != nil {
					return fmt.Errorf("failed saving output: %w", err)
				}

				if choice == "w" {
					prompt = newMessage()
					continue ATTEMPTS
				} else {
					break ATTEMPTS
				}
			}
		}
	}

	return nil
}

func newMessage() string {
	input := promptui.Prompt{
		Label: "New message",
	}

	prompt, err := input.Run()
	for err != nil {
		fmt.Fprintf(os.Stderr, "%s: please try again\n", err)
		prompt, err = input.Run()
	}

	return prompt
}

func saveOutput(cli flags, res types.Response) (err error) {
	if !cli.Get.Quiet && cli.Get.OutputFile == "" {
		input := promptui.Prompt{
			Label: "Enter file path for generated code",
		}

		cli.Get.OutputFile, err = input.Run()
		if err != nil {
			return fmt.Errorf("prompt failed: %w", err)
		}
	}

	var codeSaved, fullSaved bool

	if cli.Get.OutputFile != "" {
		f, err := os.Create(cli.Get.OutputFile)
		if err != nil {
			return fmt.Errorf(
				"failed creating output file %s: %w",
				cli.Get.OutputFile, err,
			)
		}

		fmt.Fprintln(f, res.Code)
		f.Close()

		codeSaved = true
	}

	if !cli.Get.Quiet && cli.Get.ReadmeFile == "" {
		input := promptui.Prompt{
			Label: "Enter file path for full output, or leave empty to ignore",
		}

		cli.Get.ReadmeFile, err = input.Run()
		if err != nil {
			return fmt.Errorf("prompt failed: %w", err)
		}
	}

	if cli.Get.ReadmeFile != "" {
		f, err := os.Create(cli.Get.ReadmeFile)
		if err != nil {
			return fmt.Errorf(
				"failed creating readme file %s: %w",
				cli.Get.ReadmeFile, err,
			)
		}

		fmt.Fprintln(f, res.FullOutput)
		f.Close()

		fullSaved = true
	}

	if codeSaved {
		fmt.Fprintf(os.Stderr, "Code saved successfully to %s\n", cli.Get.OutputFile)
	}
	if fullSaved {
		fmt.Fprintf(os.Stderr, "Full output saved successfully to %s\n", cli.Get.ReadmeFile)
	}

	return nil
}
