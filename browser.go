package headless

import (
	"context"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/wrk-grp/errnie"
)

/*
Browser wraps the headless Chrome instance.
*/
type Browser struct {
	ctx      context.Context
	deallocs [3]context.CancelFunc
	cancel   context.CancelFunc
	Human    *Human
}

func NewBrowser(remote string) *Browser {
	errnie.Trace()

	browser := &Browser{}
	allocCtx := context.Background()

	allocCtx, browser.deallocs[2] = chromedp.NewRemoteAllocator(
		allocCtx, remote,
	)

	ctx, dealloc := chromedp.NewContext(allocCtx)
	browser.deallocs[1] = dealloc
	browser.ctx, browser.deallocs[0] = context.WithTimeout(ctx, 30*time.Second)

	return browser
}

func (browser *Browser) Close() {
	// Call all the cancel functions on the context. They were added in reverse
	// order, so a simpel range correctly sequences the cancel calls.
	for _, dealloc := range browser.deallocs {
		dealloc()
	}
}
