package k6provision

import (
	"context"
	"fmt"
	"net/url"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/grafana/k6build/pkg/client"
	"github.com/grafana/k6build/pkg/util"

	"github.com/grafana/k6deps"
	"github.com/stretchr/testify/require"
)

func Test_depsConvert(t *testing.T) {
	t.Parallel()

	src := make(k6deps.Dependencies)

	err := src.UnmarshalText([]byte("k6>0.50;k6/x/faker>0.2.0"))

	require.NoError(t, err)

	k6Constraints, deps := depsConvert(src)

	require.Equal(t, ">0.50", k6Constraints)

	require.Equal(t, "k6/x/faker", deps[0].Name)
	require.Equal(t, src["k6/x/faker"].Constraints.String(), deps[0].Constraints)

	err = src.UnmarshalText([]byte("k6/x/faker*"))

	require.NoError(t, err)

	k6Constraints, deps = depsConvert(src)

	require.Equal(t, "*", k6Constraints)

	require.Equal(t, "k6/x/faker", deps[0].Name)
	require.Equal(t, "*", deps[0].Constraints)
}

func Test_newBuildService(t *testing.T) {
	t.Parallel()

	opts := &Options{
		CacheDir: t.TempDir(),
	}

	opts.BuildServiceURL, _ = url.Parse("http://localhost:8000")

	svc, err := newBuildService(context.Background(), opts)

	require.NoError(t, err)
	require.IsType(t, new(client.BuildClient), svc)
	require.Equal(t, "*client.BuildClient", fmt.Sprintf("%T", svc))

	srv := testWebServer(t)
	defer srv.Close()

	opts.ExtensionCatalogURL, err = url.Parse(srv.URL + "/empty-catalog.json")
	require.NoError(t, err)
	opts.BuildServiceURL = nil

	svc, err = newBuildService(context.Background(), opts)

	require.NoError(t, err)

	require.NotEqual(t, "*k6build.BuildClient", fmt.Sprintf("%T", svc))
}

func Test_newLocalBuildService(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	abs, err := filepath.Abs(filepath.Join("testdata", "minimal-catalog.json"))
	require.NoError(t, err)

	u, err := util.URLFromFilePath(abs)
	require.NoError(t, err)

	opts := &Options{ExtensionCatalogURL: u}

	_, err = newLocalBuildService(ctx, opts)
	require.NoError(t, err)

	opts.AppName = invalidAppName(t)
	_, err = newLocalBuildService(ctx, opts)
	require.Error(t, err)
}

func Test_build(t *testing.T) {
	t.Parallel()

	srv := testWebServer(t)
	defer srv.Close()

	u, err := url.Parse(srv.URL + "/minimal-catalog.json")

	require.NoError(t, err)

	ctx := context.Background()

	opts := &Options{CacheDir: t.TempDir(), ExtensionCatalogURL: u}

	loc, err := build(ctx, make(k6deps.Dependencies), opts)

	require.NoError(t, err)

	tmp := t.TempDir()

	bin := "k6"
	if runtime.GOOS == "windows" {
		bin += ".exe"
	}

	dest := filepath.Join(tmp, bin)

	err = download(ctx, loc, dest, nil)

	require.NoError(t, err)

	cmd := exec.Command(filepath.Clean(dest), "version") //nolint:gosec

	out, err := cmd.Output()

	require.NoError(t, err)
	require.True(t, strings.HasPrefix(string(out), "k6"))

	opts.ExtensionCatalogURL, err = url.Parse(srv.URL + "/empty-catalog.json")

	require.NoError(t, err)

	_, err = build(ctx, make(k6deps.Dependencies), opts)

	require.Error(t, err)

	opts.ExtensionCatalogURL, err = url.Parse(srv.URL + "/missing-catalog.json")

	require.NoError(t, err)

	_, err = build(ctx, make(k6deps.Dependencies), opts)

	require.Error(t, err)
}
