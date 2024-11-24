package downloader

import (
	"errors"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
)

type clientDownloader struct {
	resourceToDownload []string
	HttpClient         *http.Client
	path               string
}

type option func(*clientDownloader) error

const (
	ERROR_INVALID_URL      = "invalid url"
	ERROR_NO_URL_SPECIFIED = "there are no urls to download from"
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

func NewClientDownloader(opts ...option) (*clientDownloader, error) {
	c := &clientDownloader{}
	for _, opt := range opts {
		err := opt(c)
		if err != nil {
			return nil, err
		}
	}
	return c, nil
}

func (c *clientDownloader) DownloadFile() error {
	// Set up the client.
	req, err := http.NewRequest("GET", c.resourceToDownload[0], nil)
	if err != nil {
		return err
	}
	// Set the User-Agent header to mimic a browser
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:91.0) Gecko/20100101 Firefox/91.0")
	// Get the resource.
	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	// should check resp code.
	// Create file to save to.
	fileName := path.Base(c.resourceToDownload[0])
	out, err := os.Create(filepath.Join(c.path, fileName))
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

// func Main() int {
// 	saveLocation := flag.String("d", "", "name of the location/directory to save files to")
// 	flag.Parse()
// 	newClientDownloader(
// 		withFolderToSaveTo(*saveLocation),
// 		withResourceToDownload(flag.Args()),
// 	)
// 	return 0
// }
