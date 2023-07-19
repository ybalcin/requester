package requester

import (
	"errors"
	"github.com/ybalcin/requester/pkg/utility"
	"net/url"
)

type Request struct {
	URL    *url.URL
	Method string
	Body   string
}

func NewRequest(address, method, body string) (*Request, error) {
	if utility.IsStrEmpty(address) {
		return nil, errors.New("requester address cannot be empty")
	}
	if utility.IsStrEmpty(method) {
		return nil, errors.New("requester Method cannot be empty")
	}

	u, err := url.Parse(address)
	if err != nil {
		return nil, err
	}

	return &Request{
		URL:    u,
		Method: method,
		Body:   body,
	}, nil
}
