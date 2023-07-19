package cmd

import (
	"context"
	"crypto/md5"
	"errors"
	"flag"
	"fmt"
	"github.com/ybalcin/requester/pkg/requester"
	"github.com/ybalcin/requester/pkg/utility"
	"net/http"
	"os"
)

func SendRequests(ctx context.Context) error {
	parallel := flag.Int("parallel", 0, "parallel worker count")
	flag.Parse()

	urls := flag.Args()
	if len(urls) <= 0 {
		return errors.New("at least one url is required")
	}

	reqter := requester.New(requester.WithWorkerCount(*parallel), requester.WithHandleRespBody(handleRespBody))
	go func() {
		<-ctx.Done()
		reqter.Wait()
		os.Exit(0)
	}()

	requests := make([]requester.Request, len(urls))
	for i, u := range urls {
		if !utility.IsStrEmpty(u) {
			req, err := requester.NewRequest(u, http.MethodGet, "")
			if err != nil {
				return err
			}
			requests[i] = *req
		}
	}

	reqter.Do(requests...)
	reqter.Wait()

	return nil
}

func handleRespBody(req requester.Request, respBody []byte) {
	hashedBody := md5.Sum(respBody)
	fmt.Printf("%s %x\n", req.URL.String(), hashedBody)
}
