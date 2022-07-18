package katch

import (
	"context"
	"fmt"
	"time"

	"github.com/chromedp/cdproto/dom"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
)

// Input represents a Katch input
type Input struct {
	URL                   string       `json:"url" query:"url" form:"url"`
	WaitForElementVisible string       `json:"wait_element_visible" query:"wait_element_visible" form:"wait_element_visible"`
	MaxExecTime           int          `json:"max_exec_time" query:"max_exec_time" form:"max_exec_time"`
	ViewportWidth         int64        `json:"viewport_width" query:"viewport_width" form:"viewport_width"`
	ViewportHeight        int64        `json:"viewport_height" query:"viewport_height" form:"viewport_height"`
	OutputFormat          OutputFormat `json:"format" query:"format" form:"format"`
}

// OutputFormat represents a requested format
type OutputFormat string

// available output formats
const (
	OutputFormatPDF  OutputFormat = "pdf"
	OutputFormatPNG  OutputFormat = "png"
	OutputFormatJPEG OutputFormat = "jpeg"
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

	if input.WaitForElementVisible != "" {
		tasks = append(tasks, chromedp.WaitVisible(input.WaitForElementVisible, chromedp.ByQuery))
	}

	switch input.OutputFormat {
	case OutputFormatPDF:
		tasks = append(tasks, chromedp.ActionFunc(func(ctx context.Context) error {
			buf, _, err := page.PrintToPDF().Do(ctx)
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
		tasks = append(tasks, chromedp.FullScreenshot(&output, 100))
	case OutputFormatJPEG:
		tasks = append(tasks, chromedp.FullScreenshot(&output, 90))
	default:
		return nil, fmt.Errorf("unsupported output format (%s)", input.OutputFormat)
	}

	if err := chromedp.Run(ctx, tasks); err != nil {
		return nil, err
	}

	return output, nil
}
