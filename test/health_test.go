// +build e2e

package test

import (
	"fmt"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
)

func TestStatusEndpoint(t *testing.T) {
	fmt.Println("Running e2e test for status check")

	client := resty.New()
	resp, err := client.R().Get("http://localhost:8080/api/status")
	if err != nil {
		t.Fail()
	}

	assert.Equal(t, 200, resp.StatusCode())
}
