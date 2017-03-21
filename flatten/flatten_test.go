package flatten

import (
	"testing"

	"github.com/goggle/flatten/filesystem"
)

func TestGenerateFilename(t *testing.T) {
	type generateFilenameTestCase struct {
		n string
		i int
		l int
	}
	test_cases := []generateFilenameTestCase{
		{"hello", 2, 3},
		{"myFile", 24, 5},
		{"movie.mp4", 1, 1},
		{"song_23.flac", 99, 3},
	}
	expected_results := []string{
		"hello_002",
		"myFile_00024",
		"movie_1.mp4",
		"song_23_099.flac",
	}
	for i, tc := range test_cases {
		res := generateFilename(tc.n, tc.i, tc.l)
		if res != expected_results[i] {
			t.Errorf("generateFilename: expected %v, got %v", expected_results[i], res)
		}
	}
}

func TestCountFileNames(t *testing.T) {
	fs := filesystem.Filesystem{}
	fs.Init()
	fs.MkDir("/tmp/a")
	fs.MkDir("/tmp/b")
	fs.MkDir("/tmp/c")
	fs.CreateFile("/tmp/a/hello")
	fs.CreateFile("/tmp/b/hello")
	fs.CreateFile("/tmp/c/hello")
	fs.CreateFile("/tmp/c/nf.txt")
	files, _ := fs.GetFiles("/tmp", false)
	m := countFileNames(files)

	expectedLength := 2
	if len(m) != expectedLength {
		t.Errorf("countFileNames: Expected map length of %v, got %v", expectedLength, len(m))
	}

	expectedMap := map[string]int{}
	expectedMap["hello"] = 3
	expectedMap["nf.txt"] = 1
	for k, v := range expectedMap {
		tv, tok := m[k]
		if !tok {
			t.Errorf("countFileNames: Expected %v to be found in map", k)
		} else if tv != v {
			t.Errorf("countFileNames(key = %v): Expected %v, got %v.", k, v, tv)
		}
	}
}
