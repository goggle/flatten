package osabstraction

import (
	"errors"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"
)

// OSWrapper contains all the OS functions, which we want
// to be able to use.
type OSWrapper interface {
	Copy(src, dst string) error
	Move(src, dst string) error
	GetFiles(dir string, includeBaseFiles bool) ([]FileInfo, error)
	GetDirectories(dir string) ([]FileInfo, error)
	IsRegularFile(p string) bool
	IsDirectory(p string) bool
	Exists(p string) bool
	RemoveSubDirectories(p string) error
}

// FileInfo represents all the relevant information about a file
// which might be on a real file system or a fake one.
type FileInfo interface {
	IsDir() bool
	FullPath() string
	Name() string
	Directory() string
	Ext() string
	BaseName() string
	Level() int
}

// File is the path of a file.
type File string

// IsDir checks, if a file f is a directory.
func (f File) IsDir() bool {
	fp := f.FullPath()
	fi, err := os.Stat(fp)
	if err != nil {
		return false
	}
	mode := fi.Mode()
	if mode.IsDir() {
		return true
	}
	return false
}

// FullPath returns the path of a file f as a string.
func (f File) FullPath() string {
	return path.Clean(string(f))
}

// Name returns the filename.
func (f File) Name() string {
	return filepath.Base(f.FullPath())
}

// Directory returns the full path of the directory, in
// which the file lays in.
func (f File) Directory() string {
	j := strings.LastIndex(f.FullPath(), "/")
	dir := f.FullPath()[:j]
	if dir == "" {
		return "/"
	}
	return dir
}

// Ext returns the file extension if there is one, otherwise
// an empty string.
func (f File) Ext() string {
	return filepath.Ext(f.FullPath())
}

// BaseName returns the filename without the file extension.
func (f File) BaseName() string {
	filename := f.Name()
	ext := f.Ext()
	j := strings.LastIndex(filename, ext)
	return filename[:j]
}

// Level returns the depth of the file in the filesystem tree.
// The root path has level 0.
func (f File) Level() int {
	fp := f.FullPath()
	if fp == "" || fp == "/" {
		return 0
	}
	return strings.Count(fp, "/")
}

// RealOS is the data type for the operating system if we want
// to operate on the real filesystem.
type RealOS struct{}

// Copy copies a file src to dst on the real filesystem.
func (ros RealOS) Copy(src, dst string) error {
	if !ros.Exists(src) {
		return errors.New(src + " does not exist in file system")
	} else if ros.IsDirectory(src) {
		return errors.New(src + " is a directory")
	}
	if ros.Exists(dst) {
		return errors.New(dst + " already exists in file system")
	}
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}

	err = out.Sync()
	if err != nil {
		return err
	}

	cerr := out.Close()
	return cerr
}

// Move moves a file src to dst on the real filesystem.
func (ros RealOS) Move(src, dst string) error {
	return os.Rename(src, dst)
}

// GetFiles scans the underlying tree beginning at dir and returns the
// list of files found in all the subdirectories. Directories are not
// considered as files here. The option includeBaseFiles indicates,
// if the files which are directly located in dir (not in a subdirectory
// of dir) should also be added to list or not.
func (ros RealOS) GetFiles(dir string, includeBaseFiles bool) ([]FileInfo, error) {
	files := []FileInfo{}
	err := filepath.Walk(dir, func(p string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			if !includeBaseFiles {
				if path.Dir(p) != path.Clean(dir) {
					files = append(files, File(p))
				}
			} else {
				files = append(files, File(p))
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return files, nil
}

// GetDirectories does the same thing as GetFiles, but only for
// directories.
func (ros RealOS) GetDirectories(dir string) ([]FileInfo, error) {
	dirs := []FileInfo{}
	err := filepath.Walk(dir, func(p string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			if path.Clean(p) != path.Clean(dir) {
				dirs = append(dirs, File(p))
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return dirs, nil
}

// IsRegularFile returns true if a file at path p
// is a file but not a directory, otherwise false
func (ros RealOS) IsRegularFile(p string) bool {
	if stat, err := os.Stat(path.Clean(p)); err == nil {
		if !stat.IsDir() {
			return true
		}
	}
	return false
}

// IsDirectory returns true if a file located at path p
// is a directory, otherwise false.
func (ros RealOS) IsDirectory(p string) bool {
	if stat, err := os.Stat(path.Clean(p)); err == nil {
		if stat.IsDir() {
			return true
		}
	}
	return false
}

// Exists returns true if a file located at path exists
// on the filesystem, otherwise false.
func (ros RealOS) Exists(p string) bool {
	if _, err := os.Stat(path.Clean(p)); err == nil {
		return true
	}
	return false
}

// RemoveSubDirectories recurively removes all the directories
// in the parrent directory p. An error is returned, if a
// a subdirectory could not be removed, probably because it
// is not empty.
func (ros RealOS) RemoveSubDirectories(p string) error {
	for {
		subDirectories, err := ros.GetDirectories(p)
		if err != nil {
			return err
		}
		if len(subDirectories) == 0 {
			break
		}
		changed := false
		for _, sd := range subDirectories {
			ssd, err := ros.GetDirectories(sd.FullPath())
			if err != nil {
				return err
			}
			if len(ssd) == 0 {
				os.Remove(sd.FullPath())
				changed = true
			}
		}
		if !changed {
			return errors.New("could not remove all the subdirectories in " + path.Clean(p))
		}
	}
	return nil
}
