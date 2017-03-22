package osabstraction

import (
	"errors"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"
)

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

type FileInfo interface {
	IsDir() bool
	FullPath() string
	Name() string
	Directory() string
	Ext() string
	BaseName() string
	Level() int
}

type File string

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

func (f File) FullPath() string {
	return path.Clean(string(f))
}

func (f File) Name() string {
	return filepath.Base(f.FullPath())
}

func (f File) Directory() string {
	j := strings.LastIndex(f.FullPath(), "/")
	dir := f.FullPath()[:j]
	if dir == "" {
		return "/"
	}
	return dir
}

func (f File) Ext() string {
	return filepath.Ext(f.FullPath())
}

func (f File) BaseName() string {
	filename := f.Name()
	ext := f.Ext()
	j := strings.LastIndex(filename, ext)
	return filename[:j]
}

func (f File) Level() int {
	fp := f.FullPath()
	if fp == "" || fp == "/" {
		return 0
	}
	return strings.Count(fp, "/")
}

type RealOS struct{}

func (ros RealOS) Copy(src, dst string) error {
	sourceFile := File(src)
	stat, err := os.Stat(src)
	if err != nil {
		return errors.New(src + " does not exist in file system")
	}
	if stat.IsDir() {
		return errors.New(src + " is a directory")
	}
	stat, err = os.Stat(dst)
	if err == nil {
		return errors.New(dst + " already exists in file system")
	}
	in, err := os.Open(sourceFile.FullPath())
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(in, out)
	cerr := out.Close()
	if err != nil {
		return err
	}
	return cerr
}

func (ros RealOS) Move(src, dst string) error {
	return os.Rename(src, dst)
}

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

func (ros RealOS) IsRegularFile(p string) bool {
	if stat, err := os.Stat(path.Clean(p)); err == nil {
		if !stat.IsDir() {
			return true
		}
	}
	return false
}

func (ros RealOS) IsDirectory(p string) bool {
	if stat, err := os.Stat(path.Clean(p)); err == nil {
		if stat.IsDir() {
			return true
		}
	}
	return false
}

func (ros RealOS) Exists(p string) bool {
	if _, err := os.Stat(path.Clean(p)); err == nil {
		return true
	}
	return false
}

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
