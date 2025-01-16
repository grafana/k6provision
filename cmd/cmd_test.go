package cmd_test

import (
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/grafana/k6provision"
	"github.com/grafana/k6provision/cmd"
	"github.com/stretchr/testify/require"
)

func Test_New(t *testing.T) {
	root := cmd.New()

	require.Equal(t, "k6provision [flags] [script-file]", root.Use)

	t.Setenv("K6_DEPENDENCIES", "k6/x/faker>0.4.0")

	dir := t.TempDir()

	exe := filepath.Join(dir, k6provision.ExeName)

	root.SetArgs([]string{"--ignore-manifest", "-o", exe})

	err := root.Execute()
	require.NoError(t, err)

	cmd := exec.Command(exe, "version") //nolint:gosec

	contents, err := cmd.CombinedOutput()

	require.NoError(t, err)
	require.Contains(t, string(contents), "k6/x/faker")
}
