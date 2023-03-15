package headless

import (
	"net/url"

	"github.com/wrk-grp/errnie"
)

func NewURL(urlstr string) *url.URL {
	u, err := url.Parse(urlstr)
	errnie.Handles(err)
	return u
}
