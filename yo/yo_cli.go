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

func GetValue(s string) interface{} {
	if f == nil {
		return nil
	}

	l := f.Defs[s]
	if l == nil {
		return nil
	}

	return l.Value
}

func Int64G(s string) int64 {
	i := GetValue(s)
	if i == nil {
		return -1
	}

	return i.(int64)
}

func StringG(s string) string {
	i := GetValue(s)
	if i == nil {
		return ""
	}

	return i.(string)
}

func Float64G(s string) float64 {
	i := GetValue(s)
	if i == nil {
		return 0
	}

	return i.(float64)
}

func BoolG(s string) bool {
	i := GetValue(s)
	if i == nil {
		return false
	}

	return i.(bool)
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
	case Bool:
		d.Value = true
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

		if i == 0 && op == 0 {
			for {
				i, _, err := in.ReadRune()
				if err != nil {
					return string(o)
				}

				if i == 0 {
					return string(o)
				} else {
					o = append(o, i)
				}
			}
		}

		if i == ' ' && op == 0 {
			op++
			continue
		}

		if i == ' ' && op != 0 {
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

		if chr == ' ' {
			continue
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

			if f.Defs[short].Type == Bool {
				f.Defs[short].Value = true
				goto cont
			}

			data := getValue(e)

			err = f.Defs[short].parseValue(data)
			if err != nil {
				return err
			}

			goto cont
		} else {
			e.SeekR(g)
			f.Called = append(f.Called, getValue(e))
			goto cont
		}
	}
flagend:

	return nil
}

func setup(s, l string) {
	if s == "h" || l == "help" {
		panic("You cannot override the help flags")
	}

	if s == "y" || l == "yo_level" {
		panic("You can't overwrite the yo log level flags")
	}

	if f == nil {
		f = new(FlagParser)
		f.Defs = make(map[string]*Def)
		f.Routines = make(map[string]*Routine)
	}
}

func FSpew() {
	setup("", "")
	fmt.Println(spew.Sdump(f))
}

// Boolf cannot be true by default.
func Boolf(short, long, description string) {
	setup(short, long)
	f.Defs[short] = &Def{Bool, long, description, false}
}

func Durationf(short, long, description string, t time.Duration) {
	setup(short, long)
	f.Defs[short] = &Def{Duration, long, description, t}
}

func Int64f(short, long, description string, def int64) {
	setup(short, long)
	f.Defs[short] = &Def{Int64, long, description, def}
}

func Float64f(short, long, description string, def float64) {
	setup(short, long)
	f.Defs[short] = &Def{Float64, long, description, def}
}

func Stringf(short, long, description, def string) {
	setup(short, long)
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
	setup("", "")

	f.Defs["y"] = &Def{Int64, "yo_level", "log level", int64(0)}
	f.Defs["h"] = &Def{Bool, "help", "prints this message", false}

	if err := f.Parse(strings.Join(os.Args[1:], "\x00")); err != nil {
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
	if cl == nil && call == "" || call == "help" || BoolG("h") {
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
				fmt.Printf("\n    %s\n", f.Routines[v].Description)
			}
		}

		fmt.Println()
		if len(f.SortedDefs()) > 0 {
			color.Set(color.FgGreen)
			fmt.Println("Options:")
			color.Unset()

			for _, v := range f.SortedDefs() {
				fmt.Printf("  --%s, -%s\n    %s\n    Default: %v\n", f.Defs[v].Long, v, f.Defs[v].Definition, f.Defs[v].Value)
				fmt.Println()
			}
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
		for x := 1; x < len(f.Called); x++ {
			if x-1 == len(arg) {
				break
			}
			arg[x-1] = f.Called[x]
		}
		cl.Handler(arg)
	}
}

func AddSubroutine(name string, arguments []string, description string, fn func([]string)) {
	setup("", "")

	if name == "help" {
		panic("Cannot override help function")
	}
	f.Routines[name] = &Routine{fn, arguments, description}
}

func Main(description string, fn func(s []string)) {
	AddSubroutine("", nil, description, fn)
}
