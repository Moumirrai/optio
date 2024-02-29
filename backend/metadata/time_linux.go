//go:build linux
// +build linux

package metadata

import (
	"io/fs"
	"os"
	"syscall"
	"time"
)

func GetCTime(f fs.FileInfo) time.Time {
	d := f.Sys().(*syscall.Stat_t)
	cTime := time.Unix(d.Ctim.Sec, d.Ctim.Nsec)
	return cTime
}

func SetCTime(f *os.File, t time.Time) error {

	h := f.Fd()
	cTime := syscall.NsecToTimespec(t.UnixNano())
	return syscall.UtimesNanoAt(int(h), "", []syscall.Timespec{cTime}, syscall.AT_SYMLINK_NOFOLLOW)
}
