// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package runtime

import (
	"runtime/internal/atomic"
	"unsafe"
)

// GOMAXPROCS sets the maximum number of CPUs that can be executing
// simultaneously and returns the previous setting. If n < 1, it does not
// change the current setting.
// The number of logical CPUs on the local machine can be queried with NumCPU.
// This call will go away when the scheduler improves.

// GOMAXPROCS设置可以同时执行的最大CPU数，并返回以前的设置。
// 默认为runtime.NumCPU的值。
// 如果n <1，则不会更改当前设置。
// 直译：当调度程序改进时，此调用将消失。我理解：当调度程序改进时，此设置将会丢失。
// 注意：调用此函数会STW。
func GOMAXPROCS(n int) int {
	if GOARCH == "wasm" && n > 1 {
		// WebAssembly目前没有线程，因此只能使用一个CPU。
		n = 1 // WebAssembly has no threads yet, so only one CPU is possible.
	}

	// 读取当前设置
	lock(&sched.lock)
	ret := int(gomaxprocs)
	unlock(&sched.lock)
	if n <= 0 || n == ret {
		// 返回当前设置
		return ret
	}

	stopTheWorldGC("GOMAXPROCS")

	// newprocs will be processed by startTheWorld
	// newprocs将由startTheWorld处理
	newprocs = int32(n)

	startTheWorldGC()
	return ret
}

// NumCPU returns the number of logical CPUs usable by the current process.
//
// The set of available CPUs is checked by querying the operating system
// at process startup. Changes to operating system CPU allocation after
// process startup are not reflected.

// NumCPU返回当前进程可用的逻辑CPU数量。
// 通过在进程启动时查询操作系统来检查一组可用的CPU。
// 进程启动后对操作系统CPU分配数的更改不会生效。
func NumCPU() int {
	return int(ncpu)
}

// NumCgoCall returns the number of cgo calls made by the current process.
// NumCgoCall返回当前进程进行的cgo调用次数。
func NumCgoCall() int64 {
	var n int64
	for mp := (*m)(atomic.Loadp(unsafe.Pointer(&allm))); mp != nil; mp = mp.alllink {
		n += int64(mp.ncgocall)
	}
	return n
}

// NumGoroutine returns the number of goroutines that currently exist.
// NumGoroutine返回当前存在的goroutine的数量。
func NumGoroutine() int {
	return int(gcount())
}

//go:linkname debug_modinfo runtime/debug.modinfo
func debug_modinfo() string {
	return modinfo
}
