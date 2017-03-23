package filesystem

import (
	"errors"
	"sort"
	"strings"

	"github.com/goggle/flatten/osabstraction"
)

type Tree struct {
	node     osabstraction.FileInfo
	children []*Tree
}

func (t *Tree) Init(fi osabstraction.FileInfo) error {
	if t.node == nil {
		t.node = fi
		t.children = make([]*Tree, 0)
		return nil
	}
	return errors.New("tree has already been initialized")
}

func (t *Tree) InsertSuccessor(fi osabstraction.FileInfo) error {
	if t.node == nil {
		return errors.New("tree has not been initialized")
	}
	rootFullpath := t.node.FullPath()
	if rootFullpath != "/" {
		rootFullpath += "/"
	}

	if !strings.HasPrefix(fi.FullPath(), rootFullpath) {
		return errors.New(fi.FullPath() + " is not contained in " + rootFullpath)
	}

	relativePath := strings.TrimPrefix(fi.FullPath(), rootFullpath)
	elems := strings.Split(relativePath, "/")
	curr := t
	for _, elem := range elems[:len(elems)-1] {
		found := false
		for _, n := range curr.children {
			if n.node.Name() == elem {
				curr = n
				found = true
				break
			}
		}
		if !found {
			return errors.New("Could not add " + fi.FullPath() + " to tree. Could only proceed up to element " + curr.node.FullPath())
		}
	}
	for _, elem := range curr.children {
		if elem.node.Name() == fi.Name() {
			return errors.New(fi.FullPath() + " already exists in tree")
		}
	}
	newNode := Tree{node: fi, children: make([]*Tree, 0)}
	curr.children = append(curr.children, &newNode)
	return nil
}

func (t *Tree) Count() int {
	count := 0
	var countNodes func(t *Tree)
	countNodes = func(t *Tree) {
		if t != nil {
			count++
		}
		for _, child := range t.children {
			countNodes(child)
		}
	}
	countNodes(t)
	return count
}

type byName []*Tree

func (b byName) Len() int {
	return len(b)
}

func (b byName) Swap(i, j int) {
	b[i], b[j] = b[j], b[i]
}

func (b byName) Less(i, j int) bool {
	return b[i].node.Name() < b[j].node.Name()
}

func (t *Tree) Sort() {
	var so func(t *Tree)
	so = func(t *Tree) {
		sort.Sort(sort.Reverse(byName(t.children)))
		for _, child := range t.children {
			so(child)
		}
	}
	so(t)
}

func (t *Tree) Create(root osabstraction.FileInfo, osw osabstraction.OSWrapper) error {
	err := t.Init(root)
	if err != nil {
		return err
	}

	rootPath := root.FullPath()

	var files []osabstraction.FileInfo

	regularFiles, err := osw.GetFiles(rootPath, true)
	if err != nil {
		return err
	}

	directories, err := osw.GetDirectories(rootPath)
	if err != nil {
		return err
	}

	files = append(regularFiles, files...)
	files = append(directories, files...)

	// NOTE: The following code for inserting files into the tree is
	// totally inefficient...
	// It works as follows: We iterate through the files in the
	// []FileInfo slice and try to insert these files into the tree
	// as long as we succeed. We can only insert a file with the
	// InsertSuccessor method, if the parrent directory of the file
	// has already been inserted to the tree.
	for {
		changed := false
		nextFiles := []osabstraction.FileInfo{}
		for _, fi := range files {
			err := t.InsertSuccessor(fi)
			if err == nil {
				changed = true
			} else {
				nextFiles = append(nextFiles, fi)
			}
		}
		if !changed {
			return errors.New("could not insert all the files into the tree")
		}
		if len(nextFiles) == 0 {
			break
		}
		files = nextFiles
	}
	t.Sort()
	return nil
}

func (t *Tree) String() string {
	rootLine := t.node.FullPath() + "\n"
	// output := ""
	// var slice []*Tree
	// for _, el := range t.children {
	// 	slice = append(slice, el)
	// }
	// levels := []int{1}
	// for len(slice) > 0 {
	// 	curr := slice[0]
	// 	slice = slice[1:]
	// 	currLevel := levels[0]
	// 	levels = levels[1:]
	// 	for i := 0; i < currLevel; i++ {
	// 		output += "────"
	// 	}
	// 	output += curr.node.Name() + "\n"
	// 	for _, child := range curr.children {
	// 		slice = append([]*Tree{child}, slice...)
	// 		levels = append([]int{currLevel + 1}, levels...)
	//
	// 	}
	// }
	// return rootLine + output
	return rootLine
}
