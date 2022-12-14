package main

import (
	"fmt"
	"strings"
	"testing"
)

func TestFileSave(t *testing.T) {
	rle, err := NewRleFile("testdata/rats.rle")
	if err != nil {
		t.Errorf("RLE File load failed. %e", err)
	}

	rleSave := NewRLESave("testdata/ab", rle.coords, "OWNER", "DESC")
	s := rleSave.SaveFileContent()
	fmt.Println(s)
}

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
	enc, w, h := rle.Encode()
	if rle.encoded != enc {
		t.Errorf("RLE File Encode failed. \n%s\n%s", rle.encoded, enc)
	}
	if w != 12 {
		t.Errorf("RLE File Encode failed. Expected width %d Actual Width %d", 12, w)
	}
	if h != 11 {
		t.Errorf("RLE File Encode failed. Expected height %d Actual Width %d", 11, h)
	}

	saveRle := NewRLESave("RLESave", rle.coords, "owner", "desc")
	if rle.decoded != saveRle.decoded {
		t.Errorf("RLE File Encode failed. Decoded expected \n%s Actual \n%s", rle.decoded, saveRle.decoded)
	}
	if rle.encoded != saveRle.encoded {
		t.Errorf("RLE File Encode failed. Encoded expected \n%s Actual \n%s", rle.encoded, saveRle.encoded)
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
