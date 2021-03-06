package main

import (
	"encoding/json"
	"flag"
	"os"
	"path/filepath"
	"runtime"

	"github.com/asticode/go-astilectron-bundler"
	"github.com/asticode/go-astilog"
	"github.com/asticode/go-astitools/flag"
	"github.com/pkg/errors"
)

// Flags
var (
	astilectronPath   = flag.String("a", "", "the astilectron path")
	configurationPath = flag.String("c", "", "the configuration path")
	darwin            = flag.Bool("d", false, "if set, will add darwin/amd64 to the environments")
	linux             = flag.Bool("l", false, "if set, will add linux/amd64 to the environments")
	windows           = flag.Bool("w", false, "if set, will add windows/amd64 to the environments")
	environmentFilter = flag.String("e", "", "if set, will only match environments matching pattern.")
)

func main() {
	// Init
	var s = astiflag.Subcommand()
	flag.Parse()
	astilog.FlagInit()

	// Get configuration path
	var cp = *configurationPath
	var err error
	if len(cp) == 0 {
		// Get working directory path
		var wd string
		if wd, err = os.Getwd(); err != nil {
			astilog.Fatal(errors.Wrap(err, "os.Getwd failed"))
		}

		// Set configuration path
		cp = filepath.Join(wd, "bundler.json")
	}

	// Open file
	var f *os.File
	if f, err = os.Open(cp); err != nil {
		astilog.Fatal(errors.Wrapf(err, "opening file %s failed", cp))
	}
	defer f.Close()

	// Unmarshal
	var c *astibundler.Configuration
	if err = json.NewDecoder(f).Decode(&c); err != nil {
		astilog.Fatal(errors.Wrap(err, "unmarshaling configuration failed"))
	}

	// Astilectron path
	if len(*astilectronPath) > 0 {
		c.AstilectronPath = *astilectronPath
	}

	// Environments
	if *darwin {
		c.Environments = append(c.Environments, astibundler.ConfigurationEnvironment{Arch: "amd64", OS: "darwin"})
	}
	if *linux {
		c.Environments = append(c.Environments, astibundler.ConfigurationEnvironment{Arch: "amd64", OS: "linux"})
	}
	if *windows {
		c.Environments = append(c.Environments, astibundler.ConfigurationEnvironment{Arch: "amd64", OS: "windows"})
	}
	if len(c.Environments) == 0 {
		c.Environments = []astibundler.ConfigurationEnvironment{{Arch: runtime.GOARCH, OS: runtime.GOOS}}
	}

	// Environment filtering.
	if len(*environmentFilter) > 0 {
		c.EnvironmentFilter = *environmentFilter
	}

	// Build bundler
	var b *astibundler.Bundler
	if b, err = astibundler.New(c); err != nil {
		astilog.Fatal(errors.Wrap(err, "building bundler failed"))
	}

	// Handle signals
	b.HandleSignals()

	// Switch on subcommand
	switch s {
	case "bd":
		// Bind data
		if err = b.BindData(runtime.GOOS, runtime.GOARCH, ""); err != nil {
			astilog.Fatal(errors.Wrap(err, "binding data failed"))
		}
	case "cc":
		// Clear cache
		if err = b.ClearCache(); err != nil {
			astilog.Fatal(errors.Wrap(err, "clearing cache failed"))
		}
	default:
		// Bundle
		if err = b.Bundle(); err != nil {
			astilog.Fatal(errors.Wrap(err, "bundling failed"))
		}
	}
}
