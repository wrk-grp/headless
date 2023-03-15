package headless

import (
	"bytes"
	"log"
	"strings"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
	"github.com/theapemachine/wrkspc/brazil"
	"github.com/wrk-grp/errnie"
)

func (generator *spider) getImages(images []*cdp.Node) {
	for _, img := range images {
		var imgstr string

		if src := img.AttributeValue("src"); src != "" {
			imgstr = cleanURL(src)
		}

		if href := img.AttributeValue("href"); href != "" {
			imgstr = cleanURL(href)
		}

		var found bool
		for _, i := range generator.images {
			if i == imgstr {
				found = true
			}
		}

		if imgstr != "" && !found {
			generator.images = append(generator.images, imgstr)
		}
	}
}

func cleanURL(urlstr string) string {
	if urlstr[0] == '/' {
		return "https:" + urlstr
	}

	return urlstr
}

func (generator *spider) imageListener() {
	var requestID network.RequestID

	// set up a listener to watch the network events and close the channel when
	// complete the request id matching is important both to filter out
	// unwanted network events and to reference the downloaded file later
	chromedp.ListenTarget(generator.browser.ctx, func(v interface{}) {
		switch ev := v.(type) {
		case *network.EventRequestWillBeSent:
			log.Printf("EventRequestWillBeSent: %v: %v", ev.RequestID, ev.Request.URL)
			for _, img := range generator.images {
				if ev.Request.URL == img {
					generator.requests[img] = ev.RequestID
				}
			}
		case *network.EventLoadingFinished:
			log.Printf("EventLoadingFinished: %v", ev.RequestID)
			for url, reqID := range generator.requests {
				if ev.RequestID == reqID {
					buf, err := network.GetResponseBody(requestID).Do(generator.browser.ctx)
					errnie.Handles(err)

					u := NewURL(url)
					errnie.Debugs("WRITE", u.Host+u.Path)
					chunks := strings.Split(u.Path, "/")
					name := chunks[len(chunks)-1]

					if strings.Contains(name, ".jpg") || strings.Contains(name, ".jpeg") || strings.Contains(name, ".png") {
						brazil.NewFile(
							"data/images/"+u.Host+strings.Join(chunks[0:len(chunks)-1], "/"),
							name,
							bytes.NewBuffer(buf),
						)
					}
				}
			}
		}
	})
}
