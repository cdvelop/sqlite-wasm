// Copyright 2021 The Sqlite Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build freebsd
// +build freebsd

package driver

import (
	"syscall"
)

func setMaxOpenFiles(n int64) error {
	var rLimit syscall.Rlimit
	rLimit.Max = uint64(n)
	rLimit.Cur = uint64(n)
	return syscall.Setrlimit(syscall.RLIMIT_NOFILE, &rLimit)
}
