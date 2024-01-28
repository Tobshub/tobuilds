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
	c := tobuilds.Init(&efs)
	defer tobuilds.End()
	assert.NilError(t, c.Run(tobuilds.PlatformLinux, "test/test.sh", "test one"))
	assert.NilError(t, c.Run(tobuilds.PlatformLinux,
		"https://github.com/tobsdb/tobsdb/releases/download/v0.1.2-alpha/tdb-generate_Linux_x86_64", "-schema", "$TABLE a {\nb String\n}"))
	assert.NilError(t, c.Run(tobuilds.PlatformWindows,
		"https://github.com/tobsdb/tobsdb/releases/download/v0.1.2-alpha/tdb-generate_Windows_x86_64.exe", "-schema", "$TABLE a {\nb String\n}"))
	tr := c.NewArchiveTarGz(tobuilds.PlatformLinux, "test/test.tar.gz")
	assert.NilError(t, tr.Run("test.sh", "test two"))
	assert.NilError(t, tr.Run("test.sh", "test three"))
	z, err := c.NewArchiveZip(tobuilds.PlatformLinux, "test/test.zip")
	assert.NilError(t, err)
	assert.NilError(t, z.Run("test.sh", "test four"))
}
