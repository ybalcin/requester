package requester

import (
	"errors"
	"fmt"
	"github.com/ybalcin/requester/pkg/utility"
	"net/url"
	"strings"
)

type Request struct {
	URL    *url.URL
	Method string
	Body   string
}

func NewRequest(address, method, body string) (*Request, error) {
	address = strings.TrimSpace(address)

	if utility.IsStrEmpty(address) {
		return nil, errors.New("requester address cannot be empty")
	}
	if utility.IsStrEmpty(method) {
		return nil, errors.New("requester method cannot be empty")
	}
	if !utility.IsSchemeExistInURL(address) {
		address = fmt.Sprintf("http://%s", address)
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
