package main

import (
	"github.com/travis-james/go-downloader"
)

const FILE_NAME = "https://file-examples.com/storage/fef4e75e176737761a179bf/2017/10/file_example_JPG_100kB.jpg"

func main() {
	err := downloader.DownloadFile(FILE_NAME)
	if err != nil {
		panic(err)
	}
	// out, err := os.Create(FILE_NAME)
	// if err != nil {
	// 	fmt.Printf("%q", err)
	// 	os.Exit(1)
	// }
	// resp, err := http.Get("https://go.dev/dl/go1.23.3.windows-amd64.msi")
	// if err != nil {
	// 	fmt.Printf("%q", err)
	// 	os.Exit(1)
	// }
	// defer resp.Body.Close()
	// _, err = io.Copy(out, resp.Body)
	// if err != nil {
	// 	fmt.Printf("%q", err)
	// 	os.Exit(1)
	// }
}
