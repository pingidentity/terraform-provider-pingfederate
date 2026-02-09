// Copyright © 2026 Ping Identity Corporation

package api

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"slices"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const maxRetries = 4

func ExponentialBackOffRetryDelete(retryableCodes []int, f func() (*http.Response, error)) (*http.Response, error) {
	var resp *http.Response
	var err error
	backOffTime := time.Second
	var isRetryable bool

	for i := 0; i < maxRetries; i++ {
		resp, err = f()

		backOffTime, isRetryable = testForRetryable(resp, err, retryableCodes, backOffTime)

		if isRetryable {
			tflog.Info(context.Background(), fmt.Sprintf("Attempt %d failed: %v, backing off by %s.", i+1, err, backOffTime.String()))
			time.Sleep(backOffTime)
			continue
		}

		return resp, err
	}

	tflog.Info(context.Background(), fmt.Sprintf("Request failed after %d attempts", maxRetries))

	return resp, err // output the final error
}

func testForRetryable(r *http.Response, err error, retryableCodes []int, currentBackoff time.Duration) (time.Duration, bool) {
	backoff := currentBackoff

	if r != nil {
		backoff = currentBackoff * 2
		if slices.Contains(retryableCodes, r.StatusCode) {
			log.Printf("HTTP status code %d detected, available for retry", r.StatusCode)
			return backoff, true
		}
	}

	return backoff, false
}
