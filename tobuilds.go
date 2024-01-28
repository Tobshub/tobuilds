package tobuilds

import (
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

// Get a resource from the local embeded filesystem or the web
//
// Returns a closed temporary file with the contents of the resource
func (c *Ctx) GetFile(name string) (*os.File, error) {
	r, err := c.fileReadCloser(name)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	f, err := createTempFile(r)
	if err != nil {
		return nil, err
	}
	f.Close()
	return f, nil
}

func (c *Ctx) fileReadCloser(name string) (r io.ReadCloser, err error) {
	isLocal, location := isLocalFile(name), ""
	if isLocal {
		location = "local"
		r, err = c.getEmbededFile(name)
	} else {
		location = "web"
		r, err = getFromWeb(name)
	}

	if err != nil {
		return nil, err
	}

	fmt.Printf("INFO: read %s (%s)\n", name, location)
	return r, err
}

func (c *Ctx) Run(platform Platform, name string, options ...string) error {
	if !platform.isCurrent() {
		fmt.Printf("INFO: skipped run (%s %s) for different platform\n", name, options)
		return nil
	}

	f, err := c.GetFile(name)
	if err != nil {
		return err
	}

	fmt.Println("RUN:", name, options)
	return runFile(f, options)
}

// NOTE: The file must be closed
func RunFile(platform Platform, f *os.File, options ...string) error {
	if !platform.isCurrent() {
		fmt.Printf("INFO: skipped run (%s %s) for different platform\n", f.Name(), options)
		return nil
	}

	fmt.Println("RUN:", f.Name(), options)
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
	return RunFile(r.platform, f, options...)
}

var tmpDir = func() string {
	dir, err := os.MkdirTemp("", "tobuilds_tmp_dir_*")
	if err != nil {
		panic(err)
	}
	return dir
}()

func End() {
	fmt.Println("INFO: cleaning up", tmpDir)
	os.RemoveAll(tmpDir)
}

func createTempFile(src io.Reader) (*os.File, error) {
	out, err := os.CreateTemp(tmpDir, "tobuilds_tmp_*")
	if err != nil {
		return nil, err
	}

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
	fmt.Println("INFO: start download", u)
	resp, err := http.Get(u)
	if err != nil {
		return nil, err
	}
	return resp.Body, nil
}

func (c *Ctx) getEmbededFile(name string) (io.ReadCloser, error) {
	f, err := c.embedFs.Open(name)
	if err != nil {
		return nil, err
	}
	return f, nil
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

type Archive interface {
	Run(string, ...string) error
	List() ([]string, error)
}

func (c *Ctx) NewArchiveTar(platform Platform, name string) (*ArchiveTar, error) {
	if !platform.isCurrent() {
		fmt.Printf("INFO: skipped extract (%s) for different platform\n", name)
		return nil, nil
	}

	return newArchiveTar(c, platform, name)
}
