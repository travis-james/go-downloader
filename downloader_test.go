package downloader_test

import (
	"testing"

	"github.com/travis-james/go-downloader"
)

func TestDownloadFile_DownloadsFileWithoutErrors(t *testing.T) {
	url := "http://www.golang-book.com/public/pdf/gobook.pdf"
	err := downloader.DownloadFile(url)
	if err != nil {
		t.Fatal(err)
	}
}
