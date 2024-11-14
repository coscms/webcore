package captcha_go

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/admpub/log"
	"github.com/coscms/captcha"
	"github.com/webx-top/com"
	"github.com/webx-top/echo"
	"github.com/webx-top/restyclient"
)

func NewStoreAPI(apiURL string, secret string) captcha.Storer {
	return &storeAPI{
		apiURL: apiURL,
		secret: secret,
	}
}

type storeAPI struct {
	apiURL string
	secret string
}

func (a *storeAPI) Put(ctx context.Context, key string, val interface{}, timeout int64) error {
	_, err := a.request(http.MethodPost, storeAPIRequest{
		Key:     key,
		Val:     com.String(val),
		Timeout: timeout,
	})
	return err
}

func (a *storeAPI) Get(ctx context.Context, key string, value interface{}) error {
	res, err := a.request(http.MethodGet, storeAPIRequest{
		Key: key,
	})
	if err != nil {
		return err
	}
	*(value.(*[]byte)) = com.Str2bytes(res.String(`Data`))
	return nil
}

func (a *storeAPI) Delete(ctx context.Context, key string) error {
	_, err := a.request(http.MethodDelete, storeAPIRequest{
		Key: key,
	})
	return err
}

func (a *storeAPI) request(method string, body storeAPIRequest) (echo.H, error) {
	req := restyclient.Retryable()
	if len(a.secret) > 0 {
		b, _ := com.JSONEncode(body)
		token := com.Token(a.secret, b)
		req.SetAuthToken(token)
	}
	res := echo.H{}
	reqURL := a.apiURL
	if method == http.MethodGet {
		q := url.Values{}
		q.Set(`key`, body.Key)
		if len(body.Val) > 0 {
			q.Set(`val`, body.Val)
		}
		if body.Timeout > 0 {
			q.Set(`timeout`, com.String(body.Timeout))
		}
		reqURL += `?` + q.Encode()
	}
	response, err := req.SetResult(res).
		SetHeader("Content-Type", "application/json").
		SetHeader("Accept", "application/json").
		SetBody(body).Execute(method, reqURL)
	if err != nil {
		log.Errorf(`failed to connect the captcha service %s: %v`, reqURL, err)
		return nil, fmt.Errorf(`failed to connect the captcha service`)
	}
	if res.Has(`Code`) {
		if res.Int(`Code`) < 1 {
			err = errors.New(res.String(`Info`))
		}
		return res, err
	}
	if response.IsError() {
		log.Errorf(`captcha service exception %s: [%d] %v`, reqURL, response.StatusCode(), com.StripTags(com.Bytes2str(response.Body())))
		return nil, fmt.Errorf(`captcha service exception`)
	}
	return res, err
}
