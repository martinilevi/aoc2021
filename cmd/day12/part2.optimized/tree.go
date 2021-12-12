package main

import (
    "fmt"
    "io"
    "os"
)

func print(w io.Writer, prefix string, node *TreeNode) {
    if node == nil {
        return
    }

    fmt.Fprintf(w, "%s%s\n", prefix, node.Data)
    for _, v := range node.Sons {
        print(w, prefix+" ", v)
    }
}
 
func TreeTest() {
    tree := &Tree{}
    tree.Insert("a,b,c,d")
    tree.Insert("a,c,d")
    tree.Insert("a,d")
    print(os.Stdout, "", tree.Root)
    fmt.Println("Search a", tree.Search("a"))
    fmt.Println("Search b", tree.Search("b"))
    fmt.Println("Search a,b,c", tree.Search("a,b,c"))
    fmt.Println("Search a,b,c,e", tree.Search("a,b,c,e"))
    fmt.Println("Search a,b,c,d", tree.Search("a,b,c,d"))
    fmt.Println("Search a,d", tree.Search("a,d"))
    fmt.Println("Search empty", tree.Search(""))
    /*
root
 a
  b
   c
    d
  c
   d
  d
search a true
search b false
search a,b,c true
search a,b,c,e false
search a,b,c,d true
search a,d true
search empty false

    */
}
