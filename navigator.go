package headless

import (
	"github.com/chromedp/chromedp"
	"github.com/wrk-grp/errnie"
)

func (browser *Browser) Navigate(url string) {
	errnie.Trace()

	errnie.Handles(chromedp.Run(browser.ctx,
		chromedp.Navigate(url),
	))
}

func (browser *Browser) Extract(sel string) string {
	var body string

	errnie.Handles(chromedp.Run(browser.ctx,
		chromedp.Evaluate(`new XMLSerializer().serializeToString(document)`, &body),
	))

	return body
}
