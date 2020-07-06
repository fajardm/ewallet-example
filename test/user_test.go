package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/fajardm/ewallet-example/app/user/model"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"testing"
)

func createUser(request string) model.User {
	req, _ := http.NewRequest("POST", "/api/users", bytes.NewBufferString(request))
	req.Header.Add("Content-Type", "application/json")
	res, err := app.Test(req, -1)
	if err != nil {
		log.Fatal(errors.Wrap(err, "Fatal error create user"))
	}
	body, err := ioutil.ReadAll(res.Body)

	var resp struct {
		Data model.User `json:"data"`
	}
	if err = json.Unmarshal(body, &resp); err != nil {
		log.Fatal(errors.Wrap(err, "Fatal error unmarshal user"))
	}

	return resp.Data
}

func loginUser(request string) string {
	req, _ := http.NewRequest("POST", "/api/users/login", bytes.NewBufferString(request))
	req.Header.Add("Content-Type", "application/json")
	res, err := app.Test(req, -1)
	if err != nil {
		log.Fatal(errors.Wrap(err, "Fatal error login user"))
	}
	body, err := ioutil.ReadAll(res.Body)

	var resp struct {
		Data struct {
			Token string `json:"token"`
		} `json:"data"`
	}
	if err = json.Unmarshal(body, &resp); err != nil {
		log.Fatal(errors.Wrap(err, "Fatal error unmarshal user"))
	}

	return resp.Data.Token
}

func TestLoginUser(t *testing.T) {
	createUser(`{ "username": "cenna", "email": "cenna@gmail.com", "mobile_phone": "081273649506", "password": "secret" }`)

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
			description:  "test with empty username or email",
			request:      `{ "username_or_email": "", "password": "secret" }`,
			expectedCode: 400,
		},
		{
			description:  "test with empty password",
			request:      `{ "username_or_email": "cenna", "password": "" }`,
			expectedCode: 400,
		},
		{
			description:  "test with valid json",
			request:      `{ "username_or_email": "cenna", "password": "secret" }`,
			expectedCode: 200,
		},
	}

	for _, test := range cases {
		req, _ := http.NewRequest("POST", "/api/users/login", bytes.NewBufferString(test.request))
		req.Header.Add("Content-Type", "application/json")
		res, err := app.Test(req, -1)

		assert.NoError(t, err, test.description)
		assert.Equal(t, test.expectedCode, res.StatusCode, test.description)
	}
}

func TestCreateUser(t *testing.T) {
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
			request:      `{ "username": "john", "email": "john@gmail.com", "mobile_phone": "0817384956973", "password": "" }`,
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

func TestGetUser(t *testing.T) {
	user := createUser(`{ "username": "dady", "email": "dady@gmail.com", "mobile_phone": "08172637485", "password": "secret" }`)
	token := loginUser(`{ "username_or_email": "dady", "password": "secret" }`)

	cases := []struct {
		description    string
		id             string
		expectedStatus string
		expectedCode   int
	}{
		{
			description:  "test with invalid id",
			id:           "xxx",
			expectedCode: 400,
		},
		{
			description:  "test with valid id but not registered",
			id:           uuid.NewV4().String(),
			expectedCode: 404,
		},
		{
			description:  "test with valid id",
			id:           user.ID.String(),
			expectedCode: 200,
		},
	}

	for _, test := range cases {
		req, _ := http.NewRequest("GET", fmt.Sprintf("/api/users/%s", test.id), nil)
		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("Authorization", "Bearer "+token)
		res, err := app.Test(req, -1)

		assert.NoError(t, err, test.description)
		assert.Equal(t, test.expectedCode, res.StatusCode, test.description)
	}
}

func TestUpdateUser(t *testing.T) {
	user := createUser(`{ "username": "beny", "email": "beny@gmail.com", "mobile_phone": "081253647589", "password": "secret" }`)
	token := loginUser(`{ "username_or_email": "beny", "password": "secret" }`)

	cases := []struct {
		description    string
		paramID        string
		request        string
		expectedStatus string
		expectedCode   int
	}{
		{
			description:  "test with invalid id",
			paramID:      "xxx",
			request:      `{ "email": "beny2@gmail.com", "password": "secret" }`,
			expectedCode: 400,
		},
		{
			description:  "test with valid id but not registered",
			paramID:      uuid.NewV4().String(),
			request:      `{ "email": "beny2@gmail.com", "password": "secret" }`,
			expectedCode: 404,
		},
		{
			description:  "test with invalid json",
			paramID:      user.ID.String(),
			request:      `{}`,
			expectedCode: 400,
		},
		{
			description:  "test with empty email",
			paramID:      user.ID.String(),
			request:      `{ "email": "", "password": "secret" }`,
			expectedCode: 400,
		},
		{
			description:  "test with invalid email",
			paramID:      user.ID.String(),
			request:      `{ "email": "john", "password": "secret" }`,
			expectedCode: 400,
		},
		{
			description:  "test with empty password",
			paramID:      user.ID.String(),
			request:      `{ "email": "john@gmail.com", "password": "" }`,
			expectedCode: 400,
		},
		{
			description:  "test with valid json",
			paramID:      user.ID.String(),
			request:      `{ "email": "beny2@gmail.com", "password": "secret" }`,
			expectedCode: 200,
		},
	}

	for _, test := range cases {
		req, _ := http.NewRequest("PUT", fmt.Sprintf("/api/users/%s", test.paramID), bytes.NewBufferString(test.request))
		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("Authorization", "Bearer "+token)
		res, err := app.Test(req, -1)

		assert.NoError(t, err, test.description)
		assert.Equal(t, test.expectedCode, res.StatusCode, test.description)
	}
}

func TestDeletetUser(t *testing.T) {
	user := createUser(`{ "username": "dony", "email": "dony@gmail.com", "mobile_phone": "081253840698", "password": "secret" }`)
	token := loginUser(`{ "username_or_email": "dony", "password": "secret" }`)

	cases := []struct {
		description    string
		id             string
		expectedStatus string
		expectedCode   int
	}{
		{
			description:  "test with invalid id",
			id:           "xxx",
			expectedCode: 400,
		},
		{
			description:  "test with valid id but not registered",
			id:           uuid.NewV4().String(),
			expectedCode: 404,
		},
		{
			description:  "test with valid id",
			id:           user.ID.String(),
			expectedCode: 200,
		},
	}

	for _, test := range cases {
		req, _ := http.NewRequest("DELETE", fmt.Sprintf("/api/users/%s", test.id), nil)
		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("Authorization", "Bearer "+token)
		res, err := app.Test(req, -1)

		assert.NoError(t, err, test.description)
		assert.Equal(t, test.expectedCode, res.StatusCode, test.description)
	}
}
