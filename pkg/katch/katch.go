package katch

import (
	"context"
	"fmt"
	"strings"
	"time"

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
		tasks = append(tasks, infinityScrollTask(input.ScrollTimes, input.ScrollStep, scrollDelayDur))
	}

	switch input.OutputFormat {
	case OutputFormatPDF:
		tasks = append(tasks, pdfExporterTask(input.PDFLandscape, input.PDFPrintBackground, input.PDFPaperHeight, input.PDFPaperHeight, &output))
	case OutputFormatHTML:
		tasks = append(tasks, htmlExporterTask(&output))
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
