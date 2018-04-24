package lib

import (
	"io/ioutil"
	"net/http"
	"time"

	"github.com/pkg/errors"
)

type client interface {
	request(option *clientOption) ([]byte, error)
}

type clientImpl struct {
	endpoint       string
	defaultTimeout time.Duration
}

type clientOption struct {
	method  string
	path    string
	timeout time.Duration
}

func (c *clientImpl) request(option *clientOption) ([]byte, error) {
	if c == nil {
		return nil, errors.New("client is nil")
	}
	if option == nil {
		return nil, errors.New("option is nil")
	}

	url := c.endpoint + string(option.path)
	req, err := http.NewRequest(option.method, url, nil)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if option.timeout == 0 {
		option.timeout = c.defaultTimeout
	}

	client := http.Client{
		Timeout: option.timeout,
	}
	res, err := client.Do(req)
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
