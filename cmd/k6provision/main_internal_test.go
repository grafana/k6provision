package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/grafana/k6provision"
	"github.com/stretchr/testify/require"
)

//nolint:forbidigo
func Test_main(t *testing.T) {
	dir := t.TempDir()

	t.Setenv("K6_DEPENDENCIES", "k6/x/faker>0.4.0")

	exe := filepath.Join(dir, k6provision.ExeName)

	os.Args = []string{appname, "--ignore-manifest", "-o", exe}
	main()

	cmd := exec.Command(exe, "version") //nolint:gosec

	contents, err := cmd.CombinedOutput()
	require.NoError(t, err)
	require.Contains(t, string(contents), "k6/x/faker")
}
