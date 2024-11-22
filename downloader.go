package downloader

import (
	"errors"
	"flag"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
)

type clientDownloader struct {
	resourceToDownload []string
	httpClient         *http.Client
	dstFolder          string
}

type option func(*clientDownloader) error

func withResourceToDownload(stringURL []string) option {
	return func(c *clientDownloader) error {
		for _, urp := range stringURL {
			_, err := url.Parse(urp)
			if err != nil {
				return errors.New("invalid url")
			}
		}
		c.resourceToDownload = stringURL
		return nil
	}
}

func withFolderToSaveTo(saveLocation string) option {
	return func(c *clientDownloader) error {
		c.dstFolder = saveLocation
		return nil
	}
}

func newClientDownloader(opts ...option) (*clientDownloader, error) {
	c := &clientDownloader{}
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

func Main() int {
	saveLocation := flag.String("d", "", "name of the location/directory to save files to")
	flag.Parse()
	newClientDownloader(
		withFolderToSaveTo(*saveLocation),
		withResourceToDownload(flag.Args()),
	)
	return 0
}
