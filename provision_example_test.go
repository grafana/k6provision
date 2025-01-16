package k6provision_test

import (
	"bytes"
	"context"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/grafana/k6deps"
	"github.com/grafana/k6provision"
)

const depsStr = `k6=0.56.0;k6/x/faker=0.4.1;k6/x/sql=1.0.1;k6/x/sql/driver/ramsql=0.1.0`

func ExampleProvision() {
	deps := k6deps.Dependencies{}
	_ = deps.UnmarshalText([]byte(depsStr))

	exe, _ := filepath.Abs(k6provision.ExeName) // k6 or k6.exe

	_ = k6provision.Provision(context.TODO(), deps, exe, nil)

	cmd := exec.Command(exe, "version") //nolint

	out, _ := cmd.CombinedOutput()

	_, _ = os.Stdout.Write(out[bytes.Index(out, []byte("Extensions:")):]) //nolint

	// Output:
	// Extensions:
	//   github.com/grafana/xk6-faker v0.4.1, k6/x/faker [js]
	//   github.com/grafana/xk6-sql v1.0.1, k6/x/sql [js]
	//   github.com/grafana/xk6-sql v1.0.1, k6/x/sql/driver/ramsql [js]
}
