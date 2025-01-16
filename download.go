package k6provision

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"syscall"

	"github.com/grafana/k6build/pkg/util"
)

//nolint:forbidigo
func download(ctx context.Context, from *url.URL, dest string, client *http.Client) error {
	tmp, err := os.CreateTemp(filepath.Dir(dest), filepath.Base(dest)+"-*")
	if err != nil {
		return err
	}

	if from.Scheme == "file" {
		err = fileDownload(from, tmp)
	} else {
		err = httpDownload(ctx, from, tmp, client)
	}

	if err != nil {
		_ = tmp.Close()
		_ = os.Remove(tmp.Name())

		return err
	}

	if err = tmp.Close(); err != nil {
		return err
	}

	err = os.Chmod(tmp.Name(), syscall.S_IRUSR|syscall.S_IXUSR|syscall.S_IWUSR)
	if err != nil {
		return err
	}

	return os.Rename(tmp.Name(), dest)
}

//nolint:forbidigo
func fileDownload(from *url.URL, dest *os.File) error {
	filename, err := util.URLToFilePath(from)
	if err != nil {
		return err
	}

	src, err := os.Open(filepath.Clean(filename))
	if err != nil {
		return err
	}

	defer src.Close() //nolint:errcheck

	_, err = io.Copy(dest, src)

	return err
}

//nolint:forbidigo
func httpDownload(ctx context.Context, from *url.URL, dest *os.File, client *http.Client) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, from.String(), nil)
	if err != nil {
		return err
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("%w: %s", os.ErrNotExist, resp.Status)
	}

	defer resp.Body.Close() //nolint:errcheck

	_, err = io.Copy(dest, resp.Body)

	return err
}
