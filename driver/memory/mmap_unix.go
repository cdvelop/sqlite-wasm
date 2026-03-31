// Copyright 2011 Evan Shaw. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE-MMAP-GO file.

//go:build unix

// Modifications (c) 2017 The Memory Authors.

package memory

import (
	"golang.org/x/sys/unix"
	"os"
	"unsafe"
)

const pageSizeLog = 16

var (
	osPageMask = osPageSize - 1
	osPageSize = os.Getpagesize()
)

func unmap(addr uintptr, size int) error {
	return unix.Munmap((*[1 << 30]byte)(unsafe.Pointer(addr))[:size:size])
}

// pageSize aligned.
func mmap(size int) (uintptr, int, error) {
	size = roundup(size, osPageSize)
	// Ask for more so we can align the result at a pageSize boundary
	n := size + pageSize
	data, err := unix.Mmap(-1, 0, n, unix.PROT_READ|unix.PROT_WRITE, unix.MAP_PRIVATE|unix.MAP_ANON)
	if err != nil {
		return 0, 0, err
	}

	p := uintptr(unsafe.Pointer(&data[0]))
	if p&uintptr(osPageMask) != 0 {
		panic("internal error")
	}

	mod := int(p) & pageMask
	if mod != 0 { // Return the extra part before pageSize aligned block
		m := pageSize - mod
		if err := unmap(p, m); err != nil {
			unmap(p, n) // Do not leak the first mmap
			return 0, 0, err
		}

		n -= m
		p += uintptr(m)
	}

	if p&uintptr(pageMask) != 0 {
		panic("internal error")
	}

	if n > size { // Return the extra part after pageSize aligned block
		if err := unmap(p+uintptr(size), n-size); err != nil {
			// Do not error when the kernel rejects the extra part after, just return the
			// unexpectedly enlarged size.
			//
			// Fixes the bigsort.test failures on linux/s390x, see: https://gitlab.com/cznic/sqlite/-/issues/207
			size = n
		}
	}

	return p, size, nil
}
