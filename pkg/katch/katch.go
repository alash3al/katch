package katch

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/chromedp/cdproto/runtime"
	"strings"
	"time"

	"github.com/chromedp/cdproto/dom"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
)

// Input represents a Katch input
type Input struct {
	URL                string       `json:"url" query:"url" form:"url"`
	WaitFor            string       `json:"wait_for" query:"wait_for" form:"wait_for"`
	MaxExecTime        int          `json:"max_exec_time" query:"max_exec_time" form:"max_exec_time"`
	ViewportWidth      int64        `json:"viewport_width" query:"viewport_width" form:"viewport_width"`
	ViewportHeight     int64        `json:"viewport_height" query:"viewport_height" form:"viewport_height"`
	OutputFormat       OutputFormat `json:"format" query:"format" form:"format"`
	PNGFullPage        bool         `json:"png_full_page" query:"png_full_page" form:"png_full_page"`
	PDFLandscape       bool         `json:"pdf_landscape" query:"pdf_landscape" form:"pdf_landscape"`
	PDFPrintBackground bool         `json:"pdf_print_background" query:"pdf_print_background" form:"pdf_print_background"`
	PDFPaperHeight     float64      `json:"pdf_paper_height" query:"pdf_paper_height" form:"pdf_paper_height"`
	PDFPaperWidth      float64      `json:"pdf_paper_width" query:"pdf_paper_width" form:"pdf_paper_width"`
	ScrollStep         int64        `json:"scroll_step" query:"scroll_step" form:"scroll_step"`
	ScrollDelay        string       `json:"scroll_delay" query:"scroll_delay" form:"scroll_delay"`
	ScrollTimes        int64        `json:"scroll_times" query:"scroll_times" form:"scroll_times"`
}

// OutputFormat represents a requested format
type OutputFormat string

// available output formats
const (
	OutputFormatPDF  OutputFormat = "pdf"
	OutputFormatPNG  OutputFormat = "png"
	OutputFormatHTML OutputFormat = "html"
)

// Katch exports the specified url in the input
func Katch(ctx context.Context, input Input) ([]byte, error) {
	var output []byte

	ctx, cancel := chromedp.NewContext(ctx)
	defer cancel()

	if input.MaxExecTime > 0 {
		ctx, cancel = context.WithTimeout(ctx, time.Duration(input.MaxExecTime)*time.Second)
		defer cancel()
	}

	tasks := chromedp.Tasks{}

	if input.ViewportHeight > 0 && input.ViewportWidth > 0 {
		tasks = append(tasks, chromedp.EmulateViewport(input.ViewportWidth, input.ViewportHeight))
	}

	tasks = append(tasks, chromedp.Navigate(input.URL))

	if strings.TrimSpace(input.WaitFor) != "" {
		waitForDur, err := time.ParseDuration(input.WaitFor)
		if err != nil {
			return nil, err
		}

		tasks = append(tasks, chromedp.Sleep(waitForDur))
	}

	if input.ScrollStep > 0 && input.ScrollTimes != 0 {
		scrollDelayDur, err := time.ParseDuration(input.ScrollDelay)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, chromedp.ActionFunc(func(ctx context.Context) error {
			return infinityScroll(ctx, input.ScrollTimes, input.ScrollStep, scrollDelayDur)
		}))
	}

	switch input.OutputFormat {
	case OutputFormatPDF:
		tasks = append(tasks, chromedp.ActionFunc(func(ctx context.Context) error {
			pdfParams := page.PrintToPDF()
			pdfParams.Landscape = input.PDFLandscape
			pdfParams.PrintBackground = input.PDFPrintBackground
			pdfParams.PaperHeight = input.PDFPaperHeight
			pdfParams.PaperWidth = input.PDFPaperWidth

			buf, _, err := pdfParams.Do(ctx)
			if err != nil {
				return err
			}

			output = buf

			return nil
		}))
	case OutputFormatHTML:
		tasks = append(tasks, chromedp.ActionFunc(func(ctx context.Context) error {
			node, err := dom.GetDocument().Do(ctx)
			if err != nil {
				return err
			}

			outputStr, err := dom.GetOuterHTML().WithNodeID(node.NodeID).Do(ctx)
			if err != nil {
				return err
			}

			output = []byte(outputStr)

			return nil
		}))
	case OutputFormatPNG:
		if input.PNGFullPage {
			tasks = append(tasks, chromedp.FullScreenshot(&output, 100))
		} else {
			tasks = append(tasks, chromedp.CaptureScreenshot(&output))
		}
	default:
		return nil, fmt.Errorf("unsupported output format (%s)", input.OutputFormat)
	}

	if err := chromedp.Run(ctx, tasks); err != nil {
		return nil, err
	}

	return output, nil
}

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

func documentElementScrollTop(ctx context.Context, step int64) error {
	_, exception, err := runtime.Evaluate(fmt.Sprintf("document.documentElement.scrollTop += %d", step)).Do(ctx)
	if err != nil {
		return err
	}
	if exception != nil {
		return exception
	}
	return nil
}

func infinityScroll(ctx context.Context, scrollMaxTimes int64, scrollStep int64, scrollDelayDuration time.Duration) error {
	reachedBottom := false
	prevDocumentElementScrollTop := float64(0)
	scrollTimes := int64(0)
	for !reachedBottom && !(scrollMaxTimes > -1 && scrollTimes >= scrollMaxTimes) {
		if err := documentElementScrollTop(ctx, scrollStep); err != nil {
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
