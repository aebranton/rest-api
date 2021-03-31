// +build e2e

package test

import (
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
)

// TestGetUsers - tests the get all users endpoint
func TestGetUsers(t *testing.T) {
	client := resty.New()
	resp, err := client.R().Get(ROOT_URL + "api/user")
	if err != nil {
		t.Fail()
	}
	assert.Equal(t, 200, resp.StatusCode())
}

// TestCreateUser - tests creating a user, this should succeed, assuming we have started the test
// containers fresh, and the database is empty (and not trying to write the same username/email again)
// If you get an error because it is not unique, please run:
// docker-compose -f docker-compose.test.yml down
// docker-compose -f docker-compose.test.yml up --remove-orphans
func TestCreateUser(t *testing.T) {
	client := resty.New()
	resp, err := client.R().
		SetBody(`{"FirstName": "TestyUser", "LastName": "UserTesty", "Username": "testyguy",
				 "Password": "testyguy", "Email": "testyguy@example.com",
				 "Telephone": "5555555555"}`).
		Post(ROOT_URL + "api/user")
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode())
}

// TestCreateUserRejectBadEmail - Make sure that creating a user has email validation.
// In a production app id test this with various bad emails, but this gets the point accross!
func TestCreateUserRejectBadEmail(t *testing.T) {
	client := resty.New()
	resp, err := client.R().
		SetBody(`{"FirstName": "TestyUser", "LastName": "UserTesty", "Username": "testyguy2",
				 "Password": "testyguy", "Email": "testyguy@",
				 "Telephone": "5555555555"}`).
		Post(ROOT_URL + "api/user")
	assert.NoError(t, err)
	assert.Equal(t, 400, resp.StatusCode())
}

// TestCreateUserRejectBadPhone - Make sure that creating a user has phone validation.
// In a production app id test this with various bad phone numbers, but this gets the point accross!
func TestCreateUserRejectBadPhone(t *testing.T) {
	client := resty.New()
	resp, err := client.R().
		SetBody(`{"FirstName": "TestyUser", "LastName": "UserTesty", "Username": "testyguy2",
				 "Password": "testyguy", "Email": "testyguy2@example.ca",
				 "Telephone": "seven"}`).
		Post(ROOT_URL + "api/user")
	assert.NoError(t, err)
	assert.Equal(t, 400, resp.StatusCode())
}

// TestCreateUserRejectEmptyField - Make sure that creating a user has required field validation.
// In a production app id test this with submitting each field one at a time as empty, but again, this works
func TestCreateUserRejectEmptyField(t *testing.T) {
	client := resty.New()
	resp, err := client.R().
		SetBody(`{"LastName": "UserTesty", "Username": "testyguy3",
				 "Password": "testyguy", "Email": "testyguy3@example.ca",
				 "Telephone": "seven"}`).
		Post(ROOT_URL + "api/user")
	assert.NoError(t, err)
	assert.Equal(t, 400, resp.StatusCode())
}

// TestUpdateUser - make sure we can do a Put request on the user we created earlier
func TestUpdateUser(t *testing.T) {
	client := resty.New()
	resp, err := client.R().
		SetBody(`{"Telephone": "6666666666"}`).Put(ROOT_URL + "api/user/1")
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode())
}

// TestGetUser - tests getting a single user by ID
func TestGetUser(t *testing.T) {
	client := resty.New()
	resp, err := client.R().
		Get(ROOT_URL + "api/user/1")
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode())
}

// TestGetUserByUsername - Tests the get user by username endpoint
func TestGetUserByUsername(t *testing.T) {
	client := resty.New()
	resp, err := client.R().
		Get(ROOT_URL + "api/user?username=testyguy")
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode())
}

// TestDeleteUser - Tests deleting a user
func TestDeleteUser(t *testing.T) {
	client := resty.New()
	resp, err := client.R().
		Delete(ROOT_URL + "api/user/1")
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode())
}
