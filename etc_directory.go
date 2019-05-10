package etc

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func Env(key string) Path {
	return ParseSystemPath(os.Getenv(key))
}

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

// Mkdir accepts a variadic array of string elements, such as "dev", "random"
func (d Path) Mkdir(elements ...string) error {
	return d.Concat(elements...).MakeDir()
}

func (d Path) MakeDir() error {
	return os.MkdirAll(d.Render(), 0700)
}

func (d Path) MakeDirPath(path Path) error {
	return os.MkdirAll(d.GetSub(path).Render(), 0700)
}

func (d Path) String() string {
	return d.Render()
}

func (d Path) DiskStatus() *DiskStatus {
	dsk, err := GetDiskStatus(d.Render())
	if err != nil {
		panic(fmt.Errorf("%s: %s", d, err))
	}
	return dsk
}

func (d Path) Free() uint64 {
	return d.DiskStatus().Free
}

func (d Path) Stat() os.FileInfo {
	fi, _ := os.Stat(d.Render())
	return fi
}

func (d Path) Time() time.Time {
	return d.Stat().ModTime()
}

func (d Path) Size() uint64 {
	if d.IsDirectory() == false {
		return uint64(d.Stat().Size())
	}

	i, _ := dirSize(d.Render())
	return uint64(i)
}

func (d Path) DiskUsed() uint64 {
	return d.DiskStatus().Used
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

func (d Path) IsDirectory() bool {
	fi, err := os.Stat(d.Render())
	if err != nil {
		return false
	}

	return fi.IsDir()
}

func (d Path) ExistsPath(p Path) bool {
	_, err := os.Stat(d.GetSub(p).Render())
	return err == nil
}

func (d Path) Exists(p string) bool {
	return d.ExistsPath(parseNixPath([]rune(p)))
}

// Concat applies the path with a variadic list of path elements, and returns it.
func (d Path) Concat(s ...string) Path {
	y := make(Path, len(d))
	copy(y, d)
	y = append(y, s...)
	return y
}

func (d Path) Pop() (string, Path) {
	y := d
	return y[len(d)-1], y[:len(d)-2]
}

func (d Path) WriteAll(b []byte) error {
	return ioutil.WriteFile(d.Render(), b, 0700)
}

func (d Path) ReadAll() ([]byte, error) {
	return ioutil.ReadFile(d.Render())
}

func (d Path) Remove() error {
	return os.RemoveAll(d.Render())
}

func (d Path) LRU() (string, error) {
	if d.IsDirectory() == false {
		return "", fmt.Errorf("etc: cannot get LRU of a non-directory")
	}

	fis, err := ioutil.ReadDir(d.Render())
	if err != nil {
		return "", err
	}

	if len(fis) == 0 {
		return "", fmt.Errorf("etc: no directory children in %s", d)
	}

	now := time.Now()
	var oldest *os.FileInfo

	for _, v := range fis {
		if oldest == nil {
			oldest = &v
			continue
		}

		old := *oldest

		if now.Sub(v.ModTime()) > now.Sub(old.ModTime()) {
			oldest = &v
		}
	}

	old := *oldest

	return old.Name(), nil
}

func dirSize(path string) (int64, error) {
	var size int64
	err := filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			size += info.Size()
		}
		return err
	})
	return size, err
}
