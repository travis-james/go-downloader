package downloader

import (
	"errors"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
)

type Client struct {
	ResourceToDownload string
	HTTPClient         *http.Client
	DstFileName        string
}

type option func(*Client) error

func WithResourceToDownload(stringURL string) option {
	return func(c *Client) error {
		_, err := url.Parse(stringURL)
		if err != nil {
			return errors.New("invalid url")
		}
		c.ResourceToDownload = stringURL
		return nil
	}
}

func NewClient(opts ...option) (*Client, error) {
	c := &Client{}
	for _, opt := range opts {
		err := opt(c)
		if err != nil {
			return nil, err
		}
	}
	return c, nil
}

func DownloadFile(stringURL string) error {
	// Set up the client.
	req, err := http.NewRequest("GET", stringURL, nil)
	if err != nil {
		return err
	}
	// Set the User-Agent header to mimic a browser
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:91.0) Gecko/20100101 Firefox/91.0")
	// Get the resource.
	client := &http.Client{}
	resp, err := client.Do(req)
	//resp, err := http.Get(stringURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	// should check resp code.
	// Create file to save to.
	fileName := path.Base(stringURL)
	out, err := os.Create(fileName)
	if err != nil {
		return err
	}
	// Write the response to file.
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}
	return nil
}
