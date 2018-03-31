package main

import (
	"fmt"
	"os/exec"

	"github.com/superp00t/etc"

	"io/ioutil"
	"os"
	"strings"

	"github.com/ogier/pflag"
	"github.com/superp00t/etc/idl"
)

var (
	pkg   = pflag.StringP("pkg", "p", "", "package string")
	goOut = pflag.StringP("go_out", "g", "", "golang output")
)

func main() {
	pflag.Parse()

	srcFile := pflag.Arg(0)
	if srcFile == "" {
		fatalf("usage: %s <file.etcschema> --pkg=example --go_out=<example directory>\n", os.Args[0])
	}

	if *pkg == "" {
		f := strings.Split(srcFile, ".")
		*pkg = f[0]
	}

	if *goOut != "" {
		compile(srcFile, *pkg, *goOut+"/"+*pkg+".etc.go")
	}
}

func fatalf(f string, args ...interface{}) {
	fmt.Printf(f, args...)
	os.Exit(-1)
}

func compile(in, pkgName, out string) {
	prog, err := ioutil.ReadFile(in)
	if err != nil {
		fatalf("%s", err)
	}

	t, err := idl.Parse(string(prog))
	if err != nil {
		fatalf("%s", err)
	}

	t.PackageName = pkgName

	st := t.GenerateGo()

	ioutil.WriteFile(out, []byte(st), 0700)

	mc := exec.Command("gofmt", "-w", out)
	ot := etc.NewBuffer()
	mc.Stdout = ot
	mc.Stderr = ot
	c := mc.Run()
	if c != nil {
		fatalf("%s (%s)", c, ot)
	}
}
