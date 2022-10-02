package main

import (
	"io/ioutil"

	"github.com/alash3al/katch/pkg/katch"
	"github.com/urfave/cli/v2"
)

func exporter() cli.ActionFunc {
	return func(ctx *cli.Context) error {
		input := katch.Input{
			URL:                ctx.String("url"),
			WaitFor:            ctx.String("wait-for"),
			MaxExecTime:        ctx.Int("max-exec-time"),
			ViewportWidth:      ctx.Int64("viewport-width"),
			ViewportHeight:     ctx.Int64("viewport-height"),
			OutputFormat:       katch.OutputFormat(ctx.String("format")),
			PNGFullPage:        ctx.Bool("png-full-page"),
			PDFLandscape:       ctx.Bool("pdf-landscape"),
			PDFPrintBackground: ctx.Bool("pdf-background"),
			PDFPaperHeight:     ctx.Float64("pdf-paper-height"),
			PDFPaperWidth:      ctx.Float64("pdf-paper-width"),
			ScrollStep:         ctx.Int64("scroll-step"),
			ScrollDelay:        ctx.String("scroll-delay"),
			ScrollTimes:        ctx.Int64("scroll-times"),
		}

		output, err := katch.Katch(ctx.Context, input)
		if err != nil {
			return err
		}

		return ioutil.WriteFile(ctx.String("output"), output, 0777)
	}
}
