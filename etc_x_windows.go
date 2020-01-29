// +build windows

package etc

import (
	"os"
	"syscall"
	"unsafe"
)

var (
	kernel32             = syscall.MustLoadDLL("kernel32.dll")
	getDiskFreeSpaceProc = kernel32.MustFindProc("GetDiskFreeSpaceExW")
)

func ParseSystemPath(s string) Path {
	if s == "" {
		wd, _ := os.Getwd()
		return parseWinPath(wd)
	}
	return parseWinPath([]rune(s))
}

func Root() Path {
	return ParseSystemPath("C:\\")
}

func TmpDirectory() Path {
	return Env("TEMP")
}

func LocalDirectory() Path {
	return Env("USERPROFILE").Concat("AppData", "Local")
}

func HomeDirectory() Path {
	return Env("USERPROFILE")
}

func (d Path) Render() string {
	return d.RenderWin()
}

func getDiskStatus(path string) (*DiskStatus, error) {
	var freeBytes int64
	var totalBytes int64
	var availBytes int64

	_, _, err := getDiskFreeSpaceProc.Call(
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(path))),
		uintptr(unsafe.Pointer(&freeBytes)),
		uintptr(unsafe.Pointer(&totalBytes)),
		uintptr(unsafe.Pointer(&availBytes)),
	)

	if err != nil && err.Error() != "The operation completed successfully." {
		return nil, err
	}

	return &DiskStatus{
		uint64(totalBytes),
		uint64(totalBytes - freeBytes),
		uint64(freeBytes),
	}, nil
}
