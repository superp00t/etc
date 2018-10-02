package idl

import (
	"fmt"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/superp00t/etc"
)

func TestIDL(t *testing.T) {
	e := etc.Env("GOPATH")

	tsts := e.Concat("src", "github.com", "superp00t", "etc", "_tests")

	for _, tk := range []string{
		"enum.etcSchema",
	} {
		file, err := tsts.Get(tk)
		if err != nil {
			t.Fatal(err)
		}

		src := file.ToString()

		syntax, err := Parse(src)
		if err != nil {
			t.Fatal(tsts.Concat(tk).Render(), err)
		}

		fmt.Println(spew.Sdump(syntax))
	}
}
