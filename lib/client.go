package lib

import (
	"bytes"
	"io/ioutil"
	"net/http"

	"github.com/pkg/errors"
)

type client interface {
	request(option *clientOption) ([]byte, error)
}

type clientImpl struct {
	httpClient http.Client
}

type clientOption struct {
	endpoint string
	method   string
	path     string
	header   http.Header
	body     []byte
}

func (c *clientImpl) request(option *clientOption) ([]byte, error) {
	if c == nil {
		return nil, errors.New("client is nil")
	}
	if option == nil {
		return nil, errors.New("option is nil")
	}

	url := option.endpoint + string(option.path)
	req, err := http.NewRequest(option.method, url, bytes.NewReader(option.body))
	if err != nil {
		return nil, errors.WithStack(err)
	}

	req.Header = option.header
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return nil, errors.WithStack(err)
	}

	return body, nil
}
