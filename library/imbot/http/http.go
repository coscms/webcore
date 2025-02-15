package http

import (
	"errors"

	"github.com/webx-top/restyclient"
)

func Send(url string, m interface{}) ([]byte, error) {
	client := restyclient.Retryable()
	resp, err := client.SetBody(m).SetHeader(`Content-Type`, `application/json`).Post(url)
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		err = errors.New(resp.String())
		return nil, err
	}
	return resp.Body(), err
}
