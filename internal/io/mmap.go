package io

import (
	"os"

	"github.com/edsrzf/mmap-go"
)

func MapRO(f *os.File) mmap.MMap {
	if m, err := mmap.Map(f, mmap.RDONLY, 0); err != nil {
		panic(err)
	} else {
		return m
	}
}

func MapRW(f *os.File) mmap.MMap {
	if m, err := mmap.Map(f, mmap.RDWR, 0); err != nil {
		panic(err)
	} else {
		return m
	}
}

func Unmap(m mmap.MMap) {
	if err := m.Unmap(); err != nil {
		panic(err)
	}
}

func Open(filename string) *os.File {
	f, err := os.Open(filename)
	if err != nil {
		panic(err)
	} else {
		return f
	}
}

func Create(filename string, size int64) *os.File {
	f, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	err = f.Truncate(size)
	if err != nil {
		panic(err)
	}
	return f
}
