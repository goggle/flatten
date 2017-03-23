package filesystem

import (
	"fmt"
	"testing"
)

func TestTree(t *testing.T) {
	fs := Filesystem{}
	fs.Init()
	fs.MkDir("/home/goggle/test")
	fs.MkDir("/home/goggle/test/rust/tree")
	fs.MkDir("/home/goggle/test/documents/a001/a001")
	fs.CreateFile("/home/goggle/test/artisan.txt")
	fs.CreateFile("/home/goggle/test/hello.c")
	fs.CreateFile("/home/goggle/test/hello.h")
	fs.CreateFile("/home/goggle/test/railjet_data.txt")
	fs.CreateFile("/home/goggle/test/documents/cv.pdf")
	fs.CreateFile("/home/goggle/test/documents/exams.pdf")
	fs.CreateFile("/home/goggle/test/documents/form.docx")
	fs.CreateFile("/home/goggle/test/documents/a001/a001/mri_010023.csv")
	fs.CreateFile("/home/goggle/test/documents/a001/a001/mri_010024.csv")
	fs.CreateFile("/home/goggle/test/rust/main.rs")
	fs.CreateFile("/home/goggle/test/rust/ask.rs")
	fs.CreateFile("/home/goggle/test/rust/tree/tree.rs")
	fs.CreateFile("/home/goggle/test/rust/tree/doc.txt")

	tree := Tree{}
	err := tree.Create(fs["/home/goggle/test"], fs)
	if err != nil {
		t.Errorf("TestTree: No error expected, got %v", err)
	}

	expectedTreeString := `/home/goggle/test
├── artisan.txt
├── documents
│   ├── a001
│   │   └── a001
│   │       ├── mri_010023.csv
│   │       └── mri_010024.csv
│   ├── cv.pdf
│   ├── exams.pdf
│   └── form.docx
├── hello.c
├── hello.h
├── railjet_data.txt
└── rust
    ├── ask.rs
    ├── main.rs
    └── tree
        ├── doc.txt
        └── tree.rs
`
	if expectedTreeString != fmt.Sprintf("%v", tree) {
		t.Errorf("TreeTest: Expected %v, got %v", expectedTreeString, tree)
	}
}
