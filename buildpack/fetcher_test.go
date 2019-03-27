package buildpack_test

import (
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/fatih/color"
	"github.com/onsi/gomega/ghttp"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"

	"github.com/buildpack/pack/config"

	"github.com/buildpack/pack/buildpack"
	h "github.com/buildpack/pack/testhelpers"
)

func TestBuildpackFetcher(t *testing.T) {
	h.RequireDocker(t)
	color.NoColor = true
	if runtime.GOOS == "windows" {
		t.Skip("create builder is not implemented on windows")
	}
	spec.Run(t, "BuildpackFetcher", testBuildpackFetcher, spec.Parallel(), spec.Report(report.Terminal{}))
}

func testBuildpackFetcher(t *testing.T, when spec.G, it spec.S) {
	when("#FetchBuildpack", func() {
		var (
			subject *buildpack.Fetcher
		)

		it.Before(func() {
			subject = buildpack.NewFetcher(&config.Config{}, nil)
		})

		it("fetches from a relative directory", func() {
			tmpDir, err := ioutil.TempDir("", "")
			h.AssertNil(t, err)
			defer os.RemoveAll(tmpDir)

			bp := buildpack.Buildpack{
				ID:  "bp.one",
				URI: filepath.Join("testdata", "buildpack"),
			}
			h.AssertNil(t, subject.FetchBuildpack(".", &bp))
			h.AssertNotEq(t, bp.Dir, "")
			h.AssertDirContainsFileWithContents(t, bp.Dir, "bin/detect", "I come from a directory\n")
			h.AssertDirContainsFileWithContents(t, bp.Dir, "bin/build", "I come from a directory\n")
		})

		it("fetches from a relative tgz", func() {
			tmpDir, err := ioutil.TempDir("", "")
			h.AssertNil(t, err)
			defer os.RemoveAll(tmpDir)

			bp := buildpack.Buildpack{
				ID:  "bp.one",
				URI: filepath.Join("testdata", "buildpack.tgz"),
			}
			h.AssertNil(t, subject.FetchBuildpack(".", &bp))
			h.AssertNotEq(t, bp.Dir, "")
			h.AssertDirContainsFileWithContents(t, bp.Dir, "bin/detect", "I come from an archive\n")
			h.AssertDirContainsFileWithContents(t, bp.Dir, "bin/build", "I come from an archive\n")
		})

		it("fetches from an absolute directory", func() {
			tmpDir, err := ioutil.TempDir("", "")
			h.AssertNil(t, err)
			defer os.RemoveAll(tmpDir)

			absPath, err := filepath.Abs(filepath.Join("testdata", "buildpack"))
			h.AssertNil(t, err)

			bp := buildpack.Buildpack{
				ID:  "bp.one",
				URI: absPath,
			}
			h.AssertNil(t, subject.FetchBuildpack(".", &bp))
			h.AssertNotEq(t, bp.Dir, "")
			h.AssertDirContainsFileWithContents(t, bp.Dir, "bin/detect", "I come from a directory\n")
			h.AssertDirContainsFileWithContents(t, bp.Dir, "bin/build", "I come from a directory\n")
		})

		it("fetches from an absolute tgz", func() {
			tmpDir, err := ioutil.TempDir("", "")
			h.AssertNil(t, err)
			defer os.RemoveAll(tmpDir)

			absPath, err := filepath.Abs(filepath.Join("testdata", "buildpack.tgz"))
			h.AssertNil(t, err)

			bp := buildpack.Buildpack{
				ID:  "bp.one",
				URI: absPath,
			}
			h.AssertNil(t, subject.FetchBuildpack(".", &bp))
			h.AssertNotEq(t, bp.Dir, "")
			h.AssertDirContainsFileWithContents(t, bp.Dir, "bin/detect", "I come from an archive\n")
			h.AssertDirContainsFileWithContents(t, bp.Dir, "bin/build", "I come from an archive\n")
		})

		it("fetches from a 'file://' URI directory", func() {
			tmpDir, err := ioutil.TempDir("", "")
			h.AssertNil(t, err)
			defer os.RemoveAll(tmpDir)

			absPath, err := filepath.Abs(filepath.Join("testdata", "buildpack"))
			h.AssertNil(t, err)

			bp := buildpack.Buildpack{
				ID:  "bp.one",
				URI: "file://" + absPath,
			}
			h.AssertNil(t, subject.FetchBuildpack(".", &bp))
			h.AssertNotEq(t, bp.Dir, "")
			h.AssertDirContainsFileWithContents(t, bp.Dir, "bin/detect", "I come from a directory\n")
			h.AssertDirContainsFileWithContents(t, bp.Dir, "bin/build", "I come from a directory\n")
		})

		it("fetches from a 'file://' URI tgz", func() {
			tmpDir, err := ioutil.TempDir("", "")
			h.AssertNil(t, err)
			defer os.RemoveAll(tmpDir)

			absPath, err := filepath.Abs(filepath.Join("testdata", "buildpack.tgz"))
			h.AssertNil(t, err)

			bp := buildpack.Buildpack{
				ID:  "bp.one",
				URI: "file://" + absPath,
			}
			h.AssertNil(t, subject.FetchBuildpack(".", &bp))
			h.AssertNotEq(t, bp.Dir, "")
			h.AssertDirContainsFileWithContents(t, bp.Dir, "bin/detect", "I come from an archive\n")
			h.AssertDirContainsFileWithContents(t, bp.Dir, "bin/build", "I come from an archive\n")
		})

		it("fetches from a 'http(s)://' URI tgz", func() {
			server := ghttp.NewServer()
			server.AppendHandlers(func(w http.ResponseWriter, r *http.Request) {
				path := filepath.Join("testdata")
				http.ServeFile(w, r, path)
			})

			tmpDir, err := ioutil.TempDir("", "")
			h.AssertNil(t, err)
			defer os.RemoveAll(tmpDir)

			bp := buildpack.Buildpack{
				ID:  "bp.one",
				URI: "http://buildpack.tgz",
			}
			h.AssertNil(t, subject.FetchBuildpack(".", &bp))
			h.AssertNotEq(t, bp.Dir, "")
			h.AssertDirContainsFileWithContents(t, bp.Dir, "bin/detect", "I come from an archive\n")
			h.AssertDirContainsFileWithContents(t, bp.Dir, "bin/build", "I come from an archive\n")
		})
	})
}
