package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/alecthomas/kong"
	"github.com/gofireflyio/aiac/libaiac"
)

type flags struct {
	APIKey       string `help:"OpenAI API key (env: OPENAI_API_KEY)" optional:""`
	SessionToken string `help:"Session token for ChatGPT (env: CHATGPT_SESSION_TOKEN)" optional:""`
	ChatGPT      bool   `help:"Use ChatGPT instead of the OpenAI API (requires --session-token)" default:false`
    OutputFile string   `help:"Output file to push resulting code to, defaults to stdout" default:"-" type:"path"`
    ReadmeFile string   `help:"Markdown file to push explanations to" optional:"" type:"path"`
	Get          struct {
		What       []string `arg:"" help:"What to ask ChatGPT to generate"`
	} `cmd:"" help:"Generate IaC code" aliases:"generate"`
}

func main() {
    var cli flags
	cmd := kong.Parse(&cli)

	if cmd.Command() != "get <what>" {
		fmt.Fprintln(os.Stderr, "Unknown command")
		os.Exit(1)
	}

	var token string

	if !cli.ChatGPT {
        token = cli.APIKey
        if token == "" {
            var ok bool
            token, ok = os.LookupEnv("OPENAI_API_KEY")

            if !ok {
                fmt.Fprintf(os.Stderr, "You must provide an OpenAI API key\n")
                os.Exit(1)
            }
        }
	} else {
        token = cli.SessionToken
        if token == "" {
            var ok bool
            token, ok = os.LookupEnv("CHATGPT_SESSION_TOKEN")

            if !ok {
                fmt.Fprintf(os.Stderr, "You must provide a ChatGPT session token\n")
                os.Exit(1)
            }
        }
	}

	client := libaiac.NewClient(cli.ChatGPT, token)

	err := client.Ask(
		context.TODO(),
        // NOTE: we are prepending the word "generate" to the prompt, this
        // ensures the language model actually generates code. The word "get",
        // on the other hand, doesn't necessarily result in code being generated.
		fmt.Sprintf("generate %s", strings.Join(cli.Get.What, " ")),
		cli.OutputFile,
		cli.ReadmeFile,
	)
	if err != nil {
		if errors.Is(err, libaiac.ErrNoCode) {
			fmt.Fprintln(
				os.Stderr,
				"It doesn't look like ChatGPT generated any code, please make "+
					"sure that you're prompt properly guides ChatGPT to do so.",
			)
		} else {
			fmt.Fprintf(os.Stderr, "Request failed: %s\n", err)
		}
		os.Exit(1)
	}
}
