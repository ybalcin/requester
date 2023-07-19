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
	workerCount    int
	wg             sync.WaitGroup
	queue          chan Request
	httpClient     *http.Client
	handleRespBody func(req Request, body []byte)
}

// New creates new Requester instance with workerCount argument
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

func WithWorkerCount(count int) Option {
	return func(r *Requester) {
		r.workerCount = count
	}
}

func WithHandleRespBody(handler func(req Request, body []byte)) Option {
	return func(r *Requester) {
		r.handleRespBody = handler
	}
}

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
		fmt.Printf("requester error during building request: %s, error: %s\n", req.URL.String(), err.Error())
		return
	}

	resp, err := r.httpClient.Do(req)
	if err != nil {
		fmt.Printf("requester error occured with request: %s, error: %s\n", req.URL.String(), err.Error())
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("requester received non-OK HTTP status for request: %s, status: %d\n", req.URL.String(), resp.StatusCode)
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("requester error occured when reading Body with request: %s, error: %s\n", req.URL.String(), err.Error())
		return
	}

	if r.handleRespBody != nil {
		r.handleRespBody(rr, body)
	}
}

// Wait waits to complete all requests
func (r *Requester) Wait() {
	r.wg.Wait()
}
