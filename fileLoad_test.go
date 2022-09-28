package main

import (
	"strings"
	"testing"
)

func TestFileLoadPaths(t *testing.T) {
	AssertFileLoad(t, "fileLoad_test.go", "")
	AssertFileLoad(t, "fileXXX", "")
	AssertFileLoad(t, "testdata/fileXXX.bat", "")
	AssertFileLoad(t, "testdata/ibeacon.rle", "")
	AssertTailFileLoad(t, "testdata", "/grtest")
	// AssertFileLoad(t, "a/", "")
	// AssertFileLoad(t, "/a", "/")
	// AssertFileLoad(t, "/", "")
	// AssertFileLoad(t, "a/b", "a")
	// AssertFileLoad(t, "a/b/c", "a/b")
	// AssertFileLoad(t, "/a/b", "/a")
	// AssertFileLoad(t, "/a/b/c", "/a/b")

}

func AssertFileLoad(t *testing.T, path, exp string) {
	s, e := PathToParentPath(path)
	if e != nil {
		if e.Error() != exp {
			t.Errorf("PathToParent Failed. For Path '%s' Returned Error '%s' expected '%s'", path, s, exp)
			return
		}
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
