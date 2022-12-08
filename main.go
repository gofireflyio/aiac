package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/alecthomas/kong"
	"github.com/gofireflyio/aiac/libaiac"
)

var cli struct {
	SessionToken string   `help:"ChatGPT session token" required:""`
	OutputFile   string   `help:"Output file to push resulting code to, defaults to stdout" default:"-" type:"path"`
	ReadmeFile   string   `help:"Markdown file to update with explanations" type:"path"`
	Ask          []string `arg:"" help:"What to ask ChatGPT to do"`
}

func main() {
	kong.Parse(&cli)

	client := libaiac.NewClient(cli.SessionToken)

	err := client.Ask(
		context.TODO(),
		strings.Join(cli.Ask, " "),
		cli.OutputFile,
		cli.ReadmeFile,
	)
	if err != nil {
        if err == libaiac.ErrNoCode {
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
