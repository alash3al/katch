package main

import (
	"io/ioutil"

	"github.com/alash3al/katch/pkg/katch"
	"github.com/urfave/cli/v2"
)

func exporter() cli.ActionFunc {
	return func(ctx *cli.Context) error {
		input := katch.Input{
			URL:            ctx.String("url"),
			WaitFor:        ctx.String("wait-for"),
			MaxExecTime:    ctx.Int("max-exec-time"),
			ViewportWidth:  ctx.Int64("viewport-width"),
			ViewportHeight: ctx.Int64("viewport-height"),
			OutputFormat:   katch.OutputFormat(ctx.String("format")),
		}

		output, err := katch.Katch(ctx.Context, input)
		if err != nil {
			return err
		}

		return ioutil.WriteFile(ctx.String("output"), output, 0777)
	}
}
