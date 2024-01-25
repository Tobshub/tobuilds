package tobuilds

import (
	"bytes"
	"embed"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"runtime"
)

type Platform string

const (
	PlatformAny     Platform = "any"
	PlatformLinux   Platform = "linux"
	PlatformWindows Platform = "windows"
)

type Ctx struct {
	embedFs *embed.FS
}

func Init(efs *embed.FS) *Ctx { return &Ctx{efs} }

var tmpFiles []*os.File = []*os.File{}

func (c *Ctx) Run(platform Platform, value string, options ...string) {
	if platform != PlatformAny && runtime.GOOS != string(platform) {
		fmt.Println("skipped run for different platform")
		return
	}

	var reader io.Reader
	if isLocalFile(value) {
		reader = c.getEmbededFile(value)
	} else {
		r := getFromWeb(value)
		defer r.Close()
		reader = r
	}

	f, err := createTempFile(reader)
	if err != nil {
		panic(err)
	}

	tmpFiles = append(tmpFiles, f)
	runFile(f, options)
}

func End() {
	for _, f := range tmpFiles {
		os.Remove(f.Name())
	}
}

func createTempFile(src io.Reader) (*os.File, error) {
	out, err := os.CreateTemp("", "tobuilds_download_*")
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

func getFromWeb(u string) io.ReadCloser {
	resp, err := http.Get(u)
	if err != nil {
		panic(err)
	}
	return resp.Body
}

func (c *Ctx) getEmbededFile(file string) io.Reader {
	data, err := c.embedFs.ReadFile(file)
	if err != nil {
		panic(err)
	}
	return bytes.NewReader(data)
}

// TODO: implement
// detect file type(?)
// go utility to run based on ext(?)
func runFile(f *os.File, options []string) {
	buf, err := os.ReadFile(f.Name())
	if err != nil {
		panic(err)
	}
	fmt.Println(string(buf), options)
}
