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

func TestDownloadFile_WithValidURLsTheResourceIsSavedToPath(t *testing.T) {
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

	tempPath := t.TempDir()
	client, err := downloader.NewClientDownloader(
		downloader.WithPathToSaveTo(&tempPath),
		downloader.WithResourceToDownload([]string{
			ts.URL + fmt.Sprintf("/%s", testJSON),
			ts.URL + fmt.Sprintf("/%s", testJPG),
		}),
		downloader.WithHttpClient(ts.Client()),
	)
	if err != nil {
		t.Fatal(err)
	}

	// Run.
	err = client.DownloadFileAndSaveFile()
	if err != nil {
		t.Fatal(err)
	}

	// Assert.
	// Check json file.
	want, err := os.ReadFile(filepath.Join(testFolder, testJSON))
	if err != nil {
		t.Fatal(err)
	}
	// Check permissions.
	stat, err := os.Stat(filepath.Join(tempPath, testJSON))
	if err != nil {
		t.Fatal(err)
	}
	perm := stat.Mode().Perm()
	if perm != downloader.PERMISSION_600 {
		t.Errorf("want file mode 0o600, got 0o%o", perm)
	}
	// Check contents.
	got, err := os.ReadFile(filepath.Join(tempPath, testJSON))
	if err != nil {
		t.Fatal(err)
	}
	if !cmp.Equal(want, got) {
		t.Fatal(cmp.Diff(want, got))
	}
	// check jpg file.
	want, err = os.ReadFile(filepath.Join(testFolder, testJPG))
	if err != nil {
		t.Fatal(err)
	}
	// Check permissions.
	stat, err = os.Stat(filepath.Join(tempPath, testJPG))
	if err != nil {
		t.Fatal(err)
	}
	perm = stat.Mode().Perm()
	if perm != downloader.PERMISSION_600 {
		t.Errorf("want file mode 0o600, got 0o%o", perm)
	}
	// Check contents.
	got, err = os.ReadFile(filepath.Join(tempPath, testJPG))
	if err != nil {
		t.Fatal(err)
	}
	if !cmp.Equal(want, got) {
		t.Fatal(cmp.Diff(want, got))
	}
}

// TestMain is special to Go, it will alwasy execute first, purpose
// is to setup any test fixtures.
// In this case, if some test script calls exec downloader, then
// the downloader.Main function should be executed as an independent
// binary, in a sub process, just as if it were a “real” external
// command.
// func TestMain(m *testing.M) {
// 	os.Exit(testscript.RunMain(m, map[string]func() int{
// 		"downloader": downloader.Main,
// 	}))
// }
