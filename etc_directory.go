package etc

import (
	"io"
	"os"
	"runtime"
	"strings"
)

type Path []string

func splitDir(input []rune, sep rune) []string {
	var out []string
	var cur string

	for i := 0; i < len(input); i++ {
		if input[i] == sep {
			if cur != "" {
				if cur == ".." || cur == "." || cur == "..." {
					cur = ""
					continue
				}
				out = append(out, cur)
				cur = ""
			}
			continue
		}

		cur += string(input[i])
	}

	if cur != "" {
		out = append(out, cur)
	}

	return out
}

func ParseWindowsPath(s string) Path {
	return parseWinPath([]rune(s))
}

func ParseUnixPath(s string) Path {
	return parseNixPath([]rune(s))
}

func parseWinPath(r []rune) Path {
	var out []string

	p := '\\'

	if string(r) == "." {
		f, err := os.Getwd()
		if err != nil {
			panic(err)
		}
		return parseWinPath([]rune(f))
	}

	if len(r) <= 1 {
		panic(string(r))
	}

	if r[1] == ':' && r[2] == '\\' {
		out = append(out, "...")
		out = append(out, string(r[0]))
		out = append(out, splitDir(r[3:], p)...)
		return Path(out)
	}

	if r[1] == ':' && r[2] == '/' {
		out = append(out, "...")
		out = append(out, string(r[0]))
		out = append(out, splitDir(r[3:], p)...)
		return Path(out)
	}

	return splitDir(r, p)
}

func parseNixPath(r []rune) Path {
	var out []string

	if r[0] == '/' {
		out = append(out, "...", "root")
		return append(out, splitDir(r[1:], '/')...)
	}

	return splitDir(r, '/')
}

func ParseSystemPath(s string) Path {
	r := []rune(s)
	if runtime.GOOS == "windows" {
		return parseWinPath(r)
	}

	return parseNixPath(r)
}

func (d Path) RenderWin() string {
	if d[0] == "..." {
		prefix := strings.ToUpper(d[1])
		if prefix == "root" {
			prefix = "C"
		}
		out := prefix + ":" + "\\"
		return out + strings.Join(d[2:], "\\")
	} else {
		return strings.Join(d, "\\")
	}
}

func (d Path) RenderUnix() string {
	if d[0] == "..." {
		return "/" + strings.Join(d[2:], "/")
	}

	return strings.Join(d, "/")
}

func (d Path) Render() string {
	if runtime.GOOS == "windows" {
		return d.RenderWin()
	}

	return d.RenderUnix()
}

func TmpDirectory() Path {
	if runtime.GOOS == "windows" {
		return ParseSystemPath(os.Getenv("TEMP"))
	}

	return ParseSystemPath("/tmp/")
}

func LocalDirectory() Path {
	if runtime.GOOS == "windows" {
		return ParseSystemPath(os.Getenv("USERPROFILE") + "\\AppData\\Local")
	}

	return ParseSystemPath(os.Getenv("HOME") + "/.local/share/")
}

func (d Path) GetSub(p Path) Path {
	return append(d, p...)
}

func (d Path) GetSubPath(p Path) string {
	return d.GetSub(p).Render()
}

func (d Path) GetSubFile(path Path) (*Buffer, error) {
	return FileController(d.GetSubPath(path))
}

func (d Path) Get(path string) (*Buffer, error) {
	if path[0] == '/' {
		path = path[1:]
	}
	return d.GetSubFile(parseNixPath([]rune(path)))
}

func (d Path) MakeDirPath(path Path) error {
	return os.MkdirAll(d.GetSub(path).Render(), 0700)
}

func (d Path) MakeDir(path string) error {
	return d.MakeDirPath(ParseUnixPath(path))
}

func (d Path) Put(path string, data io.Reader) error {
	e, err := d.Get(path)
	if err == nil {
		e.Flush()
	} else {
		return err
	}

	_, err = io.Copy(e, data)
	if err != nil {
		e.Close()
		return err
	}

	return e.Close()
}

func (d Path) IsExtant() bool {
	_, err := os.Stat(d.Render())
	return err == nil
}

func (d Path) ExistsPath(p Path) bool {
	_, err := os.Stat(d.GetSub(p).Render())
	return err == nil
}

func (d Path) Exists(p string) bool {
	return d.ExistsPath(parseNixPath([]rune(p)))
}

func Env(variable string) Path {
	return ParseSystemPath(os.Getenv(variable))
}

func (d Path) Concat(s ...string) Path {
	y := d
	y = append(y, s...)
	return y
}

func (d Path) Pop() (string, Path) {
	y := d
	return y[len(d)-1], y[:len(d)-2]
}
