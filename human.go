package headless

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/wrk-grp/errnie"
)

type Human struct {
	tick int64
	ctx  context.Context
}

func NewHuman(ctx context.Context) *Human {
	return &Human{1, ctx}
}

func (human *Human) Wait() {
	human.delay(80 + rand.Int63n(80))
}

func (human *Human) Click(sel string) {
	errnie.Handles(chromedp.Run(human.ctx,
		chromedp.Click(sel, chromedp.NodeVisible),
	))
}

func (human *Human) Type(sel, str string) {
	for _, char := range str {
		human.tick++
		human.delay(rand.Int63n(10))

		errnie.Handles(chromedp.Run(human.ctx,
			chromedp.SendKeys(sel, string(char)),
		))
	}
}

func (human *Human) delay(add int64) {
	var t int64

	if (human.tick+10)%(rand.Int63n(32)+1) == 0 {
		t = (rand.Int63n(32) + 500) - (rand.Int63n(32) + 300)
	}

	if (human.tick+10)%(rand.Int63n(10)+1) == 0 {
		t = (rand.Int63n(32) + 300) - (rand.Int63n(32) + 200)
	}

	if t == 0 {
		t = (rand.Int63n(32) + 150) - (rand.Int63n(32) + 70)
	}

	fmt.Printf("[%d]", t)
	time.Sleep(time.Duration(t+add) * time.Millisecond)
}
