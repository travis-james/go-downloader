package downloader

import (
	"io"
	"net/http"
	"os"
	"path"
)

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
