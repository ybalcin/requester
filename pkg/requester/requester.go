package requester

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"sync"
)

const (
	defaultWorkerCount = 10
	defaultQueueSize   = 100
)

type Option func(r *Requester)

type Requester struct {
	workerCount   int
	wg            sync.WaitGroup
	queue         chan Request
	httpClient    *http.Client
	handleSuccess func(req Request, body []byte)
	handleFailure func(err error)
}

// New creates new Requester instance with Option
func New(opts ...Option) *Requester {
	r := &Requester{
		queue:      make(chan Request, defaultQueueSize),
		httpClient: &http.Client{},
	}
	for _, opt := range opts {
		opt(r)
	}
	if r.workerCount <= 0 {
		r.workerCount = defaultWorkerCount
	}

	for i := 0; i < r.workerCount; i++ {
		go r.worker()
	}

	return r
}

// WithWorkerCount option sets parallel worker count
func WithWorkerCount(count int) Option {
	return func(r *Requester) {
		r.workerCount = count
	}
}

// WithHandleSuccess option sets successful response handler
func WithHandleSuccess(handler func(req Request, body []byte)) Option {
	return func(r *Requester) {
		r.handleSuccess = handler
	}
}

// WithHandleFailure option sets failure handler
func WithHandleFailure(handler func(err error)) Option {
	return func(r *Requester) {
		r.handleFailure = handler
	}
}

// Do sends sends requests
func (r *Requester) Do(requests ...Request) {
	for _, req := range requests {
		r.wg.Add(1)
		r.queue <- req
	}
}

func (r *Requester) worker() {
	for req := range r.queue {
		r.sendRequest(req)
	}
}

func (r *Requester) sendRequest(rr Request) {
	defer r.wg.Done()

	req, err := http.NewRequest(rr.Method, rr.URL.String(), bytes.NewBufferString(rr.Body))
	if err != nil {
		r.handleFailureIfHandlerNotNil(fmt.Errorf("requester error during building request: %s, error: %s\n", req.URL.String(), err.Error()))
		return
	}

	resp, err := r.httpClient.Do(req)
	if err != nil {
		r.handleFailureIfHandlerNotNil(fmt.Errorf("requester error occured with request: %s, error: %s\n", req.URL.String(), err.Error()))
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		r.handleFailureIfHandlerNotNil(fmt.Errorf("requester received non-OK HTTP status for request: %s, status: %d\n", req.URL.String(), resp.StatusCode))
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		r.handleFailureIfHandlerNotNil(fmt.Errorf("requester error occured when reading Body with request: %s, error: %s\n", req.URL.String(), err.Error()))
		return
	}

	if r.handleSuccess != nil {
		r.handleSuccess(rr, body)
	}
}

// Wait waits to complete all requests
func (r *Requester) Wait() {
	r.wg.Wait()
}

func (r *Requester) handleFailureIfHandlerNotNil(err error) {
	if r.handleFailure != nil {
		r.handleFailure(err)
	}
}
