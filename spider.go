package headless

import (
	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
	"github.com/wrk-grp/errnie"
)

type spider struct {
	browser  *Browser
	links    []string
	images   []string
	requests map[string]network.RequestID
}

func (browser *Browser) Spider(urlstr string) {
	errnie.Trace()
	errnie.Debugs("SPIDER", urlstr)

	spider := &spider{browser: browser, requests: make(map[string]network.RequestID)}
	spider.imageListener()
	spider.generate(urlstr)
}

func (generator *spider) generate(urlstr string) {
	var (
		links  []*cdp.Node
		images []*cdp.Node
	)

	errnie.Debugs("SCRAPING", urlstr)

	errnie.Handles(chromedp.Run(generator.browser.ctx,
		chromedp.Navigate(urlstr),
		chromedp.WaitReady("body"),
		chromedp.Nodes("//img", &images),
		chromedp.Nodes("//a", &links),
	))

	generator.getImages(images)

	for _, link := range links {
		if href := link.AttributeValue("href"); href != "" {
			var found bool

			for _, l := range generator.links {
				if href == l {
					found = true
				}
			}

			if !found {
				errnie.Informs("LINK", href)
				generator.links = append(generator.links, href)
			}
		}
	}

	var task string
	task, generator.links = generator.links[0], append(generator.links[1:], generator.links[0])
	generator.generate(task)
}
