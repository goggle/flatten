package flatten

import (
	"fmt"
	"testing"

	"github.com/goggle/flatten/filesystem"
)

func TestGenerateFilename(t *testing.T) {
	type generateFilenameTestCase struct {
		n string
		i int
		l int
	}
	testCases := []generateFilenameTestCase{
		{"hello", 2, 3},
		{"myFile", 24, 5},
		{"movie.mp4", 1, 1},
		{"song_23.flac", 99, 3},
	}
	expectedResults := []string{
		"hello_002",
		"myFile_00024",
		"movie_1.mp4",
		"song_23_099.flac",
	}
	for i, tc := range testCases {
		res := generateFilename(tc.n, tc.i, tc.l)
		if res != expectedResults[i] {
			t.Errorf("generateFilename: expected %v, got %v", expectedResults[i], res)
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

func TestEvaluateAppendixLength(t *testing.T) {
	fs := filesystem.Filesystem{}
	fs.Init()
	fs.MkDir("/tmp/a")
	fs.CreateFile("/tmp/a/hello")

	expected := 0
	result := evaluateAppendixLength("/tmp", "hello", 1, fs)
	if expected != result {
		t.Errorf("evaluateAppendixLength: expected %v, got %v", expected, result)
	}

	fs.MkDir("/tmp/b")
	fs.CreateFile("/tmp/b/hello")
	expected = 1
	result = evaluateAppendixLength("/tmp", "hello", 2, fs)
	if expected != result {
		t.Errorf("evaluateAppendixLength: expected %v, got %v", expected, result)
	}

	for i := 0; i < 15; i++ {
		stri := fmt.Sprintf("%v", i)
		fs.MkDir("/tmp/c" + stri)
		fs.CreateFile("/tmp/c" + stri + "/hello")
	}
	expected = 2
	result = evaluateAppendixLength("/tmp", "hello", 17, fs)
	if expected != result {
		t.Errorf("evaluateAppendixLength: expected %v, got %v", expected, result)
	}

	for i := 0; i < 83; i++ {
		stri := fmt.Sprintf("%v", i)
		fs.MkDir("/tmp/d" + stri)
		fs.CreateFile("/tmp/d" + stri + "/hello")
	}
	expected = 3
	result = evaluateAppendixLength("/tmp", "hello", 100, fs)
	if expected != result {
		t.Errorf("evaluateAppendixLength: expected %v, got %v", expected, result)
	}

	fs.CreateFile("/tmp/hello_001")
	expected = 4
	result = evaluateAppendixLength("/tmp", "hello", 100, fs)
	if expected != result {
		t.Errorf("evaluateAppendixLength: expected %v, got %v", expected, result)
	}

	fs.CreateFile("/tmp/hello_0100")
	expected = 5
	result = evaluateAppendixLength("/tmp", "hello", 100, fs)
	if expected != result {
		t.Errorf("evaluateAppendixLength: expected %v, got %v", expected, result)
	}
}

func TestFlattenOne(t *testing.T) {
	// Test copying with different source and destination:
	fs := filesystem.Filesystem{}
	fs.Init()
	fs.MkDir("/tmp/a")
	fs.MkDir("/tmp/a/aa")
	fs.MkDir("/tmp/b")
	fs.MkDir("/tmp/c")
	fs.MkDir("/home/goggle/Downloads")

	fs.CreateFile("/home/goggle/Downloads/hello.txt")
	fs.CreateFile("/tmp/a/hello.txt")
	fs.CreateFile("/tmp/b/hello.txt")
	fs.CreateFile("/tmp/c/world.zip")
	fs.CreateFile("/tmp/a/aa/testFile.rar")

	err := Flatten(fs["/tmp"], fs["/home/goggle/Downloads"], fs, true, false)
	if err != nil {
		t.Errorf("Flatten: no error expected, got %v", err)
	}

	expectedFs := filesystem.Filesystem{}
	expectedFs.Init()
	expectedFs.MkDir("/tmp/a")
	expectedFs.MkDir("/tmp/a/aa")
	expectedFs.MkDir("/tmp/b")
	expectedFs.MkDir("/tmp/c")
	expectedFs.MkDir("/home/goggle/Downloads")

	expectedFs.CreateFile("/home/goggle/Downloads/hello.txt")
	expectedFs.CreateFile("/tmp/a/hello.txt")
	expectedFs.CreateFile("/tmp/b/hello.txt")
	expectedFs.CreateFile("/tmp/c/world.zip")
	expectedFs.CreateFile("/tmp/a/aa/testFile.rar")

	expectedFs.CreateFile("/home/goggle/Downloads/hello_1.txt")
	expectedFs.CreateFile("/home/goggle/Downloads/hello_2.txt")
	expectedFs.CreateFile("/home/goggle/Downloads/world.zip")
	expectedFs.CreateFile("/home/goggle/Downloads/testFile.rar")

	if !fs.Equal(expectedFs) {
		t.Errorf("Flatten: expected %v, got %v", expectedFs, fs)
	}

	// Test moving with different source and destination
	fs = filesystem.Filesystem{}
	fs.Init()
	fs.MkDir("/tmp/a")
	fs.MkDir("/tmp/a/aa")
	fs.MkDir("/tmp/b")
	fs.MkDir("/tmp/c")
	fs.MkDir("/home/goggle/Downloads")

	fs.CreateFile("/home/goggle/Downloads/hello.txt")
	fs.CreateFile("/tmp/a/hello.txt")
	fs.CreateFile("/tmp/b/hello.txt")
	fs.CreateFile("/tmp/c/world.zip")
	fs.CreateFile("/tmp/a/aa/testFile.rar")

	err = Flatten(fs["/tmp"], fs["/home/goggle/Downloads"], fs, false, false)
	if err != nil {
		t.Errorf("Flatten: no error expected, got %v", err)
	}

	expectedFs = filesystem.Filesystem{}
	expectedFs.Init()
	expectedFs.MkDir("/home/goggle/Downloads")
	expectedFs.MkDir("/tmp")

	expectedFs.CreateFile("/home/goggle/Downloads/hello.txt")
	expectedFs.CreateFile("/home/goggle/Downloads/hello_1.txt")
	expectedFs.CreateFile("/home/goggle/Downloads/hello_2.txt")
	expectedFs.CreateFile("/home/goggle/Downloads/world.zip")
	expectedFs.CreateFile("/home/goggle/Downloads/testFile.rar")

	if !fs.Equal(expectedFs) {
		t.Errorf("Flatten: expected %v, got %v", expectedFs, fs)
	}
}

func TestFlattenTwo(t *testing.T) {
	fs := filesystem.Filesystem{}
	fs.Init()
	fs.MkDir("/home/goggle/Downloads/songs")
	for i := 1; i <= 222; i++ {
		stri := fmt.Sprintf("%v", i)
		fs.MkDir("/home/goggle/Downloads/songs/tmp" + stri)
		fs.CreateFile("/home/goggle/Downloads/songs/tmp" + stri + "/song.flac")
	}
	err := Flatten(fs["/home/goggle/Downloads"], fs["/home/goggle/Downloads"], fs, false, false)
	if err != nil {
		t.Errorf("Flatten: no error expected, got %v", err)
	}

	expectedFs := filesystem.Filesystem{}
	expectedFs.Init()
	expectedFs.MkDir("/home/goggle/Downloads")
	for i := 1; i <= 222; i++ {
		fileStr := fmt.Sprintf("/home/goggle/Downloads/song_%03v.flac", i)
		expectedFs.CreateFile(fileStr)
	}

	if !fs.Equal(expectedFs) {
		t.Errorf("Flatten: expected %v, got %v", expectedFs, fs)
	}
}

func TestFlattenThree(t *testing.T) {
	fs := filesystem.Filesystem{}
	fs.Init()
	fs.MkDir("/tmp/a")
	fs.CreateFile("/tmp/blubb.txt")
	fs.CreateFile("/tmp/a/blubb.txt")

	err := Flatten(fs["/tmp"], fs["/tmp"], fs, true, true)
	if err != nil {
		t.Errorf("Flatten: no error expected, got %v", err)
	}

	expectedFs := filesystem.Filesystem{}
	expectedFs.Init()
	expectedFs.MkDir("/tmp/a")
	expectedFs.CreateFile("/tmp/blubb.txt")
	expectedFs.CreateFile("/tmp/blubb_1.txt")
	expectedFs.CreateFile("/tmp/blubb_2.txt")
	expectedFs.CreateFile("/tmp/a/blubb.txt")
	if !fs.Equal(expectedFs) {
		t.Errorf("Flatten: expected %v, got %v", expectedFs, fs)
	}
}

func TestFlattenFour(t *testing.T) {
	fs := filesystem.Filesystem{}
	fs.Init()
	fs.MkDir("/tmp/a")
	fs.CreateFile("/tmp/blubb.txt")
	fs.CreateFile("/tmp/a/blubb.txt")
	fs.CreateFile("/tmp/mail.py")

	err := Flatten(fs["/tmp"], fs["/tmp"], fs, false, true)
	if err != nil {
		t.Errorf("Flatten: no error expected, got %v", err)
	}

	// This behavior (moving files including the base files with identical
	// source and destination path) might look weird at first sight.
	// This might get changed in the future...
	expectedFs := filesystem.Filesystem{}
	expectedFs.Init()
	expectedFs.MkDir("/tmp")
	expectedFs.CreateFile("/tmp/blubb_1.txt")
	expectedFs.CreateFile("/tmp/blubb_2.txt")
	expectedFs.CreateFile("/tmp/mail_1.py")
	if !fs.Equal(expectedFs) {
		t.Errorf("Flatten: expected %v, got %v", expectedFs, fs)
	}

}
