package main

import (
	"fmt"
	"os/exec"
	"path/filepath"

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
	jsOut = pflag.StringP("js_out", "j", "", "javascript output")
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

	fi := loadSyntax(srcFile, *pkg)

	if *goOut != "" {
		compileGo(fi, *goOut+"/"+*pkg)
	}

	if *jsOut != "" {
		compileJS(fi, filepath.Join(*jsOut, *pkg))
	}
}

func loadSyntax(path, pkg string) *idl.Syntax {
	prog, err := ioutil.ReadFile(path)
	if err != nil {
		fatalf("%s", err)
	}

	t, err := idl.Parse(string(prog))
	if err != nil {
		fatalf("%s", err)
	}

	t.PackageName = pkg
	return t
}

func fatalf(f string, args ...interface{}) {
	fmt.Printf(f, args...)
	os.Exit(-1)
}

func compileGo(t *idl.Syntax, out string) {
	st, rpc := t.GenerateGo()

	out1 := out + ".etc.go"
	out2 := out + "-rpc.go"

	compileString(out1, st)
	gofmt(out1)

	if !exists(out2) {
		if rpc != "" {
			compileString(out2, rpc)
			ioutil.WriteFile(out2, []byte(rpc), 0700)
			gofmt(out2)
		}
	}
}

func compileString(path, data string) {
	fmt.Println("Compiling", path)
	ioutil.WriteFile(path, []byte(data), 0700)
}

// TypeScript declarations?
func compileJS(t *idl.Syntax, out string) {
	st, rpc := t.GenerateJS()

	out1 := out + ".etc.js"
	out2 := out + "-rpc.js"

	compileString(out1, st)

	if !exists(out2) {
		if rpc != "" {
			compileString(out2, rpc)
		}
	}
}

func gofmt(out string) {
	mc := exec.Command("gofmt", "-w", out)
	ot := etc.NewBuffer()
	mc.Stdout = ot
	mc.Stderr = ot
	c := mc.Run()
	if c != nil {
		fatalf("%s (%s)", c, ot)
	}
}

func exists(path string) bool {
	if _, err := os.Stat(path); err == nil {
		return true
	}

	return false
}
