package tobuilds

import (
	"bytes"
	"embed"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"runtime"
)

type Platform string

const (
	PlatformAny     Platform = "any"
	PlatformLinux   Platform = "linux"
	PlatformWindows Platform = "windows"
)

func (p Platform) isCurrent() bool { return string(p) == runtime.GOOS || p == PlatformAny }

type Ctx struct{ embedFs *embed.FS }

func Init(efs *embed.FS) *Ctx { return &Ctx{efs} }

func (c *Ctx) GetFile(name string) (*os.File, error) {
	var reader io.Reader
	var err error
	if isLocalFile(name) {
		reader, err = c.getEmbededFile(name)
	} else {
		var r io.ReadCloser
		r, err = getFromWeb(name)
		defer r.Close()
		reader = r
	}

	if err != nil {
		return nil, err
	}

	f, err := createTempFile(reader)
	return f, err
}

func (c *Ctx) Run(platform Platform, name string, options ...string) error {
	if !platform.isCurrent() {
		fmt.Println("skipped run for different platform")
		return nil
	}

	f, err := c.GetFile(name)
	if err != nil {
		return err
	}
	return runFile(f, options)
}

func (c *Ctx) RunFile(platform Platform, f *os.File, options ...string) error {
	if !platform.isCurrent() {
		fmt.Println("skipped run for different platform")
		return nil
	}

	return runFile(f, options)
}

func (c *Ctx) NewPlatformRunner(platform Platform) *Runner { return &Runner{platform, c} }

type Runner struct {
	platform Platform
	ctx      *Ctx
}

func (r *Runner) Run(name string, options ...string) error {
	return r.ctx.Run(r.platform, name, options...)
}

func (r *Runner) RunFile(f *os.File, options ...string) error {
	return r.ctx.RunFile(r.platform, f, options...)
}

var tmpDir = func() string {
	dir, err := os.MkdirTemp("", "tobuilds_tmp_dir_*")
	if err != nil {
		panic(err)
	}
	return dir
}()

func End() {
	os.RemoveAll(tmpDir)
}

func createTempFile(src io.Reader) (*os.File, error) {
	out, err := os.CreateTemp(tmpDir, "tobuilds_tmp_*")
	if err != nil {
		return nil, err
	}
	defer out.Close()

	_, err = io.Copy(out, src)
	if err != nil {
		return nil, err
	}

	return out, nil
}

func isLocalFile(value string) bool {
	u, err := url.ParseRequestURI(value)
	return err != nil || u.Scheme == "" || u.Host == ""
}

func getFromWeb(u string) (io.ReadCloser, error) {
	resp, err := http.Get(u)
	if err != nil {
		return nil, err
	}
	return resp.Body, nil
}

func (c *Ctx) getEmbededFile(file string) (io.Reader, error) {
	data, err := c.embedFs.ReadFile(file)
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(data), nil
}

// TODO: implement
func runFile(f *os.File, options []string) error {
	err := makeExecutable(f)
	if err != nil {
		return err
	}

	res, err := exec.Command(f.Name(), options...).CombinedOutput()
	fmt.Print("OUT: ", string(res))
	return err
}

func makeExecutable(f *os.File) error {
	if runtime.GOOS != string(PlatformLinux) {
		return nil
	}
	cmd := exec.Command("chmod", "u+x", f.Name())
	return cmd.Run()
}
