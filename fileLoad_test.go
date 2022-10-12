package main

import (
	"strings"
	"testing"
)

func TestFileLoadPaths(t *testing.T) {
	AssertFileLoad(t, "fileLoad_test.go", "is not a dir")
	AssertFileLoad(t, "fileXXX", "no such file or directory")
	AssertFileLoad(t, "testdata/fileXXX.bat", "no such file or directory")
	AssertFileLoad(t, "testdata/ibeacon.rle", "is not a dir")
	AssertTailFileLoad(t, "testdata", "/grtest")
}
func TestFileEncode(t *testing.T) {
	rle, err := NewRleFile("testdata/rats.rle")
	if err != nil {
		t.Errorf("RLE File load failed. %e", err)
	}
	if rle.encoded != rle.Encode() {
		t.Errorf("RLE File Encode failed. \n%s\n%s", rle.encoded, rle.Encode())
	}

}

func AssertFileLoad(t *testing.T, path, exp string) {
	s, e := PathToParentPath(path)
	if e != nil {
		if !strings.HasSuffix(e.Error(), exp) {
			t.Errorf("PathToParent Failed. For Path '%s' Returned Error '%s' expected '%s'", path, e.Error(), exp)
		}
		return
	}
	if s != exp {
		t.Errorf("PathToParent Failed. For Path '%s' Returned '%s' expected '%s'", path, s, exp)
	}
}

func AssertTailFileLoad(t *testing.T, path, exp string) {
	s, e := PathToParentPath(path)
	if e != nil {
		if e.Error() != exp {
			t.Errorf("PathToParent Failed. For Path '%s' Returned Error '%s' expected '%s'", path, s, exp)
			return
		}
	}
	if !strings.HasSuffix(s, exp) {
		t.Errorf("PathToParent Failed. For Path '%s' Returned '%s' expected suffix '%s'", path, s, exp)
	}
}
