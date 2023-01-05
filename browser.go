package headless

import (
	"context"
	"io"

	"github.com/chromedp/chromedp"
	"github.com/wrk-grp/errnie"
)

/*
Browser wraps the headless Chrome instance.
*/
type Browser struct {
	ctx    context.Context
	cancel context.CancelFunc
}

func NewBrowser() *Browser {
	ctx, cancel := chromedp.NewContext(context.Background())
	return &Browser{ctx, cancel}
}

func (browser *Browser) Read(p []byte) (n int, err error) {
	defer browser.cancel()

	errnie.Handles(chromedp.Run(
		browser.ctx,
		chromedp.Navigate("https://www.google.com"),
		chromedp.CaptureScreenshot(&p),
	))

	return len(p), io.EOF
}
