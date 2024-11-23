package downloader_test

import (
	"net/http"
	"net/http/httptest"
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

	ts := httptest.NewTLSServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			switch r.URL.Path {
			case "/test.json":
				http.ServeFile(w, r, "testdata/test.json")
			case "/test.jpg":
				http.ServeFile(w, r, "testdata/test.jpg")
			default:
				http.NotFound(w, r)
			}
		}))
	defer ts.Close()

	path := t.TempDir()
	client, err := downloader.NewClientDownloader(
		downloader.WithPathToSaveTo(path),
		downloader.WithResourceToDownload([]string{
			ts.URL + "/test.json",
			"https://www.gitlab.com/test.jpg",
		}),
	)
	if err != nil {
		t.Fatal(err)
	}

	client.HttpClient = ts.Client()
	err = client.DownloadFile()
	if err != nil {
		t.Fatal(err)
	}
}

// func TestDownloadFile_DownloadsFileWithoutErrors(t *testing.T) {
// 	url := "http://www.golang-book.com/public/pdf/gobook.pdf"
// 	err := downloader.DownloadFile(url)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// }
