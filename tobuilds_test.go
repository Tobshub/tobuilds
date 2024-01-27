package tobuilds_test

import (
	"embed"
	"testing"

	"github.com/tobshub/tobuilds"
)

//go:embed test/*
var efs embed.FS

func TestRun(t *testing.T) {
	c := tobuilds.Init(&efs)
	err := c.Run(tobuilds.PlatformLinux, "test/test.sh")
	err = c.Run(tobuilds.PlatformLinux,
		"https://github.com/tobsdb/tobsdb/releases/download/v0.1.2-alpha/tdb-generate_Linux_x86_64", "-schema", "$TABLE a {\nb String\n}")
	err = c.Run(tobuilds.PlatformWindows,
		"https://github.com/tobsdb/tobsdb/releases/download/v0.1.2-alpha/tdb-generate_Windows_x86_64.exe", "-schema", "$TABLE a {\nb String\n}")
	tobuilds.End()
	if err != nil {
		t.Fatal(err)
	}
}
