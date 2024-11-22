package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

// http://www.golang-book.com/public/pdf/gobook.pdf
const FILE_NAME = "go1.23.3.windows-amd64.msi"

func main() {
	out, err := os.Create(FILE_NAME)
	if err != nil {
		fmt.Printf("%q", err)
		os.Exit(1)
	}
	resp, err := http.Get("https://go.dev/dl/go1.23.3.windows-amd64.msi")
	defer resp.Body.Close()
	if err != nil {
		fmt.Printf("%q", err)
		os.Exit(1)
	}
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		fmt.Printf("%q", err)
		os.Exit(1)
	}
}
