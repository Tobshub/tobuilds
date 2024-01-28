package tobuilds

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
)

type ArchiveTar struct {
	ctx *Ctx
	Platform
	Name string
	tr   *tar.Reader
}

func newArchiveTar(ctx *Ctx, platform Platform, name string) (*ArchiveTar, error) {
	r, err := ctx.fileReadCloser(name)
	if err != nil {
		return nil, err
	}

	gz, err := gzip.NewReader(r)
	if err != nil {
		return nil, err
	}

	archive := ArchiveTar{
		ctx, platform,
		name, tar.NewReader(gz),
	}
	fmt.Println("INFO: extracted new archive", name)
	return &archive, nil
}

func (a *ArchiveTar) Run(name string, options ...string) error {
	for {
		hdr, err := a.tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		if hdr.Name == name {
			f, err := createTempFile(a.tr)
			if err != nil {
				return err
			}
			f.Close()
			fmt.Printf("INFO: run %s from archive\n", name)
			return RunFile(a.Platform, f, options...)
		}
	}
	return nil
}
