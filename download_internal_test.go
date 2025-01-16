package k6provision

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/grafana/k6build/pkg/util"
	"github.com/stretchr/testify/require"
)

//nolint:forbidigo
func Test_httpDownload(t *testing.T) {
	t.Parallel()

	srv := testWebServer(t)
	defer srv.Close()

	tmp := t.TempDir()
	ctx := context.Background()
	from, err := url.Parse(srv.URL + "/empty-catalog.json")

	require.NoError(t, err)

	dest, err := os.Create(filepath.Clean(filepath.Join(tmp, "catalog.json")))

	require.NoError(t, err)

	err = httpDownload(ctx, from, dest, http.DefaultClient)
	require.NoError(t, err)

	require.NoError(t, dest.Close())

	contents, err := os.ReadFile(dest.Name())

	require.NoError(t, err)
	require.Equal(t, "{}", strings.TrimSpace(string(contents)))
}

//nolint:forbidigo
func Test_fileDownload(t *testing.T) {
	t.Parallel()

	srv := testWebServer(t)
	defer srv.Close()

	tmp := t.TempDir()
	abs, err := filepath.Abs(filepath.Join("testdata", "empty-catalog.json"))

	require.NoError(t, err)

	from, err := util.URLFromFilePath(abs)

	require.NoError(t, err)

	dest, err := os.Create(filepath.Clean(filepath.Join(tmp, "catalog.json")))

	require.NoError(t, err)

	err = fileDownload(from, dest)
	require.NoError(t, err)

	require.NoError(t, dest.Close())

	contents, err := os.ReadFile(dest.Name())

	require.NoError(t, err)
	require.Equal(t, "{}", strings.TrimSpace(string(contents)))

	from, err = util.URLFromFilePath(filepath.Join(tmp, "no_such_file"))

	require.NoError(t, err)

	err = fileDownload(from, dest)

	require.Error(t, err)
}

//nolint:forbidigo
func Test_download(t *testing.T) {
	t.Parallel()

	srv := testWebServer(t)
	defer srv.Close()

	tmp := t.TempDir()
	ctx := context.Background()
	from, err := url.Parse(srv.URL + "/empty-catalog.json")

	require.NoError(t, err)

	dest := filepath.Clean(filepath.Join(tmp, "catalog.json"))

	require.NoError(t, err)

	err = download(ctx, from, dest, http.DefaultClient)
	require.NoError(t, err)

	contents, err := os.ReadFile(dest)

	require.NoError(t, err)
	require.Equal(t, "{}", strings.TrimSpace(string(contents)))

	abs, err := filepath.Abs(filepath.Join("testdata", "empty-catalog.json"))

	require.NoError(t, err)

	from, err = util.URLFromFilePath(abs)

	require.NoError(t, err)

	err = download(ctx, from, dest, http.DefaultClient)
	require.NoError(t, err)

	contents, err = os.ReadFile(dest)

	require.NoError(t, err)
	require.Equal(t, "{}", strings.TrimSpace(string(contents)))
}

func testWebServer(t *testing.T) *httptest.Server {
	t.Helper()

	return httptest.NewServer(http.FileServer(http.Dir("testdata")))
}
