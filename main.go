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
	APIKey     string `help:"OpenAI API key" required:"" env:"OPENAI_API_KEY"`
	OutputFile string `help:"Output file to push resulting code to, defaults to stdout" default:"-" type:"path" short:"o"`
	Save       bool   `help:"Save AIaC response without retry prompt" default:false short:"s"`
	Quiet      bool   `help:"Print AIaC response to stdout and exit (non-interactive mode)" default:false short:"q"`
	Get        struct {
		What []string `arg:"" help:"Which IaC template to generate"`
	} `cmd:"" help:"Generate IaC code" aliases:"generate"`
}

func main() {
	if len(os.Args) < 2 {
		os.Args = append(os.Args, "--help")
	}
	var cli flags
	cmd := kong.Parse(&cli)

	if cmd.Command() != "get <what>" {
		fmt.Fprintln(os.Stderr, "Unknown command")
		os.Exit(1)
	}

	client := libaiac.NewClient(cli.APIKey)

	err := client.Ask(
		context.TODO(),
		// NOTE: we are prepending the word "generate" to the prompt, this
		// ensures the language model actually generates code. The word "get",
		// on the other hand, doesn't necessarily result in code being generated.
		fmt.Sprintf("generate %s", strings.Join(cli.Get.What, " ")),
		!cli.Save,
		cli.Quiet,
		cli.OutputFile,
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Request failed: %s\n", err)
		os.Exit(1)
	}
}
