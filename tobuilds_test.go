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
	c.Run(tobuilds.PlatformAny, "test/test.sh")
}
