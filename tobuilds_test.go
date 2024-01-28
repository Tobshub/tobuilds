package tobuilds_test

import (
	"embed"
	"testing"

	"github.com/tobshub/tobuilds"
	"gotest.tools/assert"
)

//go:embed test/*
var efs embed.FS

func TestRun(t *testing.T) {
	ctx := tobuilds.Init(&efs)
	defer tobuilds.End()
	assert.NilError(t, ctx.Run(tobuilds.PlatformLinux, "test/test.sh", "test one"))
	assert.NilError(t, ctx.Run(tobuilds.PlatformLinux,
		"https://github.com/tobsdb/tobsdb/releases/download/v0.1.2-alpha/tdb-generate_Linux_x86_64", "-schema", "$TABLE a {\nb String\n}"))
	assert.NilError(t, ctx.Run(tobuilds.PlatformWindows,
		"https://github.com/tobsdb/tobsdb/releases/download/v0.1.2-alpha/tdb-generate_Windows_x86_64.exe", "-schema", "$TABLE a {\nb String\n}"))
	tr := ctx.NewArchiveTarGz(tobuilds.PlatformLinux, "test/test.tar.gz")
	if tr != nil {
		assert.NilError(t, tr.Run("test.sh", "test two"))
		assert.NilError(t, tr.Run("test.sh", "test three"))
	}
	z, err := ctx.NewArchiveZip(tobuilds.PlatformWindows, "test/test.zip")
	assert.NilError(t, err)
	if z != nil {
		assert.NilError(t, z.Run("test.sh", "test four"))
	}
}
