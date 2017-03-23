package filesystem

import (
	"errors"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"

	"github.com/goggle/flatten/osabstraction"
)

// Filesystem is the data type for a fake filesystem in the memory.
type Filesystem map[string]DummyFile

// Init initializes the filesystem by setting a root entry.
func (fs Filesystem) Init() {
	fs["/"] = DummyFile{Path: "/", IsDirectory: true}
}

// MkDir adds a directory dir to the filesystem. The parrent directories
// of dir do not need to be in the filesystem yet.
func (fs Filesystem) MkDir(dir string) error {
	cleanPath := filepath.Clean(dir)
	if strings.Index(cleanPath, "/") != 0 {
		return errors.New("invalid path: " + cleanPath)
	}
	_, alreadyExists := fs[cleanPath]
	if alreadyExists {
		return errors.New(dir + " already exists in filesystem")
	}
	chain := strings.Split(cleanPath, "/")
	currentPath := "/"
	for _, name := range chain {
		if name == "" {
			continue
		}
		currentPath += name
		file, exists := fs[currentPath]
		if exists && !file.IsDir() {
			return errors.New(currentPath + " is a file, not a directory")
		} else if !exists {
			break
		}
		currentPath += "/"
	}
	currentPath = "/"
	for _, name := range chain {
		if name == "" {
			continue
		}
		currentPath += name
		_, exists := fs[currentPath]
		if !exists {
			fs[currentPath] = DummyFile{Path: currentPath, IsDirectory: true}
		}
		currentPath += "/"
	}
	return nil
}

// CreateFile adds a file to the filesystem.
func (fs Filesystem) CreateFile(fpath string) error {
	cleanPath := filepath.Clean(fpath)
	if strings.Index(cleanPath, "/") != 0 {
		return errors.New("invalid path: " + cleanPath)
	}
	file := DummyFile{Path: cleanPath, IsDirectory: false}
	directory := file.Directory()
	entry, exists := fs[directory]
	if !exists {
		err := fs.MkDir(directory)
		if err != nil {
			return err
		}
	} else if !entry.IsDir() {
		return errors.New(directory + " is not a directory")
	}
	_, exists = fs[cleanPath]
	if exists {
		return errors.New(cleanPath + " already exists in file system")
	}
	fs[cleanPath] = file
	return nil
}

// RemoveDirectory removes an empty directory from the filesystem.
func (fs Filesystem) RemoveDirectory(path string) error {
	cleanPath := filepath.Clean(path)
	if cleanPath == "/" {
		return errors.New("cannot remove root directory")
	}
	entry, exists := fs[cleanPath]
	if !exists {
		return errors.New(cleanPath + " does not exist in file system")
	} else if !entry.IsDir() {
		return errors.New(cleanPath + " is not a directory")
	}
	for _, v := range fs {
		if v.Directory() == cleanPath {
			return errors.New(cleanPath + " is not empty!")
		}
	}
	delete(fs, cleanPath)
	return nil
}

// RemoveFile removes a regular file from the filesystem.
func (fs Filesystem) RemoveFile(path string) error {
	cleanPath := filepath.Clean(path)
	entry, exists := fs[cleanPath]
	if !exists {
		return errors.New(cleanPath + " does not exist in file system")
	} else if entry.IsDir() {
		return errors.New(cleanPath + " is a directory")
	}
	delete(fs, cleanPath)
	return nil
}

// Copy copies a file from source to destination on the filesystem.
func (fs Filesystem) Copy(source string, destination string) error {
	sourcePath := filepath.Clean(source)
	file, exists := fs[sourcePath]
	if !exists {
		return errors.New(sourcePath + " does not exist in file system")
	}
	if file.IsDir() {
		return fs.MkDir(destination)
	}
	return fs.CreateFile(destination)
}

// Move moves a file from source to destination on the filesystem.
func (fs Filesystem) Move(source string, destination string) error {
	err := fs.Copy(source, destination)
	if err != nil {
		return err
	}
	file, _ := fs[filepath.Clean(source)]
	if file.IsDir() {
		return fs.RemoveDirectory(source)
	}
	return fs.RemoveFile(source)

}

// Dirs returns a list of all the directories on the filesystem.
func (fs Filesystem) Dirs() []string {
	dirs := []string{}
	for _, v := range fs {
		if v.IsDir() {
			dirs = append(dirs, v.FullPath())
		}
	}
	return dirs
}

// RealFiles returns a list of all the regular files (not directories)
// on the filesystem.
func (fs Filesystem) RealFiles() []string {
	files := []string{}
	for _, v := range fs {
		if !v.IsDir() {
			files = append(files, v.FullPath())
		}
	}
	return files
}

// GetFiles returns all the regular files located at dir (if includeBaseFiles),
// and in the subdirectories of dir.
func (fs Filesystem) GetFiles(dir string, includeBaseFiles bool) ([]osabstraction.FileInfo, error) {
	files := []osabstraction.FileInfo{}
	dir = path.Clean(dir)
	for _, v := range fs {
		if !v.IsDir() {
			if v.Directory() == dir && includeBaseFiles {
				files = append(files, v)
			} else if v.Directory() != dir && strings.HasPrefix(v.FullPath(), dir) {
				files = append(files, v)
			}
		}
	}
	return files, nil
}

// GetDirectories returns all the directories in the dir subtree.
func (fs Filesystem) GetDirectories(dir string) ([]osabstraction.FileInfo, error) {
	files := []osabstraction.FileInfo{}
	dir = path.Clean(dir)
	for _, v := range fs {
		if v.IsDir() {
			if strings.HasPrefix(v.FullPath(), dir) && v.FullPath() != dir {
				files = append(files, v)
			}
		}
	}
	return files, nil
}

type byLevel []osabstraction.FileInfo

func (bl byLevel) Len() int {
	return len(bl)
}
func (bl byLevel) Swap(i, j int) {
	bl[i], bl[j] = bl[j], bl[i]
}
func (bl byLevel) Less(i, j int) bool {
	return bl[i].Level() < bl[j].Level()
}

// RemoveSubDirectories removes all the directories in the p subtree.
func (fs Filesystem) RemoveSubDirectories(p string) error {
	p = filepath.Clean(p)
	if !fs.IsDirectory(p) {
		return errors.New(p + " is not a directory")
	}
	dirs, err := fs.GetDirectories(p)
	if err != nil {
		return err
	}
	sort.Sort(byLevel(dirs))
	for i := len(dirs) - 1; i >= 0; i-- {
		err := fs.RemoveDirectory(dirs[i].FullPath())
		if err != nil {
			return err
		}
	}
	return nil
}

// AddFromRealFilesystem adds all the files and directories from
// a given path p to the simulated filesystem fs.
func (fs Filesystem) AddFromRealFilesystem(p string) error {
	p = path.Clean(p)
	err := filepath.Walk(p, func(path string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if fi.IsDir() {
			// FIXME: ignoring errors:
			fs.MkDir(path)
		} else {
			// FIXME: ignoring errors:
			fs.CreateFile(path)
		}
		return nil
	})
	return err
}

// IsRegularFile returns true if a file p is a regular file
// on the filesystem fs, otherwise false
func (fs Filesystem) IsRegularFile(p string) bool {
	df, exists := fs[path.Clean(p)]
	if exists && !df.IsDir() {
		return true
	}
	return false
}

// IsDirectory returns true if a file p is a directory
// in the filesystem fs, otherwise false.
func (fs Filesystem) IsDirectory(p string) bool {
	df, exists := fs[path.Clean(p)]
	if exists && df.IsDir() {
		return true
	}
	return false
}

// Exists returns true if a file p exists in the filesystem
// fs, otherwise false.
func (fs Filesystem) Exists(p string) bool {
	if fs.IsRegularFile(p) || fs.IsDirectory(p) {
		return true
	}
	return false
}

// Equal tests, if two filesystems have exactly the same structure.
func (fs Filesystem) Equal(f Filesystem) bool {
	if len(fs) != len(f) {
		return false
	}
	for k, v := range fs {
		vv, ok := f[k]
		if !ok {
			return false
		}
		// NOTE: in our file system (which is only used for simulation purposes),
		// two files are considered as equal, if they have the same path and
		// are of the same type (i.e. either regular file or directory)
		if v.FullPath() != vv.FullPath() || v.IsDir() != vv.IsDir() {
			return false
		}
	}
	return true
}

// DummyFile represents a file on a simulated filesystem.
type DummyFile struct {
	Path        string
	IsDirectory bool
}

// IsDir checks if df is a directory.
func (df DummyFile) IsDir() bool {
	return df.IsDirectory
}

// FullPath returns the full path of df.
func (df DummyFile) FullPath() string {
	return filepath.Clean(df.Path)
}

// Name returns the name of df.
func (df DummyFile) Name() string {
	return filepath.Base(df.Path)
}

// Directory returns the directory in which the dummy file df
// is located.
func (df DummyFile) Directory() string {
	clean := filepath.Clean(df.Path)
	j := strings.LastIndex(clean, "/")
	dir := clean[:j]
	if dir == "" {
		return "/"
	}
	return dir
}

// Ext returns the extension of df.
func (df DummyFile) Ext() string {
	return filepath.Ext(df.Path)
}

// BaseName returns the filename of df without its extension.
func (df DummyFile) BaseName() string {
	filename := df.Name()
	ext := df.Ext()
	j := strings.LastIndex(filename, ext)
	return filename[:j]
}

// Level returns the tree depth of the branch, on which
// df is located in the simulated filesystem.
func (df DummyFile) Level() int {
	clean := filepath.Clean(df.Path)
	if clean == "" || clean == "/" {
		return 0
	}
	return strings.Count(clean, "/")
}
