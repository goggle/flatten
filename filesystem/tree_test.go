package filesystem

import (
	"fmt"
	"testing"

	"github.com/goggle/flatten/osabstraction"
)

func TestTree(t *testing.T) {
	fs := Filesystem{}
	fs.Init()
	fs.MkDir("/tmp/a/b/c")
	fs.CreateFile("/tmp/a/hello.c")
	fs.CreateFile("/tmp/a/song.mp3")
	fs.CreateFile("/tmp/a/watch.mp4")
	fs.CreateFile("/tmp/a/arrifix.txt")
	fs.CreateFile("/tmp/a/b/gogol.go")

	// tr := Tree{}
	// tr.Init(fs["/tmp"])
	// tr.InsertSuccessor(fs["/tmp/a"])
	// tr.InsertSuccessor(fs["/tmp/a/b"])
	// tr.InsertSuccessor(fs["/tmp/a/hello.c"])
	// tr.InsertSuccessor(fs["/tmp/a/song.mp3"])
	// tr.InsertSuccessor(fs["/tmp/a/watch.mp4"])
	// tr.InsertSuccessor(fs["/tmp/a/b/c"])
	// tr.InsertSuccessor(fs["/tmp/a/b/gogol.go"])
	// tr.InsertSuccessor(fs["/tmp/a/arrifix.txt"])
	// tr.Sort()
	// fmt.Printf(tr.String())
	// tr.Create(fs["/tmp"], fs)
	// fmt.Printf(tr.String())

	tt := Tree{}
	osw := osabstraction.RealOS{}
	root := osabstraction.File("/home/alex/baz")
	err := tt.Create(root, osw)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(tt.Count())
	// fmt.Printf(tt.String())
	// fmt.Println(tt.children)
}
