package tobuilds

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
)

type ArchiveTarGz struct {
	ctx *Ctx
	Platform
	name string
}

func (a *ArchiveTarGz) newArchiveTarGzReader() (*tar.Reader, error) {
	r, err := a.ctx.fileReadCloser(a.name)
	if err != nil {
		return nil, err
	}

	gz, err := gzip.NewReader(r)
	if err != nil {
		return nil, err
	}

	tr := tar.NewReader(gz)
	return tr, nil
}

func newArchiveTarGz(ctx *Ctx, platform Platform, name string) *ArchiveTarGz {
	fmt.Println("INFO: registered new archive", name)
	return &ArchiveTarGz{ctx, platform, name}
}

func (a *ArchiveTarGz) Run(name string, options ...string) error {
	tr, err := a.newArchiveTarGzReader()
	if err != nil {
		return err
	}

	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}

		if err != nil {
			return err
		}

		if hdr.Name == name {
			f, err := createTempFile(tr)
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

func (a *ArchiveTarGz) List() ([]string, error) {
	tr, err := a.newArchiveTarGzReader()
	if err != nil {
		return nil, err
	}
	var names []string
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		names = append(names, hdr.Name)
	}
	return names, nil
}
