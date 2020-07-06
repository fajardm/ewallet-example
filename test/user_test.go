package main

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestStoreUser(t *testing.T) {
	cases := []struct {
		description    string
		request        string
		expectedStatus string
		expectedCode   int
	}{
		{
			description:  "test with invalid json",
			request:      `{}`,
			expectedCode: 400,
		},
		{
			description:  "test with empty username",
			request:      `{ "username": "", "email": "john@gmail.com", "mobile_phone": "0817384956973", "password": "secret" }`,
			expectedCode: 400,
		},
		{
			description:  "test with empty email",
			request:      `{ "username": "john", "email": "", "mobile_phone": "0817384956973", "password": "secret" }`,
			expectedCode: 400,
		},
		{
			description:  "test with invalid email",
			request:      `{ "username": "john", "email": "john", "mobile_phone": "0817384956973", "password": "secret" }`,
			expectedCode: 400,
		},
		{
			description:  "test with empty mobile phone",
			request:      `{ "username": "john", "email": "john@gmail.com", "mobile_phone": "", "password": "secret" }`,
			expectedCode: 400,
		},
		{
			description:  "test with empty password",
			request:      `{ "username": "john", "email": "john@gmail.com", "mobile_phone": "", "password": "" }`,
			expectedCode: 400,
		},
		{
			description:  "test with valid json",
			request:      `{ "username": "john", "email": "john@gmail.com", "mobile_phone": "0817384956973", "password": "secret" }`,
			expectedCode: 201,
		},
	}

	for _, test := range cases {
		req, _ := http.NewRequest("POST", "/api/users", bytes.NewBufferString(test.request))
		req.Header.Add("Content-Type", "application/json")
		res, err := app.Test(req, -1)

		assert.NoError(t, err, test.description)
		assert.Equal(t, test.expectedCode, res.StatusCode, test.description)
	}
}
