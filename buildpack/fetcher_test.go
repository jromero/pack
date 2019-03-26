package buildpack_test

import (
	"github.com/buildpack/pack/buildpack"
	h "github.com/buildpack/pack/testhelpers"
	"github.com/fatih/color"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"testing"
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
			subject = buildpack.NewFetcher(nil, nil)
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

		it.Pend("fetches from an absolute directory", func() {
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

		it.Pend("fetches from an absolute tgz", func() {
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

		it.Pend("fetches from a 'file://' URI directory", func() {
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

		it.Pend("fetches from a 'file://' URI tgz", func() {
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

		it.Pend("fetches from a 'http(s)://' URI tgz", func() {
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

		//		when("a relative directory", func() {
		//
		//			it("supports absolute directories as well as archives", func() {
		//				mockImage := mocks.NewMockImage(mockController)
		//				mockFetcher.EXPECT().FetchLocalImage("some/build").Return(mockImage, nil)
		//				mockImage.EXPECT().Rename("myorg/mybuilder")
		//
		//				absPath, err := filepath.Abs("testdata/used-to-test-various-uri-schemes/buildpack")
		//				h.AssertNil(t, err)
		//
		//				f, err := ioutil.TempFile("", "*.toml")
		//				h.AssertNil(t, err)
		//				ioutil.WriteFile(f.Name(), []byte(fmt.Sprintf(`[[buildpacks]]
		//id = "some.bp.with.no.uri.scheme"
		//uri = "%s"
		//
		//[[buildpacks]]
		//id = "some.bp.with.no.uri.scheme.and.tgz"
		//uri = "%s.tgz"
		//
		//[[groups]]
		//buildpacks = [
		//  { id = "some.bp.with.no.uri.scheme", version = "1.2.3" },
		//  { id = "some.bp.with.no.uri.scheme.and.tgz", version = "1.2.4" },
		//]
		//
		//[[groups]]
		//buildpacks = [
		//  { id = "some.bp1", version = "1.2.3" },
		//]
		//
		//[stack]
		//id = "com.example.stack"
		//build-image = "some/build"
		//run-image = "some/run"
		//`, absPath, absPath)), 0644)
		//				f.Name()
		//
		//				flags := pack.CreateBuilderFlags{
		//					RepoName:        "myorg/mybuilder",
		//					BuilderTomlPath: f.Name(),
		//					Publish:         false,
		//					NoPull:          true,
		//				}
		//
		//				builderConfig, err := factory.BuilderConfigFromFlags(context.TODO(), flags)
		//				h.AssertNil(t, err)
		//
		//				h.AssertDirContainsFileWithContents(t, builderConfig.Buildpacks[0].Dir, "bin/detect", "I come from a directory")
		//				h.AssertDirContainsFileWithContents(t, builderConfig.Buildpacks[1].Dir, "bin/build", "I come from an archive")
		//			})
		//		})
		//
		//		when("a buildpack location uses file:// uris", func() {
		//			it("supports absolute directories as well as archives", func() {
		//				mockImage := mocks.NewMockImage(mockController)
		//				mockFetcher.EXPECT().FetchLocalImage("some/build").Return(mockImage, nil)
		//				mockImage.EXPECT().Rename("myorg/mybuilder")
		//
		//				absPath, err := filepath.Abs("testdata/used-to-test-various-uri-schemes/buildpack")
		//				h.AssertNil(t, err)
		//
		//				f, err := ioutil.TempFile("", "*.toml")
		//				h.AssertNil(t, err)
		//				ioutil.WriteFile(f.Name(), []byte(fmt.Sprintf(`[[buildpacks]]
		//id = "some.bp.with.no.uri.scheme"
		//uri = "file://%s"
		//
		//[[buildpacks]]
		//id = "some.bp.with.no.uri.scheme.and.tgz"
		//uri = "file://%s.tgz"
		//
		//[[groups]]
		//buildpacks = [
		//  { id = "some.bp.with.no.uri.scheme", version = "1.2.3" },
		//  { id = "some.bp.with.no.uri.scheme.and.tgz", version = "1.2.4" },
		//]
		//
		//[[groups]]
		//buildpacks = [
		//  { id = "some.bp1", version = "1.2.3" },
		//]
		//
		//[stack]
		//id = "com.example.stack"
		//build-image = "some/build"
		//run-image = "some/run"
		//`, absPath, absPath)), 0644)
		//				f.Name()
		//
		//				flags := pack.CreateBuilderFlags{
		//					RepoName:        "myorg/mybuilder",
		//					BuilderTomlPath: f.Name(),
		//					Publish:         false,
		//					NoPull:          true,
		//				}
		//
		//				builderConfig, err := factory.BuilderConfigFromFlags(context.TODO(), flags)
		//				h.AssertNil(t, err)
		//
		//				h.AssertDirContainsFileWithContents(t, builderConfig.Buildpacks[0].Dir, "bin/detect", "I come from a directory")
		//				h.AssertDirContainsFileWithContents(t, builderConfig.Buildpacks[1].Dir, "bin/build", "I come from an archive")
		//			})
		//		})
		//
		//		when("a buildpack location uses http(s):// uris", func() {
		//			var (
		//				server *http.Server
		//			)
		//			it.Before(func() {
		//				port := 1024 + rand.Int31n(65536-1024)
		//				fs := http.FileServer(http.Dir("testdata"))
		//				server = &http.Server{Addr: fmt.Sprintf("127.0.0.1:%d", port), Handler: fs}
		//				go func() {
		//					err := server.ListenAndServe()
		//					if err != http.ErrServerClosed {
		//						t.Fatalf("could not create http server: %v", err)
		//					}
		//				}()
		//				serverReady := false
		//				for i := 0; i < 10; i++ {
		//					resp, err := http.Get(fmt.Sprintf("http://%s/used-to-test-various-uri-schemes/buildpack.tgz", server.Addr))
		//					if err == nil && resp.StatusCode >= 200 && resp.StatusCode < 300 {
		//						serverReady = true
		//						break
		//					}
		//					t.Logf("Waiting for server to become ready on %s. Currently %v\n", server.Addr, err)
		//					time.Sleep(1 * time.Second)
		//				}
		//				if !serverReady {
		//					t.Fatal("http server does not seem to be up")
		//				}
		//			})
		//			it("downloads and extracts the archive", func() {
		//				mockImage := mocks.NewMockImage(mockController)
		//				mockFetcher.EXPECT().FetchLocalImage("some/build").Return(mockImage, nil)
		//				mockImage.EXPECT().Rename("myorg/mybuilder")
		//
		//				f, err := ioutil.TempFile("", "*.toml")
		//				h.AssertNil(t, err)
		//				ioutil.WriteFile(f.Name(), []byte(fmt.Sprintf(`[[buildpacks]]
		//id = "some.bp.with.no.uri.scheme"
		//uri = "http://%s/used-to-test-various-uri-schemes/buildpack.tgz"
		//
		//[[groups]]
		//buildpacks = [
		//  { id = "some.bp.with.no.uri.scheme", version = "1.2.3" },
		//]
		//
		//[[groups]]
		//buildpacks = [
		//  { id = "some.bp1", version = "1.2.3" },
		//]
		//
		//[stack]
		//id = "com.example.stack"
		//build-image = "some/build"
		//run-image = "some/run"
		//`, server.Addr)), 0644)
		//				f.Name()
		//
		//				flags := pack.CreateBuilderFlags{
		//					RepoName:        "myorg/mybuilder",
		//					BuilderTomlPath: f.Name(),
		//					Publish:         false,
		//					NoPull:          true,
		//				}
		//
		//				builderConfig, err := factory.BuilderConfigFromFlags(context.TODO(), flags)
		//				h.AssertNil(t, err)
		//
		//				h.AssertDirContainsFileWithContents(t, builderConfig.Buildpacks[0].Dir, "bin/build", "I come from an archive")
		//			})
		//			it.After(func() {
		//				if server != nil {
		//					ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)
		//					server.Shutdown(ctx)
		//				}
		//			})
		//		})
	})
}
