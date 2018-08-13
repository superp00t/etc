package yo

import (
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/fatih/color"
	"github.com/superp00t/etc"
)

var (
	f *FlagParser
)

type Type int

const (
	String Type = iota
	Bool
	Int64
	Float64
	Duration
)

type Def struct {
	Type       Type
	Long       string
	Definition string
	Value      interface{}
}

type Routine struct {
	Handler     func([]string)
	Arguments   []string
	Description string
}

type FlagParser struct {
	Called   []string
	Defs     map[string]*Def
	Routines map[string]*Routine
}

func (d *Def) parseValue(s string) error {
	var err error
	switch d.Type {
	case String:
		d.Value = s
	case Int64:
		d.Value, err = strconv.ParseInt(s, 0, 64)
	case Float64:
		d.Value, err = strconv.ParseFloat(s, 64)
	case Duration:
		d.Value, err = time.ParseDuration(s)
	default:
		return fmt.Errorf("unknown type %d", d.Type)
	}
	return err
}

func getValue(in *etc.Buffer) string {
	o := []rune{}

	op := 0

	for {
		i, _, err := in.ReadRune()
		if err != nil {
			break
		}

		if i == '=' {
			continue
		}

		if i == ' ' && op == 0 {
			op++
			continue
		}

		if i == ' ' {
			break
		}

		if i == 0 {
			break
		}
		o = append(o, i)
		op++
	}

	return string(o)
}

func (f *FlagParser) Parse(s string) error {
	e := etc.FromString(s)

	for {
	cont:
		g := e.Rpos()
		chr, _, err := e.ReadRune()
		if err != nil {
			goto flagend
		}

		if chr == '-' {
			chr2, _, err := e.ReadRune()
			if err != nil {
				goto flagend
			}

			// --option
			if chr2 == '-' {
				optName := []rune{}

				for {
					chr3, _, err := e.ReadRune()
					if err == io.EOF {
						break
					}

					if err != nil {
						goto flagend
					}

					if chr3 == ' ' {
						break
					}

					optName = append(optName, chr3)
				}

				name := string(optName)

				var found *Def = nil

				for _, v := range f.Defs {
					if v.Long == name {
						found = v
						break
					}
				}

				if found == nil {
					return fmt.Errorf("unknown flag name %s", name)
				}

				if found.Type == Bool {
					found.Value = true
					goto cont
				}

				data := getValue(e)

				err := found.parseValue(data)
				if err != nil {
					return err
				}

				goto cont
			}

			short := string(chr2)

			if f.Defs[short] == nil {
				return fmt.Errorf("unknown flag shorthand %s", short)
			}

		} else {
			e.SeekR(g)
			f.Called = append(f.Called, getValue(e))
			goto cont
		}
	}
flagend:
	return nil
}

func setup() {
	if f == nil {
		f = new(FlagParser)
		f.Defs = make(map[string]*Def)
		f.Routines = make(map[string]*Routine)
	}
}

func FSpew() {
	fmt.Println(spew.Sdump(f))
}

// Boolf cannot be true by default.
func Boolf(short, long, description string) {
	setup()
	f.Defs[short] = &Def{Bool, long, description, false}
}

func Durationf(short, long, description string, t time.Duration) {
	setup()
	f.Defs[short] = &Def{Duration, long, description, t}
}

func Int64f(short, long, description string, def int64) {
	setup()
	f.Defs[short] = &Def{Int64, long, description, def}
}

func Float64f(short, long, description string, def float64) {
	setup()
	f.Defs[short] = &Def{Float64, long, description, def}
}

func Stringf(short, long, description, def string) {
	setup()
	f.Defs[short] = &Def{String, long, description, def}
}

func (s *FlagParser) SortedRoutines() []string {
	var r []string

	for k := range s.Routines {
		r = append(r, k)
	}

	sort.Strings(r)

	return r
}

func (s *FlagParser) SortedDefs() []string {
	var r []string

	for k := range s.Defs {
		r = append(r, k)
	}

	sort.Strings(r)

	return r
}

func Init() {
	if err := f.Parse(strings.Join(os.Args[1:], " ")); err != nil {
		Fatal(err)
	}

	call := ""

	if len(f.Called) != 0 {
		call = f.Called[0]
	} else {
		f.Called = append(f.Called, "")
	}

	exeNameL := strings.Split(os.Args[0], string(os.PathSeparator))
	exe := exeNameL[len(exeNameL)-1]
	cl := f.Routines[call]
	if cl == nil && call == "" || call == "help" {
		color.Set(color.FgGreen)
		fmt.Printf("%s:\n\n", exe)
		color.Unset()
		for _, v := range f.SortedRoutines() {
			if v == "" {
				color.Set(color.FgCyan)
				fmt.Printf("  %s\n\n", f.Routines[v].Description)
				color.Set(color.FgGreen)
				fmt.Printf("Subcommands:\n\n")
				color.Unset()
			} else {
				fmt.Printf("  %s ", exe)
				color.Set(color.FgCyan)
				fmt.Printf("%s ", v)
				color.Unset()
				if len(f.Routines[v].Arguments) > 0 {
					color.Set(color.FgHiBlue)
					for _, arg := range f.Routines[v].Arguments {
						fmt.Printf("[%s] ", arg)
					}
					color.Unset()
				}
				fmt.Printf("\n\n    %s\n\n", f.Routines[v].Description)
			}
		}

		fmt.Println()
		color.Set(color.FgGreen)
		fmt.Println("Options:")
		color.Unset()
		fmt.Println()

		for _, v := range f.SortedDefs() {
			fmt.Printf("  --%s, -%s\n\n   %s\n", f.Defs[v].Long, v, f.Defs[v].Definition)
			fmt.Println()
		}
		return
	}

	if cl == nil {
		Fatal("No handler for routine: \"" + call + "\"")
	}

	if call == "" {
		cl.Handler(nil)
	} else {
		arg := make([]string, len(cl.Arguments))
		for x := 2; x < len(os.Args); x++ {
			arg[x-2] = os.Args[x]
		}
		cl.Handler(arg)
	}
}

func AddSubroutine(name string, arguments []string, description string, fn func([]string)) {
	if name == "help" {
		panic("Cannot override help function")
	}
	f.Routines[name] = &Routine{fn, arguments, description}
}

func Main(description string, fn func(s []string)) {
	AddSubroutine("", nil, description, fn)
}