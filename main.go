package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/alecthomas/kong"
	"github.com/gofireflyio/aiac/v2/libaiac"
)

type flags struct {
	APIKey     string        `help:"OpenAI API key" required:"" env:"OPENAI_API_KEY"`
	ListModels struct{}      `cmd:"" help:"List supported models"`
	OutputFile string        `help:"Output file to push resulting code to, defaults to stdout" default:"-" type:"path" short:"o"` //nolint: lll
	Quiet      bool          `help:"Print AIaC response to stdout and exit (non-interactive mode)" default:"false" short:"q"`
	Save       bool          `help:"Save AIaC response without retry prompt" default:"false" short:"s"`
	Full       bool          `help:"Return full output, including explanations, if any" default:"false" short:"f"`
	Model      libaiac.Model `help:"Model to use, default to \"gpt-3.5-turbo\""`
	Get        struct {
		What []string `arg:"" help:"Which IaC template to generate"`
	} `cmd:"" help:"Generate IaC code" aliases:"generate"`
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

	if ctx.Command() == "list-models" {
		for _, model := range libaiac.SupportedModels {
			fmt.Println(model)
		}

		os.Exit(0)
	}

	if ctx.Command() != "get <what>" {
		fmt.Fprintln(os.Stderr, "Unknown command")
		os.Exit(1)
	}

	client := libaiac.NewClient(cli.APIKey).
		SetFull(cli.Full)

	if cli.Model != "" {
		client.SetModel(cli.Model)
	}

	err = client.Ask(
		context.TODO(),
		// NOTE: we are prepending the word "generate" to the prompt, this
		// ensures the language model actually generates code. The word "get",
		// on the other hand, doesn't necessarily result in code being generated.
		fmt.Sprintf("generate sample code for a %s", strings.Join(cli.Get.What, " ")),
		!cli.Save,
		cli.Quiet,
		cli.OutputFile,
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Request failed: %s\n", err)
		os.Exit(1)
	}
}
