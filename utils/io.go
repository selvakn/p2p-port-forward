package utils

import (
	"io"
)

func Sync(source1 io.ReadWriteCloser, source2 io.ReadWriteCloser) {
	defer source1.Close()
	defer source2.Close()

	go io.Copy(source1, source2)
	io.Copy(source2, source1)
}