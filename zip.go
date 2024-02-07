package tobuilds

import (
	"archive/zip"
	"fmt"
)

type ArchiveZip struct {
	ctx *Ctx
	Platform
	name string
	z    *zip.ReadCloser
}

func newArchiveZip(ctx *Ctx, platform Platform, name string) (*ArchiveZip, error) {
	fmt.Println("INFO: registered new zip archive", name)
	f, err := ctx.GetFile(name)
	if err != nil {
		return nil, err
	}

	z, err := zip.OpenReader(f.Name())
	if err != nil {
		return nil, err
	}

	fmt.Println("INFO: read from archive", name)
	return &ArchiveZip{ctx, platform, name, z}, nil
}

func (a *ArchiveZip) Run(name string, options ...string) error {
	if !a.Platform.isCurrent() {
		fmt.Printf("INFO: skipped run (%s %s) for different platform\n", name, options)
		return nil
	}

	for _, f := range a.z.File {
		if f.Name == name {
			r, err := f.Open()
			if err != nil {
				return err
			}

			f, err := createTempFile(r)
			if err != nil {
				return err
			}
			f.Close()
			fmt.Printf("INFO: run %s from archive %s\n", name, a.name)
			return RunFile(a.Platform, f, options...)
		}
	}
	return nil
}

func (a *ArchiveZip) List() []*zip.File { return a.z.File }
