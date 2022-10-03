package main

import (
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:     "katch",
		Usage:    "a simple interface for headless chrome",
		Commands: []*cli.Command{},
	}

	app.Commands = append(app.Commands, &cli.Command{
		Name:   "serve",
		Usage:  "starts a simple http interface server",
		Action: server(),
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "listen",
				Aliases: []string{"l"},
				Value:   ":3000",
			},
		},
	})

	app.Commands = append(app.Commands, &cli.Command{
		Name:   "export",
		Usage:  "exports the specified url as pdf, png, jpeg and html",
		Action: exporter(),
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "url",
				Usage:    "the url to capture",
				Aliases:  []string{"u"},
				Required: true,
			},
			&cli.StringFlag{
				Name:  "wait-for",
				Usage: "duration to wait before processing, example: 5s, 10s, 1m ... etc",
			},
			&cli.Int64Flag{
				Name:  "max-exec-time",
				Usage: "the max execution time in seconds",
				Value: 30,
			},
			&cli.Int64Flag{
				Name:    "viewport-width",
				Usage:   "the viewport width",
				Aliases: []string{"vw"},
			},
			&cli.Int64Flag{
				Name:    "viewport-height",
				Usage:   "the viewport height",
				Aliases: []string{"vh"},
			},
			&cli.StringFlag{
				Name:    "format",
				Usage:   "the output format",
				Aliases: []string{"f"},
			},
			&cli.StringFlag{
				Name:     "output",
				Usage:    "the output filename",
				Aliases:  []string{"o"},
				Required: true,
			},
			&cli.BoolFlag{
				Name:  "png-full-page",
				Usage: "whether to take a full page screenshot regardless viewport dimensions (in case of png output only)",
				Value: true,
			},
			&cli.BoolFlag{
				Name:  "pdf-landscape",
				Usage: "whether to print the page in landscape mode or not",
				Value: true,
			},
			&cli.BoolFlag{
				Name:  "pdf-background",
				Usage: "whether to include the page background(s) in the result or not",
				Value: true,
			},
			&cli.Float64Flag{
				Name:  "pdf-paper-height",
				Usage: "the pdf paper height, 0 means chrome's default",
				Value: 0,
			},
			&cli.Float64Flag{
				Name:  "pdf-paper-width",
				Usage: "the pdf paper width, 0 means chrome's default",
				Value: 0,
			},
			&cli.Int64Flag{
				Name:  "scroll-step",
				Usage: "the scroll step (try to read more at: https://developer.mozilla.org/en-US/docs/Web/API/Element/scrollTop)",
			},
			&cli.Int64Flag{
				Name:  "scroll-times",
				Usage: "how many times do you want to scroll? 0 means don't scroll, -1 means do infinity scroll",
				Value: 0,
			},
			&cli.StringFlag{
				Name:  "scroll-delay",
				Usage: "the delay duration between each scroll step, accepts duration formats (1s, 100ms, 1h, ... etc)",
			},
		},
	})

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
