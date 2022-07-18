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
				Usage: "the element to let chrome wait for before processing",
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
		},
	})

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
