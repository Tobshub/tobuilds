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

var pendingFiles []*os.File = []*os.File{}

func (c *Ctx) Run(platform Platform, value string, options ...string) {
	if platform != PlatformAny && runtime.GOOS != string(platform) {
		fmt.Println("skipped run for different platform")
		return
	}

	var f *os.File
	if isLocalFile(value) {
		f = c.getEmbededFile(value)
	} else {
		f = getFromWeb(value)
	}
	pendingFiles = append(pendingFiles, f)
	runFile(f, options)
}

// TODO: listen for program end then run
func cleanUpFiles() {
	for _, f := range pendingFiles {
		os.Remove(f.Name())
	}
}

func getFromWeb(u string) *os.File {
	resp, err := http.Get(u)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()
	f, err := createTempFile(resp.Body)
	if err != nil {
		panic(err)
	}
	return f
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

func (c *Ctx) getEmbededFile(file string) *os.File {
	data, err := c.embedFs.ReadFile(file)
	if err != nil {
		panic(err)
	}

	f, err := createTempFile(bytes.NewReader(data))
	if err != nil {
		panic(err)
	}
	return f
}

// TODO: implement
func runFile(f *os.File, options []string) {
	buf, err := os.ReadFile(f.Name())
	if err != nil {
		panic(err)
	}
	fmt.Println(string(buf), options)
}
