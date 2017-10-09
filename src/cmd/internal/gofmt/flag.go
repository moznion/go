// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gofmt

import "flag"

type boolFlag struct {
	Name  string
	Value bool
}

type stringFlag struct {
	Name  string
	Value string
}

type Flag struct {
	List        boolFlag
	Write       boolFlag
	RewriteRule stringFlag
	SimplifyAST boolFlag
	DoDiff      boolFlag
	AllErrors   boolFlag
	Cpuprofile  stringFlag
}

func InitGofmtFlag(fs *flag.FlagSet) *Flag {
	f := Flag{
		List:        boolFlag{Name: "l"},
		Write:       boolFlag{Name: "w"},
		RewriteRule: stringFlag{Name: "r"},
		SimplifyAST: boolFlag{Name: "s"},
		DoDiff:      boolFlag{Name: "d"},
		AllErrors:   boolFlag{Name: "e"},
		Cpuprofile:  stringFlag{Name: "cpuprofile"},
	}

	// main operation modes
	fs.BoolVar(&f.List.Value, f.List.Name, false, "list files whose formatting differs from gofmt's")
	fs.BoolVar(&f.Write.Value, f.Write.Name, false, "write result to (source) file instead of stdout")
	fs.StringVar(&f.RewriteRule.Value, f.RewriteRule.Name, "", "rewrite rule (e.g., 'a[b:len(a)] -> a[b:]')")
	fs.BoolVar(&f.SimplifyAST.Value, f.SimplifyAST.Name, false, "simplify code")
	fs.BoolVar(&f.DoDiff.Value, f.DoDiff.Name, false, "display diffs instead of rewriting files")
	fs.BoolVar(&f.AllErrors.Value, f.AllErrors.Name, false, "report all errors (not just the first 10 on different lines)")

	// debugging
	fs.StringVar(&f.Cpuprofile.Value, f.Cpuprofile.Name, "", "write cpu profile to this file")

	return &f
}
