package k6provision

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"syscall"

	"github.com/adrg/xdg"
	"github.com/gregjones/httpcache"
	"github.com/gregjones/httpcache/diskcache"
)

// Options contains the optional parameters of the Command function.
type Options struct {
	// AppName contains the name of the application. It is used to define the default value of CacheDir.
	// If empty, it defaults to os.Args[0].
	AppName string
	// CacheDir specifies the name of the directory where the cacheable files can be cached.
	// Its default is determined based on the XDG Base Directory Specification.
	// https://specifications.freedesktop.org/basedir-spec/basedir-spec-latest.html
	CacheDir string
	// Client is used during HTTP communication with the build service.
	// If absent, http.DefaultClient will be used.
	Client *http.Client
	// ExtensionCatalogURL contains the URL of the k6 extension catalog to be used.
	// If absent, DefaultExtensionCatalogURL will be used.
	ExtensionCatalogURL *url.URL
	// BuildServiceURL contains the URL of the k6 build service to be used.
	// If the value is not nil, the k6 binary is built using the build service instead of the local build.
	BuildServiceURL *url.URL
}

// DefaultExtensionCatalogURL contains the address of the default k6 extension catalog.
const DefaultExtensionCatalogURL = "https://registry.k6.io/catalog.json"

func (o *Options) appname() string {
	if o != nil && len(o.AppName) > 0 {
		return o.AppName
	}

	return filepath.Base(os.Args[0]) //nolint:forbidigo
}

func (o *Options) client() (*http.Client, error) {
	if o != nil && o.Client != nil {
		return o.Client, nil
	}

	cachedir, err := o.cacheDir()
	if err != nil {
		return nil, err
	}

	dir := filepath.Join(cachedir, "http")

	err = os.MkdirAll(dir, syscall.S_IRUSR|syscall.S_IWUSR|syscall.S_IXUSR) //nolint:forbidigo
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrCache, err.Error())
	}

	transport := httpcache.NewTransport(diskcache.New(dir))

	return &http.Client{Transport: transport}, nil
}

func (o *Options) extensionCatalogURL() *url.URL {
	if o != nil && o.ExtensionCatalogURL != nil {
		return o.ExtensionCatalogURL
	}

	loc, _ := url.Parse(DefaultExtensionCatalogURL)

	return loc
}

func (o *Options) cacheDir() (string, error) {
	var xdgdir string
	if o != nil {
		xdgdir = o.CacheDir
	}

	if len(xdgdir) == 0 {
		dir, err := xdg.CacheFile(o.appname())
		if err != nil {
			return "", fmt.Errorf("%w: %s", ErrCache, err.Error())
		}

		xdgdir = dir
	}

	err := os.MkdirAll(xdgdir, syscall.S_IRUSR|syscall.S_IWUSR|syscall.S_IXUSR) //nolint:forbidigo
	if err != nil {
		return "", fmt.Errorf("%w: %s", ErrCache, err.Error())
	}

	return xdgdir, nil
}
