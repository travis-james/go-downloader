package downloader

import (
	"io"
	"net/http"
	"os"
	"path"
)

func DownloadFile(stringURL string) error {
	// Get the resource.
	resp, err := http.Get(stringURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
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
