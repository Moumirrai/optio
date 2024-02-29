//go:build windows
// +build windows

package metadata

import (
	"io/fs"
	"os"
	"syscall"
	"time"
)

func GetCTime(f fs.FileInfo) time.Time {
	d := f.Sys().(*syscall.Win32FileAttributeData)
	cTime := time.Unix(0, d.CreationTime.Nanoseconds())
	return cTime
}

func SetCTime(f *os.File, t time.Time) error {

	h := f.Fd()
	cTime := syscall.NsecToFiletime(t.UnixNano())
	return syscall.SetFileTime(syscall.Handle(h), &cTime, nil, nil)
}
