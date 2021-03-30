// +build e2e

package test

import (
	"fmt"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
)

func TestGetUsers(t *testing.T) {
	client := resty.New()
	resp, err := client.R().Get(ROOT_URL + "api/user")
	if err != nil {
		t.Fail()
	}
	assert.Equal(t, 200, resp.StatusCode())
}

func TestCreateUser(t *testing.T) {
	client := resty.New()
	resp, err := client.R().
		SetBody(`{"FirstName": "TestyUser", "LastName": "UserTesty", "Username": "testyguy",
				 "Password": "testyguy", "Email": "testyguy@example.com",
				 "Telephone": "5555555555"}`).
		Post(ROOT_URL + "api/user")
	fmt.Println(resp.String())
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode())
}

func TestValidateUser(t *testing.T) {
	client := resty.New()
	resp, err := client.R().
		SetBody(`{"FirstName": "TestyUser", "LastName": "UserTesty", "Username": "testyguy",
				 "Password": "testyguy", "Email": "testyguy@",
				 "Telephone": "5555555555"}`).
		Post(ROOT_URL + "api/user")

	assert.Error(t, err)
	assert.Equal(t, 400, resp.StatusCode())
}
