// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package fmtcmd implements the ``go fmt'' command.
package fmtcmd

import (
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"sync"

	"cmd/go/internal/base"
	"cmd/go/internal/cfg"
	"cmd/go/internal/load"
	"cmd/go/internal/str"
	"cmd/internal/gofmt"
)

var gofmtFlag *gofmt.Flag

func init() {
	base.AddBuildFlagsNX(&CmdFmt.Flag)
	gofmtFlag = gofmt.InitGofmtFlag(&CmdFmt.Flag)
}

var CmdFmt = &base.Command{
	Run:       runFmt,
	UsageLine: "fmt [-n] [-x] [packages]",
	Short:     "run gofmt on package sources",
	Long: `
Fmt runs the command 'gofmt -l -w' on the packages named
by the import paths. It prints the names of the files that are modified.

For more about gofmt, see 'go doc cmd/gofmt'.
For more about specifying packages, see 'go help packages'.

The -n flag prints commands that would be executed.
The -x flag prints commands as they are executed.

To run gofmt with specific options, run gofmt itself.

See also: go fix, go vet.
	`,
}

func runFmt(cmd *base.Command, args []string) {
	gofmt := gofmtPath()
	procs := runtime.GOMAXPROCS(0)
	var wg sync.WaitGroup
	wg.Add(procs)
	fileC := make(chan string, 2*procs)

	var gofmtOptions []string

	rv := reflect.ValueOf(*gofmtFlag)
	fieldNum := rv.Type().NumField()
	for i := 0; i < fieldNum; i++ {
		f := rv.Field(i)

		var ok bool
		var name string
		if name, ok = f.FieldByName("Name").Interface().(string); !ok {
			panic("Name is not string")
		}
		name = "-" + name

		// TODO Should treat "-l" and "-w" as default

		v := f.FieldByName("Value")
		if sv, ok := v.Interface().(string); ok {
			if len(sv) > 0 {
				gofmtOptions = append(gofmtOptions, name, sv)
			}
		} else if bv, ok := v.Interface().(bool); ok {
			if bv {
				gofmtOptions = append(gofmtOptions, name)
			}
		} else {
			panic("Value type is not suitable")
		}
	}

	for i := 0; i < procs; i++ {
		go func() {
			defer wg.Done()
			for file := range fileC {
				base.Run(str.StringList(gofmt, gofmtOptions, file))
			}
		}()
	}
	for _, pkg := range load.Packages(args) {
		// Use pkg.gofiles instead of pkg.Dir so that
		// the command only applies to this package,
		// not to packages in subdirectories.
		files := base.RelPaths(pkg.InternalAllGoFiles())
		for _, file := range files {
			fileC <- file
		}
	}
	close(fileC)
	wg.Wait()
}

func gofmtPath() string {
	gofmt := "gofmt"
	if base.ToolIsWindows {
		gofmt += base.ToolWindowsExtension
	}

	gofmtPath := filepath.Join(cfg.GOBIN, gofmt)
	if _, err := os.Stat(gofmtPath); err == nil {
		return gofmtPath
	}

	gofmtPath = filepath.Join(cfg.GOROOT, "bin", gofmt)
	if _, err := os.Stat(gofmtPath); err == nil {
		return gofmtPath
	}

	// fallback to looking for gofmt in $PATH
	return "gofmt"
}
