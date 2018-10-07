package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/blitzlog/errors"
	"github.com/blitzlog/proto/api"
	"github.com/blitzlog/proto/common"
)

type Client struct {
	address string
	http    *http.Client
}

func New(address string) *Client {
	return &Client{
		address: address,
		http:    &http.Client{},
	}
}

func (c *Client) CreateToken(extToken, provider string) (
	*api.CreateTokenResponse, error) {

	method, route, status := "POST", "v1/tokens", 200

	createTokenRequest := &api.CreateTokenRequest{
		ExternalToken: extToken,
		Provider:      provider,
	}
	createTokenResponse := new(api.CreateTokenResponse)

	return createTokenResponse, c.do(method, route, "",
		createTokenRequest, createTokenResponse, status)
}

func (c *Client) CreateIdvAccount(token string) (
	*api.CreateAccountResponse, error) {

	method, route, status := "POST", "v1/accounts", 200

	createAccountRequest := &api.CreateAccountRequest{
		Type: "individual",
	}
	createAccountResponse := new(api.CreateAccountResponse)

	return createAccountResponse, c.do(method, route, token,
		createAccountRequest, createAccountResponse, status)
}

func (c *Client) GetAccount(accountId, token string) (
	*api.GetAccountResponse, error) {

	method, route, status := "GET", "v1/accounts/%s", 200

	route = fmt.Sprintf(route, accountId)
	getAccountResponse := new(api.GetAccountResponse)

	return getAccountResponse, c.do(method, route, token,
		nil, getAccountResponse, status)
}

func (c *Client) CreateOrgAccount(accountName, token string) (
	*api.CreateOrgResponse, error) {

	method, route, status := "POST", "v1/accounts", 200

	createOrgRequest := &api.CreateAccountRequest{
		Type: "organization",
		Name: accountName,
	}
	createOrgResponse := new(api.CreateOrgResponse)

	return createOrgResponse, c.do(method, route, token,
		createOrgRequest, createOrgResponse, status)
}

func (c *Client) CreateKey(accountId, token string) (
	*api.CreateKeyResponse, error) {

	method, route, status := "POST", "v1/accounts/%s/keys", 200

	route = fmt.Sprintf(route, accountId)
	createKeyRequest := &api.CreateKeyRequest{}
	createKeyResponse := new(api.CreateKeyResponse)

	return createKeyResponse, c.do(method, route, token,
		createKeyRequest, createKeyResponse, status)
}

func (c *Client) GetKeys(accountId, token string) (
	*api.GetKeysResponse, error) {

	method, route, status := "GET", "v1/accounts/%s/keys", 200

	route = fmt.Sprintf(route, accountId)
	getKeysResponse := new(api.GetKeysResponse)

	return getKeysResponse, c.do(method, route, token,
		nil, getKeysResponse, status)
}

func (c *Client) GetLogs(accountId, token string) (
	*api.GetLogsResponse, error) {

	method, route, status := "POST", "v1/accounts/%s/logs", 200

	// create route
	route = fmt.Sprintf(route, accountId)

	// create get log request
	getLogsRequest := &api.GetLogsRequest{
		Query: &common.LogQuery{
			Page: &common.QueryPage{
				Size: 1000,
			},
		},
	}

	// object to write response to
	getLogsResponse := new(api.GetLogsResponse)

	return getLogsResponse, c.do(method, route, token,
		getLogsRequest, getLogsResponse, status)
}

func (c *Client) do(method, route, token string, in, out interface{},
	expectedStatus int) error {

	// create url
	url := fmt.Sprintf("%s/%s", c.address, route)

	// marshal payload
	payload, err := json.Marshal(in)
	if err != nil {
		return errors.Wrap(err, "error marshaling body")
	}

	// create http request
	req, err := http.NewRequest(method, url, bytes.NewReader(payload))
	if err != nil {
		return errors.Wrap(err, "error creating request")
	}

	// set headers
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Auth-Token", token)

	// make http request
	resp, err := c.http.Do(req)
	if err != nil {
		return errors.Wrap(err, "error making api request")
	}

	// read response body
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.Wrap(err, "error unmarshaling response")
	}

	// check response status
	if resp.StatusCode != expectedStatus {
		return errors.New("unexpected response: %d: %s",
			resp.StatusCode, body)
	}

	// unmarshal response to output
	err = json.Unmarshal([]byte(body), out)
	if err != nil {
		return errors.Wrap(err, "error unmarshalling response")
	}

	return nil
}
