// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package unix_test

import (
	"bytes"
	"io/ioutil"
	"os"
	"path"
	"testing"

	"golang.org/x/sys/unix"
)

// stringsFromByteSlice converts a sequence of attributes to a []string.
// On Darwin, each entry is a NULL-terminated string.
func stringsFromByteSlice(buf []byte) []string {
	var result []string
	off := 0
	for i, b := range buf {
		if b == 0 {
			result = append(result, string(buf[off:i]))
			off = i + 1
		}
	}
	return result
}

func TestClonefile(t *testing.T) {
	file, err := ioutil.TempFile("", "TestCloneFile")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(file.Name())

	data := []byte("This is a test")
	file.Write(data)
	file.Close()

	clonedName := file.Name() + "-cloned"
	err = unix.Clonefile(file.Name(), clonedName, 0)
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(clonedName)

	clonedData, err := ioutil.ReadFile(clonedName)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(data, clonedData) {
		t.Fatal(err)
	}
}

func TestClonefileatWithCwd(t *testing.T) {
	file, err := ioutil.TempFile("", "TestCloneFileat")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(file.Name())

	data := []byte("This is a test")
	file.Write(data)
	file.Close()

	clonedName := file.Name() + "-cloned"
	err = unix.Clonefileat(unix.AT_FDCWD, file.Name(), unix.AT_FDCWD, clonedName, 0)
	if err != nil {
		t.Fatal(err)
	}

	defer os.Remove(clonedName)

	clonedData, err := ioutil.ReadFile(clonedName)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(data, clonedData) {
		t.Fatal(err)
	}
}

func TestClonefileatWithRelativePaths(t *testing.T) {
	srcDir, err := ioutil.TempDir("", "src")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(srcDir)

	dstDir, err := ioutil.TempDir("", "dest")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dstDir)

	srcFd, err := unix.Open(srcDir, unix.O_RDONLY|unix.O_DIRECTORY, 0)
	if err != nil {
		t.Fatal(err)
	}
	defer unix.Close(srcFd)

	dstFd, err := unix.Open(dstDir, unix.O_RDONLY|unix.O_DIRECTORY, 0)
	if err != nil {
		t.Fatal(err)
	}
	defer unix.Close(dstFd)

	srcFile, err := ioutil.TempFile(srcDir, "TestCloneFileat")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(srcFile.Name())

	dstFile, err := ioutil.TempFile(dstDir, "TestCloneFileat")
	if err != nil {
		t.Fatal(err)
	}
	os.Remove(dstFile.Name())

	data := []byte("This is a test")
	srcFile.Write(data)
	srcFile.Close()

	src := path.Base(srcFile.Name())
	dst := path.Base(dstFile.Name())
	err = unix.Clonefileat(srcFd, src, dstFd, dst, 0)
	if err != nil {
		t.Fatal(err)
	}

	clonedData, err := ioutil.ReadFile(dstFile.Name())
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(data, clonedData) {
		t.Fatal(err)
	}
}

func TestFclonefileat(t *testing.T) {
	file, err := ioutil.TempFile("", "TestCloneFile")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(file.Name())

	data := []byte("This is a test")
	file.Write(data)
	file.Close()

	fd, err := unix.Open(file.Name(), unix.O_RDONLY, 0)
	if err != nil {
		t.Fatal(err)
	}
	defer unix.Close(fd)

	dstFile, err := ioutil.TempFile("", "TestFcloneFileat")
	if err != nil {
		t.Fatal(err)
	}
	os.Remove(dstFile.Name())

	err = unix.Fclonefileat(fd, unix.AT_FDCWD, dstFile.Name(), 0)
	if err != nil {
		t.Fatal(err)
	}

	clonedData, err := ioutil.ReadFile(dstFile.Name())
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(data, clonedData) {
		t.Fatal(err)
	}
}
