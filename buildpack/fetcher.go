package buildpack

import (
	"crypto/sha256"
	"fmt"
	"github.com/buildpack/pack/archive"
	"github.com/buildpack/pack/config"
	"github.com/buildpack/pack/logging"
	"github.com/pkg/errors"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
)

// TODO : test this by itself, currently it is tested in create_builder_test.go
// TODO : attempt to use this during build with the --buildpack flag to get tar.gz buildpacks
// TODO : think of a better name for this construct
// TODO : probably don't need a config here
type Fetcher struct {
	Config *config.Config
	Logger *logging.Logger
}

func NewFetcher(cfg *config.Config, logger *logging.Logger) *Fetcher {
	return &Fetcher{
		Config: cfg,
		Logger: logger,
	}
}

// TODO : pass builder dir (local buildpack search dir) into the constructor ???
// TODO : bp should be builder.BuildpackMetadata and this should fetch buildpack.Metadata
func (f *Fetcher) FetchBuildpack(builderDir string, bp *Buildpack) error {
	asURL, err := url.Parse(bp.URI)
	if err != nil {
		return err
	}

	switch asURL.Scheme {
	case "", "file":
		path := asURL.Path
		if !asURL.IsAbs() && !filepath.IsAbs(path) {
			path = filepath.Join(builderDir, path)
		}

		if filepath.Ext(path) == ".tgz" {
			file, err := os.Open(path)
			if err != nil {
				return errors.Wrapf(err, "could not open file to untar: %q", path)
			}
			defer file.Close()

			tmpDir, err := ioutil.TempDir("", fmt.Sprintf("create-builder-%s-", bp.EscapedID()))
			if err != nil {
				return fmt.Errorf(`failed to create temporary directory: %s`, err)
			}

			if err = archive.ExtractTarGZ(file, tmpDir); err != nil {
				return err
			}
			bp.Dir = tmpDir
		} else {
			bp.Dir = path
		}
	case "http", "https":
		uriDigest := fmt.Sprintf("%x", sha256.Sum256([]byte(bp.URI)))
		cachedDir := filepath.Join(f.Config.Path(), "dl-cache", uriDigest)
		_, err := os.Stat(cachedDir)
		if os.IsNotExist(err) {
			if err = os.MkdirAll(cachedDir, 0744); err != nil {
				return err
			}
		}
		etagFile := cachedDir + ".etag"
		bytes, err := ioutil.ReadFile(etagFile)
		etag := ""
		if err == nil {
			etag = string(bytes)
		}

		reader, etag, err := f.downloadAsStream(bp.URI, etag)
		if err != nil {
			return errors.Wrapf(err, "failed to download from %q", bp.URI)
		} else if reader == nil {
			// can use cached content
			bp.Dir = cachedDir
			break
		}
		defer reader.Close()

		if err = archive.ExtractTarGZ(reader, cachedDir); err != nil {
			return err
		}

		if err = ioutil.WriteFile(etagFile, []byte(etag), 0744); err != nil {
			return err
		}

		bp.Dir = cachedDir
	default:
		return fmt.Errorf("unsupported protocol in URI %q", bp.URI)
	}

	return nil
}

func (f *Fetcher) downloadAsStream(uri string, etag string) (io.ReadCloser, string, error) {
	c := http.Client{}
	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		return nil, "", err
	}
	if etag != "" {
		req.Header.Set("If-None-Match", etag)
	}
	if resp, err := c.Do(req); err != nil {
		return nil, "", err
	} else {
		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			f.Logger.Verbose("Downloading from %q\n", uri)
			return resp.Body, resp.Header.Get("Etag"), nil
		} else if resp.StatusCode == 304 {
			f.Logger.Verbose("Using cached version of %q\n", uri)
			return nil, etag, nil
		} else {
			return nil, "", fmt.Errorf("could not download from %q, code http status %d", uri, resp.StatusCode)
		}
	}
}
