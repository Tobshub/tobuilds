# Tobuilds

A simple system for getting your app (and it's external dependencies) to your users.

This is still a rough idea for now, and I'm hacking at it as I go (no pun intended).

The idea is to be able to setup per os run instructions and easily pass command line options to the runners.

## Usage

```go
// main.go
package main

import (
    "embed"
    "github.com/tobshub/tobuilds"
)

// Embed build files into script
//go:embed build/*
var efs embed.FS

func main() {
    // Create a new context with the embed file system
    ctx := tobuilds.Init(efs)
    // `End` cleans up tmp files so make sure to call this :D
    defer tobuilds.End()
    // Run an embeded script on linux machines
    if err := ctx.Run(tobuilds.PlatFormLinux, "build/build.sh"); err != nil {
        panic(err)
    }
    // Run a remote script on windows machines
    if err := ctx.Run(tobuilds.PlatFormWindows, "https://ftp.fake.com/remote/build/script.exe"); err != nil {
        panic(err)
    }
    // Access an archive
    tr := ctx.NewArchiveTarGz(tobuilds.PlatformLinux, "build/example.tar.gz")
    if tr != nil {
        // run a file from the archive
        if err := tr.Run("example.sh", "arg1"); err != nil {
            panic(err)
        }
    }
}
```

## TODO
- [ ] Support for running msi files
