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
	"github.com/gofireflyio/aiac/v5/libaiac"
	"github.com/gofireflyio/aiac/v5/libaiac/types"
	"github.com/manifoldco/promptui"
)

type flags struct {
	Config     string   `help:"Configuration file path" type:"path" short:"c"`
	Backend    string   `help:"Backend to use" short:"b"`
	OutputFile string   `help:"Output file to push resulting code to" optional:"" type:"path" short:"o"`         //nolint: lll
	ReadmeFile string   `help:"Readme file to push entire Markdown output to" optional:"" type:"path" short:"r"` //nolint: lll
	Quiet      bool     `help:"Non-interactive mode, print/save output and exit" default:"false" short:"q"`      //nolint: lll
	Full       bool     `help:"Print full Markdown output to stdout" default:"false" short:"f"`                  //nolint: lll
	Model      string   `help:"Model to use" short:"m"`
	What       []string `arg:"" optional:"" help:"Which IaC template to generate"`
	Clipboard  bool     `help:"Copy generated code to clipboard (in --quiet mode)"`
	ListModels bool     `help:"List supported models and exit"`
	Timeout    int      `help:"Generate code timeout in second" default:"60"`
	Version    bool     `help:"Print aiac version and exit"`
}

func main() {
	var cli flags
	parser := kong.Must(
		&cli,
		kong.Name("aiac"),
		kong.Description("Artificial Intelligence Infrastructure-as-Code Generator."),
		kong.ConfigureHelp(kong.HelpOptions{
			FlagsLast: true,
		}),
	)

	_, err := parser.Parse(os.Args[1:])
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	if cli.Version {
		fmt.Fprintf(os.Stdout, "aiac version %s\n", libaiac.Version)
		os.Exit(0)
	}

	aiac, err := libaiac.New(cli.Config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed loading aiac client: %s\n", err)
		os.Exit(1)
	}

	if cli.ListModels {
		err := printModels(aiac, cli)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed listing models: %s\n", err)
			os.Exit(1)
		}

		os.Exit(0)
	}

	err = generateCode(aiac, cli)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}
	os.Exit(0)
}

func printModels(aiac *libaiac.Aiac, cli flags) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	models, err := aiac.ListModels(ctx, cli.Backend)
	if err != nil {
		return err
	}

	for _, model := range models {
		fmt.Println(model)
	}

	return nil
}

var errInvalidInput = errors.New("invalid input, please try again")

func generateCode(aiac *libaiac.Aiac, cli flags) error { //nolint: funlen, cyclop
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(cli.Timeout)*time.Second)
	defer cancel()

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

	// If the prompt starts with the word "get" or "generate", remove it. This
	// is here for backwards compatibility purposes, as previous versions used
	// these words as command names (that weren't truly part of the prompt), so
	// people may be used to adding them and we don't want them to actually be
	// in the prompt.
	if strings.ToLower(cli.What[0]) == "get" ||
		strings.ToLower(cli.What[0]) == "generate" {
		cli.What = cli.What[1:]
	}

	// NOTE: we are prepending the string "generate sample code for a..."
	// to the prompt, this is meant to ensure that the language model
	// actually generates code.
	prompt := fmt.Sprintf("Generate sample code for a %s", strings.Join(cli.What, " "))

	if cli.ReadmeFile != "" || cli.Full {
		prompt = fmt.Sprintf(
			"Generate sample code for a %s. Include explanations.",
			strings.Join(cli.What, " "),
		)
	}

	var res types.Response

	chat, err := aiac.Chat(ctx, cli.Backend, cli.Model)
	if err != nil {
		return fmt.Errorf("failed starting chat: %w", err)
	}

ATTEMPTS:
	for {
		spin.Start()

		res, err = chat.Send(ctx, prompt)

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
			if cli.Full {
				stdoutOutput = res.FullOutput
			}

			fmt.Fprintln(os.Stdout, stdoutOutput)

			if cli.Quiet {
				if cli.Clipboard {
					clipboard.WriteAll(stdoutOutput)
				}
				break ATTEMPTS
			}

			options = append(
				[][2]string{
					{"s", "save and exit"},
					{"w", "save and chat"},
					{"c", "continue chatting"},
				},
				options...,
			)
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
	if !cli.Quiet && cli.OutputFile == "" {
		input := promptui.Prompt{
			Label: "Enter file path for generated code",
		}

		cli.OutputFile, err = input.Run()
		if err != nil {
			return fmt.Errorf("prompt failed: %w", err)
		}
	}

	var codeSaved, fullSaved bool

	if cli.OutputFile != "" {
		f, err := os.Create(cli.OutputFile)
		if err != nil {
			return fmt.Errorf(
				"failed creating output file %s: %w",
				cli.OutputFile, err,
			)
		}

		fmt.Fprintln(f, res.Code)
		f.Close()

		codeSaved = true
	}

	if !cli.Quiet && cli.ReadmeFile == "" {
		input := promptui.Prompt{
			Label: "Enter file path for full output, or leave empty to ignore",
		}

		cli.ReadmeFile, err = input.Run()
		if err != nil {
			return fmt.Errorf("prompt failed: %w", err)
		}
	}

	if cli.ReadmeFile != "" {
		f, err := os.Create(cli.ReadmeFile)
		if err != nil {
			return fmt.Errorf(
				"failed creating readme file %s: %w",
				cli.ReadmeFile, err,
			)
		}

		fmt.Fprintln(f, res.FullOutput)
		f.Close()

		fullSaved = true
	}

	if codeSaved {
		fmt.Fprintf(os.Stderr, "Code saved successfully to %s\n", cli.OutputFile)
	}
	if fullSaved {
		fmt.Fprintf(os.Stderr, "Full output saved successfully to %s\n", cli.ReadmeFile)
	}

	return nil
}
