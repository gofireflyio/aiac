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
	APIKey     string `help:"OpenAI API key" optional:"" env:"OPENAI_API_KEY"`
	OutputFile string `help:"Output file to push resulting code to, defaults to stdout" default:"-" type:"path" short:"o"`
	ReadmeFile string `help:"Markdown file to push explanations to (available only in ChatGPT mode)" optional:"" type:"path" short:"r"`
	Save       bool   `help:"Save AIaC response without retry prompt" default:false short:"s"`
	Quiet      bool   `help:"Print AIaC response to stdout and exit (non-interactive mode)" default:false short:"q"`
	Get        struct {
		What []string `arg:"" help:"Which IaC template to generate"`
	} `cmd:"" help:"Generate IaC code" aliases:"generate"`

	// ChatGPT authentication is experimental
	ChatGPT             bool   `help:"Use ChatGPT mode instead of the OpenAI API (requires --session-token)" default:false hidden:""`
	SessionToken        string `help:"Session token for ChatGPT mode" optional:"" hidden:"" env:"CHATGPT_SESSION_TOKEN"`
	CloudflareClearance string `help:"Cloudflare clearance token for ChatGPT mode" optional:"" hidden:"" env:"CLOUDFLARE_CLEARANCE_TOKEN"`
	CloudflareBm        string `help:"Cloudflare bm token for ChatGPT mode" optional:"" hidden:"" env:"CLOUDFLARE_BM_TOKEN"`
	UserAgent           string `help:"Cloudflare tokens user agent ChatGPT mode" optional:"" hidden:"" env:"USER_AGENT"`
}

func main() {
	if len(os.Args) < 2 {
		os.Args = append(os.Args, "--help")
	}
	var cli flags
	cmd := kong.Parse(&cli)

	if cmd.Command() != "get <what>" {
		_, _ = fmt.Fprintln(os.Stderr, "Unknown command")
		os.Exit(1)
	}

	var token string

	if !cli.ChatGPT {
		token = cli.APIKey
		if token == "" {
			_, _ = fmt.Fprintf(os.Stderr, "You must provide an OpenAI API key\n")
			os.Exit(1)
		}
	} else {
		token = cli.SessionToken
		if token == "" {
			var ok bool
			token, ok = os.LookupEnv("CHATGPT_SESSION_TOKEN")

			if !ok {
				_, _ = fmt.Fprintf(os.Stderr, "You must provide a ChatGPT session token\n")
				os.Exit(1)
			}
		}
	}

	client := libaiac.NewClient(&libaiac.AIACClientInput{
		ChatGPT:             cli.ChatGPT,
		Token:               token,
		CloudflareClearance: cli.CloudflareClearance,
		CloudflareBm:        cli.CloudflareBm,
		UserAgent:           cli.UserAgent,
	})

	shouldRetry := !cli.Save

	err := client.Ask(
		context.TODO(),
		// NOTE: we are prepending the word "generate" to the prompt, this
		// ensures the language model actually generates code. The word "get",
		// on the other hand, doesn't necessarily result in code being generated.
		fmt.Sprintf("generate %s", strings.Join(cli.Get.What, " ")),
		shouldRetry,
		cli.Quiet,
		cli.OutputFile,
		cli.ReadmeFile,
	)
	if err != nil {
		if errors.Is(err, libaiac.ErrNoCode) {
			_, _ = fmt.Fprintln(
				os.Stderr,
				"It doesn't look like ChatGPT generated any code, please make "+
					"sure that you're prompt properly guides ChatGPT to do so.",
			)
		} else {
			_, _ = fmt.Fprintf(os.Stderr, "Request failed: %s\n", err)
		}
		os.Exit(1)
	}
}
