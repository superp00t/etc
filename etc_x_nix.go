// +build !windows

package etc

import (
	"os"
	"syscall"
)

func ParseSystemPath(s string) Path {
	return parseNixPath([]rune(s))
}

func TmpDirectory() Path {
	return ParseSystemPath("/tmp/")
}

func LocalDirectory() Path {
	return ParseSystemPath(os.Getenv("HOME") + "/.local/share/")
}

func HomeDirectory() Path {
	return ParseSystemPath(os.Getenv("HOME"))
}

func (d Path) Render() string {
	return d.RenderUnix()
}

func getDiskStatus(path string) (*DiskStatus, error) {
	var disk DiskStatus
	fs := syscall.Statfs_t{}
	err := syscall.Statfs(path, &fs)
	if err != nil {
		return nil, err
	}
	disk.All = fs.Blocks * uint64(fs.Bsize)
	disk.Free = fs.Bfree * uint64(fs.Bsize)
	disk.Used = disk.All - disk.Free
	return &disk, nil
}
