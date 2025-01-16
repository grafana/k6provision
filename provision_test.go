package k6provision_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/grafana/k6deps"
	"github.com/grafana/k6provision"
	"github.com/stretchr/testify/require"
)

func TestProvision(t *testing.T) {
	t.Parallel()

	srv := testWebServer(t)
	defer srv.Close()

	u, err := url.Parse(srv.URL + "/minimal-catalog.json")
	require.NoError(t, err)

	ctx := context.Background()

	tmpdir := t.TempDir()

	opts := &k6provision.Options{CacheDir: tmpdir, ExtensionCatalogURL: u}

	exe := filepath.Join(tmpdir, k6provision.ExeName)

	deps := k6deps.Dependencies{}

	err = deps.UnmarshalText([]byte("k6=v0.56.0;k6/x/faker=0.4.1"))
	require.NoError(t, err)

	err = k6provision.Provision(ctx, deps, exe, opts)
	require.NoError(t, err)

	cmd := exec.Command(exe, "version") //nolint:gosec

	contents, err := cmd.CombinedOutput()

	require.NoError(t, err)
	require.Contains(t, string(contents), "k6/x/faker")
}

func TestProvision_errors(t *testing.T) {
	t.Parallel()

	srv := testWebServer(t)
	defer srv.Close()

	u, err := url.Parse(srv.URL + "/missing-catalog.json")
	require.NoError(t, err)

	ctx := context.Background()

	tmpdir := t.TempDir()

	exe := filepath.Join(tmpdir, k6provision.ExeName)

	deps := k6deps.Dependencies{}

	err = deps.UnmarshalText([]byte("k6=v0.56.0;k6/x/faker=0.4.1"))
	require.NoError(t, err)

	opts := &k6provision.Options{CacheDir: tmpdir, ExtensionCatalogURL: u}

	err = k6provision.Provision(ctx, deps, exe, opts)
	require.Error(t, err)
	require.ErrorIs(t, err, k6provision.ErrBuild)

	opts = &k6provision.Options{AppName: invalidAppName(t)}

	err = k6provision.Provision(ctx, deps, exe, opts)
	require.Error(t, err)
	require.ErrorIs(t, err, k6provision.ErrCache)
}

func testWebServer(t *testing.T) *httptest.Server {
	t.Helper()

	return httptest.NewServer(http.FileServer(http.Dir("testdata")))
}

func invalidAppName(t *testing.T) string {
	t.Helper()

	return strings.Repeat("too long", 2048)
}
