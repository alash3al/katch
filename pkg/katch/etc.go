package katch

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/chromedp/cdproto/dom"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/cdproto/runtime"
	"github.com/chromedp/chromedp"
	"time"
)

func extractDocumentElementScrollTop(ctx context.Context) (float64, error) {
	result, exception, err := runtime.Evaluate("document.documentElement.scrollTop").Do(ctx)
	if err != nil {
		return 0, err
	}
	if exception != nil {
		return 0, exception
	}

	var val float64

	err = json.Unmarshal(result.Value, &val)

	return val, err
}

func setDocumentElementScrollTop(ctx context.Context, step int64) error {
	_, exception, err := runtime.Evaluate(fmt.Sprintf("document.documentElement.scrollTop += %d", step)).Do(ctx)
	if err != nil {
		return err
	}
	if exception != nil {
		return exception
	}
	return nil
}

func sleep(ctx context.Context, d time.Duration) error {
	timer := time.NewTimer(d)
	select {
	case <-ctx.Done():
		if !timer.Stop() {
			<-timer.C
		}
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}

func infinityScrollTask(scrollMaxTimes int64, scrollStep int64, scrollDelayDuration time.Duration) chromedp.ActionFunc {
	return func(ctx context.Context) error {
		reachedBottom := false
		prevDocumentElementScrollTop := float64(0)
		scrollTimes := int64(0)

		for !reachedBottom && !(scrollMaxTimes > -1 && scrollTimes >= scrollMaxTimes) {
			if err := setDocumentElementScrollTop(ctx, scrollStep); err != nil {
				return err
			}

			if err := sleep(ctx, scrollDelayDuration); err != nil {
				return err
			}

			newDocumentElementScrollTop, err := extractDocumentElementScrollTop(ctx)
			if err != nil {
				return err
			}

			if newDocumentElementScrollTop == prevDocumentElementScrollTop {
				reachedBottom = true
			}

			prevDocumentElementScrollTop = newDocumentElementScrollTop
			scrollTimes++
		}

		return nil
	}
}

func pdfExporterTask(landscape bool, printBackground bool, paperHeight float64, paperWidth float64, output *[]byte) chromedp.ActionFunc {
	return func(ctx context.Context) error {
		pdfParams := page.PrintToPDF()
		pdfParams.Landscape = landscape
		pdfParams.PrintBackground = printBackground
		pdfParams.PaperHeight = paperHeight
		pdfParams.PaperWidth = paperWidth

		buf, _, err := pdfParams.Do(ctx)
		if err != nil {
			return err
		}

		*output = buf

		return nil
	}
}

func htmlExporterTask(output *[]byte) chromedp.ActionFunc {
	return func(ctx context.Context) error {
		node, err := dom.GetDocument().Do(ctx)
		if err != nil {
			return err
		}

		outputStr, err := dom.GetOuterHTML().WithNodeID(node.NodeID).Do(ctx)
		if err != nil {
			return err
		}

		*output = []byte(outputStr)

		return nil
	}
}
