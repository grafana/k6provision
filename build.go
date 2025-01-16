package k6provision

import (
	"context"
	"log/slog"
	"net/url"
	"os"
	"path/filepath"
	"runtime"

	"github.com/grafana/k6build"
	"github.com/grafana/k6build/pkg/client"
	"github.com/grafana/k6build/pkg/local"
	"github.com/grafana/k6deps"
)

const (
	platform = runtime.GOOS + "/" + runtime.GOARCH

	k6module = "k6"
)

func depsConvert(deps k6deps.Dependencies) (string, []k6build.Dependency) {
	bdeps := make([]k6build.Dependency, 0, len(deps))
	k6constraint := "*"

	for _, dep := range deps {
		if dep.Name == k6module {
			k6constraint = dep.GetConstraints().String()
			continue
		}

		bdeps = append(bdeps, k6build.Dependency{Name: dep.Name, Constraints: dep.GetConstraints().String()})
	}

	return k6constraint, bdeps
}

func build(ctx context.Context, deps k6deps.Dependencies, opts *Options) (*url.URL, error) {
	svc, err := newBuildService(ctx, opts)
	if err != nil {
		return nil, err
	}

	k6constraints, bdeps := depsConvert(deps)

	artifact, err := svc.Build(ctx, platform, k6constraints, bdeps)
	if err != nil {
		return nil, err
	}

	return url.Parse(artifact.URL)
}

func newBuildService(ctx context.Context, opts *Options) (k6build.BuildService, error) {
	if opts != nil && opts.BuildServiceURL != nil {
		return newBuildServiceClient(opts)
	}

	return newLocalBuildService(ctx, opts)
}

func newLocalBuildService(ctx context.Context, opts *Options) (k6build.BuildService, error) {
	client, err := opts.client()
	if err != nil {
		return nil, err
	}

	file, err := os.CreateTemp("", "catalog-*.json") //nolint:forbidigo
	if err != nil {
		return nil, err
	}

	if err := file.Close(); err != nil {
		return nil, err
	}

	catfile := file.Name()

	if err := download(ctx, opts.extensionCatalogURL(), catfile, client); err != nil {
		return nil, err
	}

	cachedir, err := opts.cacheDir()
	if err != nil {
		return nil, err
	}

	conf := local.Config{
		Opts: local.Opts{
			GoOpts: local.GoOpts{
				Env:       map[string]string{"GOWORK": "off"},
				CopyGoEnv: true,
			},
			Verbose: slog.Default().Enabled(ctx, slog.LevelDebug),
		},
		Catalog:  catfile,
		StoreDir: filepath.Join(cachedir, "build"),
	}

	return local.NewBuildService(ctx, conf)
}

func newBuildServiceClient(opts *Options) (k6build.BuildService, error) {
	return client.NewBuildServiceClient(client.BuildServiceClientConfig{
		URL: opts.BuildServiceURL.String(),
	})
}
