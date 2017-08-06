package main

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"path"
	"runtime"

	"github.com/pkg/errors"
)

type Client struct {
	URL        *url.URL
	HTTPClient *http.Client
	AuthToken  string
	Logger     *log.Logger
}

func NewClient(urlStr, authToken string, logger *log.Logger) (*Client, error) {
	discardLogger := log.New(ioutil.Discard, "", log.LstdFlags)
	if logger == nil {
		logger = discardLogger
	}

	parsedURL, err := url.ParseRequestURI(urlStr)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to parse url: %s", urlStr)
	}

	httpClient := &http.Client{}

	client := Client{
		URL:        parsedURL,
		HTTPClient: httpClient,
		AuthToken:  authToken,
		Logger:     logger,
	}
	return &client, nil
}

var userAgent = fmt.Sprintf("gh-webhook-rd (%s)", runtime.Version())

func (c *Client) newRequest(ctx context.Context, method, spath string, body io.Reader) (*http.Request, error) {
	u := *c.URL
	u.Path = path.Join(c.URL.Path, spath)

	req, err := http.NewRequest(method, u.String(), body)
	if err != nil {
		return nil, err
	}

	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("X-Rundeck-Auth-Token", c.AuthToken)

	return req, nil
}

func (c *Client) RunJob(ctx context.Context, jobID string) (bool, error) {
	spath := fmt.Sprintf("/api/1/job/%s/run", jobID)
	req, err := c.newRequest(ctx, "POST", spath, nil)
	if err != nil {
		return false, err
	}

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return false, err
	}

	if res.StatusCode != http.StatusOK {
		return false, errors.New(fmt.Sprintf("unexpected status code: %d", res.StatusCode))
	}

	c.Logger.Println("Succeeded")

	defer res.Body.Close()

	return true, nil
}
