// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build amd64 && darwin
// +build amd64,darwin

package unix

import (
	"errors"
	"strconv"
	"strings"
	"syscall"
)

func setTimespec(sec, nsec int64) Timespec {
	return Timespec{Sec: sec, Nsec: nsec}
}

func setTimeval(sec, usec int64) Timeval {
	return Timeval{Sec: sec, Usec: int32(usec)}
}

func SetKevent(k *Kevent_t, fd, mode, flags int) {
	k.Ident = uint64(fd)
	k.Filter = int16(mode)
	k.Flags = uint16(flags)
}

func (iov *Iovec) SetLen(length int) {
	iov.Len = uint64(length)
}

func (msghdr *Msghdr) SetControllen(length int) {
	msghdr.Controllen = uint32(length)
}

func (msghdr *Msghdr) SetIovlen(length int) {
	msghdr.Iovlen = int32(length)
}

func (cmsg *Cmsghdr) SetLen(length int) {
	cmsg.Len = uint32(length)
}

var buggyVersion = [...]int{6153, 141, 1}

func xnuKernelBug25397314(syscallName string) (bool, error) {
	// Workaround for a kernel bug in macOS Catalina when using the kern.procargs2 syscall
	// More information about this bug can be found here:
	// https://github.com/apple-oss-distributions/xnu/blob/xnu-7195.50.7.100.1/bsd/kern/kern_sysctl.c#L1552-#L1592

	if !strings.Contains(syscallName, "kern.procargs2") {
		return false, nil
	}

	uname := Utsname{}
	err := Uname(&uname)
	if err != nil {
		return false, err
	}

	return buggyKernel(uname)
}

func buggyKernel(uname Utsname) (bool, error) {
	versionStr := string(uname.Version[:])
	xnuVersionStartStr := "xnu-"
	xnuVersionStartIndex := strings.Index(versionStr, xnuVersionStartStr) + len(xnuVersionStartStr)
	if xnuVersionStartIndex == (-1 + len(xnuVersionStartStr)) {
		return false, errors.New("could not find xnu version number in uname")
	}
	xnuVersionEndIndex := strings.Index(versionStr[xnuVersionStartIndex:], "~") + xnuVersionStartIndex
	if xnuVersionEndIndex == -1 {
		return false, errors.New("could not find xnu version number terminus")
	}
	xnuVersionStr := versionStr[xnuVersionStartIndex:xnuVersionEndIndex]
	xnuVersionStrSplit := strings.Split(xnuVersionStr, ".")
	xnuVersion := make([]int, len(xnuVersionStrSplit))
	var err error
	for i := 0; i < len(xnuVersionStrSplit); i++ {
		xnuVersion[i], err = strconv.Atoi(xnuVersionStrSplit[i])
		if err != nil {
			return false, err
		}
	}

	for i := 0; i < len(buggyVersion); i++ {
		if i >= len(xnuVersion) {
			return true, nil
		} else if buggyVersion[i] > xnuVersion[i] {
			return true, nil
		} else if buggyVersion[i] < xnuVersion[i] {
			return false, nil
		}
	}

	return len(buggyVersion) == len(xnuVersion), nil

}

func Syscall9(num, a1, a2, a3, a4, a5, a6, a7, a8, a9 uintptr) (r1, r2 uintptr, err syscall.Errno)

//sys	Fstat(fd int, stat *Stat_t) (err error) = SYS_FSTAT64
//sys	Fstatat(fd int, path string, stat *Stat_t, flags int) (err error) = SYS_FSTATAT64
//sys	Fstatfs(fd int, stat *Statfs_t) (err error) = SYS_FSTATFS64
//sys	getfsstat(buf unsafe.Pointer, size uintptr, flags int) (n int, err error) = SYS_GETFSSTAT64
//sys	Lstat(path string, stat *Stat_t) (err error) = SYS_LSTAT64
//sys	ptrace1(request int, pid int, addr uintptr, data uintptr) (err error) = SYS_ptrace
//sys	ptrace1Ptr(request int, pid int, addr unsafe.Pointer, data uintptr) (err error) = SYS_ptrace
//sys	Stat(path string, stat *Stat_t) (err error) = SYS_STAT64
//sys	Statfs(path string, stat *Statfs_t) (err error) = SYS_STATFS64
