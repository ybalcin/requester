package requester

import (
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"sync"
	"testing"
	"time"
)

// external assertion package not used because of external packages usage restriction

func TestWithWorkerCount(t *testing.T) {
	reqter := New(WithWorkerCount(3))

	if reqter.workerCount != 3 {
		t.Errorf("worker count expected: %d, but got: %d", 3, reqter.workerCount)
	}
}

func TestWithHandleRespBody(t *testing.T) {
	reqter := New(WithHandleSuccess(func(req Request, body []byte) {}))

	if reqter.handleSuccess == nil {
		t.Errorf("handleSuccess expected to set, but nil")
	}
}

func TestWithHandleFailure(t *testing.T) {
	reqter := New(WithHandleFailure(func(err error) {
	}))

	if reqter.handleFailure == nil {
		t.Errorf("handleFailure expected to set, but nil")
	}
}

func TestRequester_Do(t *testing.T) {
	t.Run("handle GET request", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			writer.WriteHeader(http.StatusOK)
			writer.Write([]byte("GET OK"))
		}))
		defer ts.Close()

		parsedURL, _ := url.Parse(ts.URL)
		req := Request{
			URL:    parsedURL,
			Method: http.MethodGet,
		}

		reqter := New(WithHandleSuccess(func(req Request, body []byte) {
			b := string(body)
			if b != "GET OK" {
				t.Errorf("response body expected: 'GET OK', but got: %s", b)
			}
		}))

		reqter.Do(req)
		reqter.Wait()
	})

	t.Run("handle POST request", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			body, _ := io.ReadAll(request.Body)
			writer.WriteHeader(http.StatusOK)
			writer.Write([]byte("POST OK: " + string(body)))
		}))
		defer ts.Close()

		parsedURL, _ := url.Parse(ts.URL)
		req := Request{
			URL:    parsedURL,
			Method: http.MethodPost,
			Body:   "HELLO",
		}

		reqter := New(WithHandleSuccess(func(req Request, body []byte) {
			b := string(body)
			if b != "POST OK: HELLO" {
				t.Errorf("response body expected: 'POST OK: HELLO', but got: %s", b)
			}
		}))

		reqter.Do(req)
		reqter.Wait()
	})

	t.Run("handle non-OK response", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			writer.WriteHeader(http.StatusServiceUnavailable)
		}))
		defer ts.Close()

		parsedURL, _ := url.Parse(ts.URL)
		req := Request{
			URL:    parsedURL,
			Method: http.MethodGet,
			Body:   "",
		}

		reqter := New(WithHandleSuccess(func(req Request, body []byte) {
			t.Errorf("response status expected non-OK, but acted as OK")
		}))

		reqter.Do(req)
		reqter.Wait()
	})

	t.Run("verify worker count", func(t *testing.T) {
		var mu sync.Mutex
		activeWorkers := 0
		maxActiveWorkers := 0

		ts := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			mu.Lock()
			activeWorkers++
			if activeWorkers > maxActiveWorkers {
				maxActiveWorkers = activeWorkers
			}
			mu.Unlock()

			time.Sleep(100 * time.Millisecond)

			mu.Lock()
			activeWorkers--
			mu.Unlock()

			writer.WriteHeader(http.StatusOK)
		}))
		defer ts.Close()

		parsedURL, _ := url.Parse(ts.URL)

		reqter := New(WithWorkerCount(3))

		reqs := make([]Request, 10)
		for i := 0; i < 10; i++ {
			reqs[i] = Request{
				URL:    parsedURL,
				Method: http.MethodGet,
			}
		}

		reqter.Do(reqs...)
		reqter.Wait()

		if maxActiveWorkers != 3 {
			t.Errorf("worker count expected 3, but got: %d", maxActiveWorkers)
		}
	})
}

func TestNewRequest(t *testing.T) {
	t.Run("should return error if address is empty or whitespace", func(t *testing.T) {
		_, err := NewRequest("", http.MethodGet, "")
		if err == nil || err.Error() != "requester address cannot be empty" {
			t.Errorf("expected error: %s, but got: %s", "requester address cannot be empty", err.Error())
		}

		_, err = NewRequest(" ", http.MethodGet, "")
		if err == nil || err.Error() != "requester address cannot be empty" {
			t.Errorf("expected error: %s, but got: %s", "requester address cannot be empty", err.Error())
		}
	})

	t.Run("should return error if method is empty", func(t *testing.T) {
		_, err := NewRequest("google.com", "", "")
		if err == nil || err.Error() != "requester method cannot be empty" {
			t.Errorf("expected error: %s, but got: %s", "requester method cannot be empty", err.Error())
		}

		_, err = NewRequest("google.com", " ", "")
		if err == nil || err.Error() != "requester method cannot be empty" {
			t.Errorf("expected error: %s, but got: %s", "requester method cannot be empty", err.Error())
		}
	})

	t.Run("should set scheme if not exist in address", func(t *testing.T) {
		req, err := NewRequest("google.com", http.MethodGet, "")
		if err != nil {
			t.Errorf("error expected to be nil, but got: %s", err.Error())
		}

		if req == nil || req.URL.String() != "http://google.com" {
			t.Errorf("expected address is http://google.com, but got: %s", req.URL.String())
		}
	})
}
