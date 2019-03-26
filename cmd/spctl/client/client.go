//Copyright (c) 2018, Oracle and/or its affiliates. All rights reserved.
//Licensed under the Universal Permissive License (UPL) Version 1.0 as shown at http://oss.oracle.com/licenses/upl.

package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"strings"

	"github.com/oracle/speedle/pkg/httputils"
)

type Client struct {
	PMSEndpoint string
	HTTPClient  *http.Client
}

func setAuthorizationHeader(req *http.Request, token string) {
	req.Header.Set("Authorization", "Bearer "+token)

}

func (c *Client) delete(u *url.URL, token string) error {
	req, err := http.NewRequest("DELETE", u.String(), nil)
	if err != nil {
		return err
	}

	setAuthorizationHeader(req, token)
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		fmt.Printf("Error happens: %v\n", err)
		return err
	}
	switch resp.StatusCode {
	case http.StatusNoContent:
		return nil
	case http.StatusUnauthorized, http.StatusForbidden:
		fmt.Println("Authentication or authorization failed. Please specify correct token using '--token' flag.")
		return errors.New(resp.Status)
	case http.StatusBadRequest:
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err == nil {
			var errorDetail httputils.ErrorResponse
			err = json.Unmarshal(body, &errorDetail)
			if err == nil {
				return errors.New(fmt.Sprintf("%s: %s", resp.Status, errorDetail.Error))
			}
		}
		return errors.New(resp.Status)
	default:
		return errors.New(resp.Status)
	}
}

func (c *Client) Delete(paths []string, token string) error {
	u, err := c.pmsURL(paths)
	if err != nil {
		return err
	}
	return c.delete(u, token)
}

func (c *Client) get(u *url.URL, paths []string, params url.Values, token string) ([]byte, error) {
	if params != nil {
		q := u.Query()
		for name, val := range params {
			q.Set(name, val[0])
		}
		u.RawQuery = q.Encode()
	}

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}

	setAuthorizationHeader(req, token)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		fmt.Printf("resp is : %v\n err is : %v\n", resp, err)
		return nil, err
	}
	switch resp.StatusCode {
	case http.StatusNotFound:
		return nil, fmt.Errorf("%s not found", strings.Join(paths, " "))
	case http.StatusOK:
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		return body, nil
	case http.StatusUnauthorized, http.StatusForbidden:
		fmt.Println("Authentication or authorization failed. Please specify correct token using '--token' flag.")
		return nil, errors.New(resp.Status)
	case http.StatusBadRequest:
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err == nil {
			var errorDetail httputils.ErrorResponse
			err = json.Unmarshal(body, &errorDetail)
			if err == nil {
				return nil, errors.New(fmt.Sprintf("%s: %s", resp.Status, errorDetail.Error))
			}
		}
		return nil, errors.New(resp.Status)
	default:
		return nil, errors.New(resp.Status)
	}
}
func (c *Client) Get(paths []string, params url.Values, token string) ([]byte, error) {
	u, err := c.pmsURL(paths)
	if err != nil {
		return nil, err
	}
	return c.get(u, paths, params, token)
}

func (c *Client) post(u *url.URL, paths []string, payload io.Reader, token string) (string, error) {
	req, err := http.NewRequest("POST", u.String(), payload)
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	setAuthorizationHeader(req, token)
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		fmt.Printf("Error happens: %v\n", err)
		return "", err
	}
	switch resp.StatusCode {
	case http.StatusCreated, http.StatusOK:
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return "", err
		}
		return string(body), nil
	case http.StatusUnauthorized, http.StatusForbidden:
		fmt.Println("Authentication or authorization failed. Please specify correct token using '--token' flag.")
		return "", errors.New(resp.Status)
	case http.StatusBadRequest:
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err == nil {
			var errorDetail httputils.ErrorResponse
			err = json.Unmarshal(body, &errorDetail)
			if err == nil {
				return "", errors.New(fmt.Sprintf("%s: %s", resp.Status, errorDetail.Error))
			}
		}
		return "", errors.New(resp.Status)
	default:
		return "", errors.New(resp.Status)
	}
}

func (c *Client) Post(paths []string, payload io.Reader, token string) (string, error) {
	u, err := c.pmsURL(paths)
	if err != nil {
		return "", err
	}
	return c.post(u, paths, payload, token)
}

func getURL(baseURL string, paths []string) (*url.URL, error) {
	u, err := url.Parse(baseURL)
	if err != nil {
		fmt.Println("Fail to parse URL", err)
		return nil, err
	}
	for i := range paths {
		if len(paths[i]) > 0 {
			u.Path = path.Join(u.Path, paths[i])
		}
	}
	return u, nil
}

func (c *Client) pmsURL(paths []string) (*url.URL, error) {
	return getURL(c.PMSEndpoint, paths)
}
