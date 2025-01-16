// Package cmd contains k6provision cobra command factory function.
package cmd

import (
	"context"
	_ "embed"
	"net/url"
	"runtime"
	"strings"

	"github.com/grafana/k6deps"
	"github.com/grafana/k6provision"
	"github.com/spf13/cobra"
)

type options struct {
	k6deps.Options
	buildServiceURL     string
	extensionCatalogURL string
	output              string
	env                 string
}

//go:embed help.md
var help string

// New creates new cobra command for provision command.
func New() *cobra.Command {
	opts := new(options)

	opts.output = "k6"

	if runtime.GOOS == "windows" {
		opts.output += ".exe"
	}

	cmd := &cobra.Command{
		Use:   "k6provision [flags] [script-file]",
		Short: "Provision k6 with extensions.",
		Long:  help,
		Args:  cobra.MaximumNArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			return prov(opts, args)
		},
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	flags := cmd.Flags()

	flags.SortFlags = false

	flags.StringVarP(&opts.output, "output", "o", opts.output, "output file")
	flags.StringVarP(&opts.Env.Name, "env", "e", k6deps.EnvDependencies, "environment variable to analyze")
	flags.StringVar(&opts.Manifest.Name, "manifest", "",
		"manifest file to analyze (default the 'package.json' nearest to the script-file)")
	flags.BoolVar(&opts.Env.Ignore, "ingnore-env", false,
		"ignore "+k6deps.EnvDependencies+" environment variable processing")
	flags.BoolVar(&opts.Manifest.Ignore, "ignore-manifest", false, "disable package.json detection and processing")

	flags.StringVar(
		&opts.extensionCatalogURL,
		"extension-catalog-url",
		k6provision.DefaultExtensionCatalogURL,
		"URL of the k6 extension catalog to be used",
	)
	flags.StringVar(
		&opts.buildServiceURL,
		"build-service-url",
		opts.buildServiceURL,
		"URL of the k6 build service to be used",
	)

	return cmd
}

func prov(opts *options, args []string) error {
	if len(args) > 0 {
		opts.Script.Name = args[0]
	}

	if !opts.Script.Ignore && strings.HasSuffix(opts.Script.Name, ".tar") {
		opts.Archive.Name = opts.Script.Name
		opts.Script.Ignore = true
		opts.Env.Ignore = true
		opts.Manifest.Ignore = true
	}

	if len(opts.env) > 0 {
		opts.Env.Contents = []byte(opts.env)
		opts.Env.Name = k6deps.EnvDependencies
	}

	deps, err := k6deps.Analyze(&opts.Options)
	if err != nil {
		return err
	}

	popts := new(k6provision.Options)

	if len(opts.buildServiceURL) > 0 {
		popts.BuildServiceURL, err = url.Parse(opts.buildServiceURL)
		if err != nil {
			return err
		}
	}

	if len(opts.extensionCatalogURL) > 0 {
		popts.ExtensionCatalogURL, err = url.Parse(opts.extensionCatalogURL)
		if err != nil {
			return err
		}
	}

	return k6provision.Provision(context.Background(), deps, opts.output, popts)
}
