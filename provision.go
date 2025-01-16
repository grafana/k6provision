package k6provision

import (
	"context"
	"errors"
	"fmt"

	"github.com/grafana/k6deps"
)

// ExeName contains the name of the k6 executable. Its value is "k6.exe" for Windows, otherwise "k6".
const ExeName = k6Exe

// Provision provisions a k6 executable with extensions based on dependencies.
func Provision(ctx context.Context, deps k6deps.Dependencies, dest string, opts *Options) error {
	loc, err := build(ctx, deps, opts)
	if err != nil {
		if errors.Is(err, ErrCache) {
			return err
		}

		return fmt.Errorf("%w: %s", ErrBuild, err.Error())
	}

	client, err := opts.client()
	if err != nil {
		return fmt.Errorf("%w: %s", ErrDownload, err.Error())
	}

	if err = download(ctx, loc, dest, client); err != nil {
		return fmt.Errorf("%w: %s", ErrDownload, err.Error())
	}

	return nil
}

var (
	// ErrDownload is returned if an error occurs during download.
	ErrDownload = errors.New("download error")
	// ErrBuild is returned if an error occurs during build.
	ErrBuild = errors.New("build error")
	// ErrCache is returned if an error occurs during cache handling.
	ErrCache = errors.New("cache error")
)
