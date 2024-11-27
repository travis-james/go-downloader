package downloader

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
)

type clientDownloader struct {
	resourceToDownload []string
	httpClient         *http.Client
	path               string
}

type option func(*clientDownloader) error

const (
	DEFAULT_USER_AGENT     = "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:91.0) Gecko/20100101 Firefox/91.0"
	ERROR_INVALID_URL      = "invalid url"
	ERROR_NO_URL_SPECIFIED = "there are no urls to download from"
	ERROR_STATUS_NOT_OK    = "response did not return 200 ok, received"
)

func isValidURL(stringURL string) bool {
	parsedURL, err := url.Parse(stringURL)
	if err != nil || parsedURL.Host == "" || parsedURL.Scheme == "" {
		return false
	}
	return true
}

func WithResourceToDownload(stringURL []string) option {
	return func(c *clientDownloader) error {
		if len(stringURL) == 0 {
			return errors.New(ERROR_NO_URL_SPECIFIED)
		}
		for _, urp := range stringURL {
			isValid := isValidURL(urp)
			if !isValid {
				return errors.New(ERROR_INVALID_URL)
			}
		}
		c.resourceToDownload = stringURL
		return nil
	}
}

func WithPathToSaveTo(saveLocation string) option {
	return func(c *clientDownloader) error {
		c.path = saveLocation
		return nil
	}
}

func WithHttpClient(client *http.Client) option {
	return func(c *clientDownloader) error {
		c.httpClient = client
		return nil
	}
}

func NewClientDownloader(opts ...option) (*clientDownloader, error) {
	c := &clientDownloader{
		httpClient: http.DefaultClient,
	}
	for _, opt := range opts {
		err := opt(c)
		if err != nil {
			return nil, err
		}
	}
	return c, nil
}

func (c clientDownloader) DownloadFile() error {
	for _, resource := range c.resourceToDownload {
		// Set up the client.
		req, err := http.NewRequest("GET", resource, nil)
		if err != nil {
			return err
		}
		// Set the User-Agent header to mimic a browser
		req.Header.Set("User-Agent", DEFAULT_USER_AGENT)
		// Get the resource.
		resp, err := c.httpClient.Do(req)
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			return errors.New(fmt.Sprintf("%s %d", ERROR_STATUS_NOT_OK, resp.StatusCode))
		}

		// Create the directory if it doesn't exist
		err = os.MkdirAll(c.path, os.ModePerm) // 0777 gives full permissions to everyone...
		if err != nil {
			return err
		}

		// Create file to save to.
		fileName := path.Base(resource)
		out, err := os.Create(filepath.Join(c.path, fileName))
		if err != nil {
			return err
		}
		// Write the response to file.
		_, err = io.Copy(out, resp.Body)
		if err != nil {
			return err
		}
	}
	return nil
}

func Main() int {
	pathToSaveTo := flag.String("path", "", "the name of the path/directory to save the resources to if no argument is supplied it will save to current working directory")
	flag.Usage = func() {
		fmt.Printf("Usage: %s [-path] [URLs...]\n", os.Args[0])
		fmt.Println("Download files at a specified URL(s)\nFlags:")
		flag.PrintDefaults()
	}
	flag.Parse()

	cd, err := NewClientDownloader(
		WithPathToSaveTo(*pathToSaveTo),
		WithResourceToDownload(flag.Args()),
	)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}

	err = cd.DownloadFile()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	return 0
}
