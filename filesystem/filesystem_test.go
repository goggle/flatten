package filesystem_test

import (
	"testing"

	"github.com/goggle/flatten/filesystem"
)

func matchLists(l1, l2 []string) bool {
	if len(l1) != len(l2) {
		return false
	}
	for i, _ := range l1 {
		found := false
		for j, _ := range l2 {
			if l1[i] == l2[j] {
				found = true
				break
			}
		}
		if found {
			continue
		} else {
			return false
		}
	}
	return true
}

func noErrorExpected(t *testing.T, err error) {
	if err != nil {
		t.Errorf("No error expected, got %v", err)
	}
}

func TestFilesystemOne(t *testing.T) {
	fs := make(filesystem.Filesystem)

	// Init file system
	fs.Init()
	dirs := fs.Dirs()
	expected := []string{"/"}
	if !matchLists(dirs, expected) {
		t.Errorf("Init(): Expected to have initial file system, containing only root directory, got %v", dirs)
	}

	// Add nested folder to file system
	err := fs.MkDir("/home/goggle/Downloads/test")
	// fmt.Println(fs)

	noErrorExpected(t, err)
	dirs = fs.Dirs()
	expected = []string{
		"/",
		"/home",
		"/home/goggle",
		"/home/goggle/Downloads",
		"/home/goggle/Downloads/test"}
	if !matchLists(dirs, expected) {
		t.Errorf("MkDir: Expected %v, got %v", expected, dirs)
	}

	// Add some files
	err = fs.CreateFile("/home/goggle/Downloads/movie.mp4")
	noErrorExpected(t, err)
	err = fs.CreateFile("/home/goggle/Downloads/test/lichess.tar.gz")
	noErrorExpected(t, err)
	files := fs.RealFiles()
	expected = []string{
		"/home/goggle/Downloads/movie.mp4",
		"/home/goggle/Downloads/test/lichess.tar.gz"}
	if !matchLists(files, expected) {
		t.Errorf("CreateFile: Expected %v, got %v", expected, files)
	}
	dirs = fs.Dirs()
	expected = []string{
		"/",
		"/home",
		"/home/goggle",
		"/home/goggle/Downloads",
		"/home/goggle/Downloads/test"}
	if !matchLists(dirs, expected) {
		t.Errorf("CreateFile (directory check): Expected %v, got %v", expected, dirs)
	}

	// Copy a file:
	err = fs.Copy("/home/goggle/Downloads/movie.mp4", "/home/goggle/Downloads/test/Bang.mp4")
	noErrorExpected(t, err)
	files = fs.RealFiles()
	expected = []string{
		"/home/goggle/Downloads/movie.mp4",
		"/home/goggle/Downloads/test/lichess.tar.gz",
		"/home/goggle/Downloads/test/Bang.mp4"}
	if !matchLists(files, expected) {
		t.Errorf("Copy: Expected %v, got %v", expected, files)
	}
	dirs = fs.Dirs()
	expected = []string{
		"/",
		"/home",
		"/home/goggle",
		"/home/goggle/Downloads",
		"/home/goggle/Downloads/test"}
	if !matchLists(dirs, expected) {
		t.Errorf("Copy (directory check): Expected %v, got %v", expected, dirs)
	}

	// Move a file:
	err = fs.Move("/home/goggle/Downloads/test/lichess.tar.gz", "/home/goggle/Downloads/chess.tar.gz")
	noErrorExpected(t, err)
	files = fs.RealFiles()
	expected = []string{
		"/home/goggle/Downloads/movie.mp4",
		"/home/goggle/Downloads/chess.tar.gz",
		"/home/goggle/Downloads/test/Bang.mp4"}
	if !matchLists(files, expected) {
		t.Errorf("Move: Expected %v, got %v", expected, files)
	}
	dirs = fs.Dirs()
	expected = []string{
		"/",
		"/home",
		"/home/goggle",
		"/home/goggle/Downloads",
		"/home/goggle/Downloads/test"}
	if !matchLists(dirs, expected) {
		t.Errorf("Move (directory check): Expected %v, got %v", expected, dirs)
	}

	// Remove the files:
	err = fs.RemoveFile("/home/goggle/Downloads/movie.mp4")
	noErrorExpected(t, err)
	err = fs.RemoveFile("/home/goggle/Downloads/chess.tar.gz")
	noErrorExpected(t, err)
	err = fs.RemoveFile("/home/goggle/Downloads/test/Bang.mp4")
	noErrorExpected(t, err)
	files = fs.RealFiles()
	expected = []string{}
	if !matchLists(files, expected) {
		t.Errorf("RemoveFile: Expected %v, got %v", expected, files)
	}
	dirs = fs.Dirs()
	expected = []string{
		"/",
		"/home",
		"/home/goggle",
		"/home/goggle/Downloads",
		"/home/goggle/Downloads/test"}
	if !matchLists(dirs, expected) {
		t.Errorf("RemoveFile (directory check): Expected %v, got %v", expected, dirs)
	}

	// Remove the directories:
	err = fs.RemoveDirectory("/home/goggle/Downloads/test")
	noErrorExpected(t, err)
	err = fs.RemoveDirectory("/home/goggle/Downloads")
	noErrorExpected(t, err)
	err = fs.RemoveDirectory("/home/goggle")
	noErrorExpected(t, err)
	err = fs.RemoveDirectory("/home")
	noErrorExpected(t, err)
	dirs = fs.Dirs()
	expected = []string{"/"}
	// if !matchLists(dirs, expected) {
	// 	t.Errorf("RemoveDirectoy: Expected %v, got %v", expected, dirs)
	// }
}

func TestDummyFile(t *testing.T) {
	df := filesystem.DummyFile{Path: "/home/goggle/test/my_song.flac", IsDirectory: false}

	isDir := df.IsDir()
	expectedIsDir := false
	if isDir != expectedIsDir {
		t.Errorf("IsDir: Expected %v, got %v", expectedIsDir, isDir)
	}

	fullpath := df.FullPath()
	expectedFullpath := "/home/goggle/test/my_song.flac"
	if fullpath != expectedFullpath {
		t.Errorf("Fullpath: Expected %v, got %v", expectedFullpath, fullpath)
	}

	name := df.Name()
	expectedName := "my_song.flac"
	if name != expectedName {
		t.Errorf("Name: Expected %v, got %v", expectedName, name)
	}

	directory := df.Directory()
	expectedDirectory := "/home/goggle/test"
	if directory != expectedDirectory {
		t.Errorf("Directory: Expected %v, got %v", expectedDirectory, directory)
	}

	ext := df.Ext()
	expectedExt := ".flac"
	if ext != expectedExt {
		t.Errorf("Ext: Expected %v, got %v", expectedExt, ext)
	}

	basename := df.BaseName()
	expectedBaseName := "my_song"
	if basename != expectedBaseName {
		t.Errorf("BaseName: Expected %v, got %v", expectedBaseName, basename)
	}

	level := df.Level()
	expectedLevel := 4
	if level != expectedLevel {
		t.Errorf("Level: Expected %v, got %v", expectedLevel, level)
	}
}
