package osabstraction

import (
	"os"
	"path"
	"path/filepath"
	"strings"
)

type OSWrapper interface {
	// Copy(src, dst FileInfo) error
	// Move(src, dst FileInfo) error
	Copy(src, dst string) error
	Move(src, dst string) error
	Stat(name string) (FileInfo, error)
	IsNotExist(err error) bool
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
