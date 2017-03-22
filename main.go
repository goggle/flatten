package main

import (
	"fmt"

	"github.com/goggle/flatten/osabstraction"
)

func main() {
	// path := "/home/alex/test/"
	// results := flatten.GetNestedFiles(path)
	// for _, elem := range results {
	// 	fmt.Println(elem.Name())
	// }
	// i := 22
	// fmt.Printf("%05d\n", i)

	osw := osabstraction.RealOS{}
	err := osw.Copy("/home/alex/test/hello/readme.txt", "/home/alex/test/hello/dont.txt")
	if err != nil {
		fmt.Println(err)
	}

}
