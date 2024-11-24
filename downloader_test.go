package downloader_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/travis-james/go-downloader"
)

func TestNewClientDownloader_InvalidURL(t *testing.T) {
	t.Parallel()
	_, err := downloader.NewClientDownloader(
		downloader.WithResourceToDownload([]string{"notaurl"}),
	)
	want := downloader.ERROR_INVALID_URL
	got := err.Error()
	if !cmp.Equal(want, got) {
		t.Fatal(cmp.Diff(want, got))
	}
}

func TestNewClientDownloader_NoURLSpecified(t *testing.T) {
	t.Parallel()
	_, err := downloader.NewClientDownloader(
		downloader.WithResourceToDownload([]string{}),
	)
	want := downloader.ERROR_NO_URL_SPECIFIED
	got := err.Error()
	if !cmp.Equal(want, got) {
		t.Fatal(cmp.Diff(want, got))
	}
}

func TestDownloadFile(t *testing.T) {
	t.Parallel()
	// Setup.
	var (
		testJSON   = "test.json"
		testJPG    = "test.jpg"
		testFolder = "testdata"
	)

	ts := httptest.NewTLSServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			switch r.URL.Path {
			case fmt.Sprintf("/%s", testJSON):
				http.ServeFile(w, r, fmt.Sprintf("%s/%s", testFolder, testJSON))
			case fmt.Sprintf("/%s", testJPG):
				http.ServeFile(w, r, fmt.Sprintf("%s/%s", testFolder, testJPG))
			default:
				http.NotFound(w, r)
			}
		}))
	defer ts.Close()

	path := t.TempDir()
	client, err := downloader.NewClientDownloader(
		downloader.WithPathToSaveTo(path),
		downloader.WithResourceToDownload([]string{
			ts.URL + fmt.Sprintf("/%s", testJSON),
			ts.URL + fmt.Sprintf("/%s", testJPG),
		}),
	)
	if err != nil {
		t.Fatal(err)
	}
	client.HttpClient = ts.Client()

	// Run.
	err = client.DownloadFile()
	if err != nil {
		t.Fatal(err)
	}

	// Assert.
	// Check json file.
	want, err := os.ReadFile(filepath.Join(path, testJSON))
	if err != nil {
		t.Fatal(err)
	}
	got, err := os.ReadFile(filepath.Join(path, testJSON))
	if err != nil {
		t.Fatal(err)
	}
	if !cmp.Equal(want, got) {
		t.Fatal(cmp.Diff(want, got))
	}
	// check jpg file.
	want, err = os.ReadFile(filepath.Join(path, testJPG))
	if err != nil {
		t.Fatal(err)
	}
	got, err = os.ReadFile(filepath.Join(path, testJPG))
	if err != nil {
		t.Fatal(err)
	}
	if !cmp.Equal(want, got) {
		t.Fatal(cmp.Diff(want, got))
	}
}

// func TestDownloadFile_DownloadsFileWithoutErrors(t *testing.T) {
// 	url := "http://www.golang-book.com/public/pdf/gobook.pdf"
// 	err := downloader.DownloadFile(url)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// }
