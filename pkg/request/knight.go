package request

import (
	"assedio/pkg/model"
	"fmt"
	"github.com/fatih/color"
	"net/http"
	"net/url"
	"time"
)

type Fighter interface {
	Hit(urls chan *url.URL, results model.Slice) error
}

type Knight struct{}

func (k *Knight) Hit(urls chan *url.URL, results model.Slice) error {
	for url := range urls {
		uri := url.String()
		client := &http.Client{
			Timeout: 1 * time.Minute,
		}
		startReqTime := time.Now()
		resp, err := client.Get(uri)
		requestTotTimeDt := time.Now().Sub(startReqTime)

		logUrl := fmt.Sprintf("%s", url.Path)

		if err != nil {
			color.Red(fmt.Sprintf("[ERR] - Chiamata %s", logUrl))

			results.Append(model.Record{
				Status:   "ERROR",
				Duration: requestTotTimeDt,
				Url:      url,
				Error:    true,
			})

			continue
		}

		if resp.StatusCode == 200 {
			color.Cyan(fmt.Sprintf("[%s] %d ms Chiamata %s", resp.Status, requestTotTimeDt.Milliseconds(), logUrl))
			results.Append(model.Record{
				Status:   resp.Status,
				Duration: requestTotTimeDt,
				Url:      url,
				Error:    false,
			})
		} else {
			color.Magenta(fmt.Sprintf("[%s] %d ms Chiamata %s", resp.Status, requestTotTimeDt.Milliseconds(), logUrl))
			results.Append(model.Record{
				Status:   resp.Status,
				Duration: requestTotTimeDt,
				Url:      url,
				Error:    true,
			})
		}
	}
	return nil
}
