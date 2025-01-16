[![GitHub Release](https://img.shields.io/github/v/release/grafana/k6provision)](https://github.com/grafana/k6provision/releases/)
[![Go Reference](https://pkg.go.dev/badge/github.com/grafana/k6provision.svg)](https://pkg.go.dev/github.com/grafana/k6provision)
[![Go Report Card](https://goreportcard.com/badge/github.com/grafana/k6provision)](https://goreportcard.com/report/github.com/grafana/k6provision)
[![GitHub Actions](https://github.com/grafana/k6provision/actions/workflows/test.yml/badge.svg)](https://github.com/grafana/k6provision/actions/workflows/test.yml)
[![codecov](https://codecov.io/gh/grafana/k6provision/graph/badge.svg?token=2I3GH2T5YE)](https://codecov.io/gh/grafana/k6provision)
![GitHub Downloads](https://img.shields.io/github/downloads/grafana/k6provision/total)

<h1 name="title">k6provision</h1>

**Provision k6 with extensions**

The purpose of k6provision is to provision k6 executables with extensions based on dependencies.

k6provision is primarily used as a [go library](https://pkg.go.dev/github.com/grafana/k6provision). In addition, it also contains a [command-line tool](#cli), which is suitable for provisioning k6 executable based on the dependencies of k6 test scripts.

The command line tool can be integrated into other command line tools as a subcommand. For this purpose, the library also contains the functionality of the command line tool as a factrory function that returns [cobra.Command](https://pkg.go.dev/github.com/spf13/cobra#Command).

## Install

Precompiled binaries can be downloaded and installed from the [Releases](https://github.com/grafana/k6provision/releases) page.

If you have a go development environment, the installation can also be done with the following command:

```
go install github.com/grafana/k6provision/cmd/k6provision@latest
```

## Usage

Using provisioning is extremely easy. You simply need to pass in the dependencies (name, version constraints pairs) and the filename where you want to save the k6 executable. That's it.

```go file=provision_example_test.go
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
```

## CLI

<!-- #region cli -->
## k6provision

Provision k6 with extensions.

### Synopsis

Analyze the k6 test script and provision k6 with extensions based on dependencies.

### Sources

Dependencies can come from three sources: k6 test script, manifest file, `K6_DEPENDENCIES` environment variable. Instead of these three sources, a k6 archive can also be specified, which can contain all three sources.

### Output

By default, the k6 executable is created in the current directory with the name `k6` (`k6.exe` on Windows). This can be overridden by using the `--o/-output` flag.


```
k6provision [flags] [script-file]
```

### Flags

```
  -o, --output string                  output file (default "k6")
  -e, --env string                     environment variable to analyze (default "K6_DEPENDENCIES")
      --manifest string                manifest file to analyze (default the 'package.json' nearest to the script-file)
      --ingnore-env                    ignore K6_DEPENDENCIES environment variable processing
      --ignore-manifest                disable package.json detection and processing
      --extension-catalog-url string   URL of the k6 extension catalog to be used (default "https://registry.k6.io/catalog.json")
      --build-service-url string       URL of the k6 build service to be used
  -h, --help                           help for k6provision
```

<!-- #endregion cli -->

## Contribute

If you want to contribute or help with the development of **k6provision**, start by 
reading [CONTRIBUTING.md](CONTRIBUTING.md).
