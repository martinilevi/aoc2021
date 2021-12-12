package main

import (
	"bufio"
	"fmt"
	"os"
	"runtime"
	"strings"
	"unicode"
)

func isError(err error) bool {
	if err != nil {
		_, filename, line, _ := runtime.Caller(1)
		fmt.Printf("%s:%d %s\n", filename, line, err.Error())
	}
	return (err != nil)
}

func IsLower(s string) bool {
	for _, r := range s {
		if !unicode.IsLower(r) && unicode.IsLetter(r) {
			return false
		}
	}
	return true
}

type GraphNode struct {	
	Connections []string
	Visited bool
	Small bool
	Start bool
	End bool
}

func AddNode(grafo map[string]*GraphNode, name string) {
	start := false
	end := false
	small := false
	if name == "start" {
		start = true
	} else if name == "end" {
		end = true
	} else if IsLower(name) {
		small = true
	}
	grafo[name] = &GraphNode{nil,false, small, start, end}
}

func AddLink(grafo map[string]*GraphNode, n1, n2 string) {
	grafo[n1].Connections = append(grafo[n1].Connections, n2)
	grafo[n2].Connections = append(grafo[n2].Connections, n1)
}

func Copy(grafo map[string]*GraphNode) (copy map[string]*GraphNode) {
	copy = map[string]*GraphNode{}
	for k, v := range grafo {
		copy[k]=&GraphNode{
			Visited: false,
			Small: v.Small,
			Start: v.Start,
			End: v.End,
		}
		for _, name := range v.Connections {
			copy[k].Connections = append(copy[k].Connections, name)
		}
	}
	return
}

//TODO: This function could be optmized if possible paths were stored in a tree, search would be log(n) instead of O(n)
func IsSubPathPresent(paths map[string]bool, path string) bool {
	for p, _ := range paths {
		if strings.HasPrefix(p, path) {
			return true
		}
	}
	return false
}

func IsAnySmallCaveRepeated(path string) bool {
	pathParts := strings.Split(path,",")
	for _, v := range pathParts {
		if v == "start" {
			continue
		}
		if v == "end" {
			continue
		}
		if !IsLower(v) {
			continue
		}
		target := v
		cnt := 0
		for _, w := range pathParts {
			if w == target {
				cnt = cnt + 1
			}
		}
		if cnt > 1 {
			return true
		}
	}
	return false
}

func IsNodeOnCurrentPath(name string, path string) bool {
	pathParts := strings.Split(path,",")
	for _, v := range pathParts {
		if name == v {
			return true
		}
	}
	return false
}

func HasBigBrothers(grafo map[string]*GraphNode, name string) (hasBigBrothers bool) {
	for _, vv := range grafo[name].Connections {
		if grafo[vv].Start {
			continue
		}
		if !grafo[vv].Small {
			hasBigBrothers = true
			break
		}
	}
	return
}

func HasSmallBrothers(grafo map[string]*GraphNode, name string) (hasSmallBrothers bool) {
	for _, vv := range grafo[name].Connections {
		if grafo[vv].Start {
			continue
		}
		if grafo[vv].Small {
			hasSmallBrothers = true
			break
		}
	}
	return
}

func PathLevel(path string) int {
	pathParts := strings.Split(path,",")
	return len(pathParts)
}

func Tabs(t int) (res string) {
	for x:=0; x<t; x++ {
		res = res + "\t"
	}
	return
}

func VisitNode(grafo map[string]*GraphNode, name, path string, paths map[string]bool) {
	//fmt.Println(path)
	node := grafo[name]
	node.Visited = true
	for _, v := range node.Connections {
		//start
		if v == "start" {
			//won't go back to start
			continue
		}

		//end
		if v == "end" {
			if ! paths[path+",end"] {
				//fmt.Println(">>>",path+",end")
				paths[path+",end"]=true
			}
			if len(paths) % 1000 == 0 {
				fmt.Println(len(paths))
			}
			continue
		}

		//small cave
		if grafo[v].Small {
			if IsNodeOnCurrentPath(v, path) && IsAnySmallCaveRepeated(path) {
				//small caves may be repeated once
				continue
			}

			/*
			if HasBigBrothers(grafo, v) {
				VisitNode(grafo, v, path+","+v, paths)
			} else {
				if HasSmallBrothers(grafo, v) && !IsAnySmallCaveRepeated(path) {
					VisitNode(grafo, v, path+","+v, paths)
				}
			}*/
			VisitNode(grafo, v, path+","+v, paths)
			continue
		}

		//big cave
		if grafo[v].Visited {
			hasUnvisitedSiblings := false
			for _, vv := range grafo[v].Connections {
				if grafo[vv].Start {
					continue
				}
				if !grafo[vv].Visited {
					hasUnvisitedSiblings = true
					break
				}
			}
			if hasUnvisitedSiblings || !IsSubPathPresent(paths, path+","+v) {
				VisitNode(grafo, v, path+","+v, paths)
			}
		} else {
			//big one not visited, visit
			VisitNode(grafo, v, path+","+v, paths)
		}
	}
}

func Walk(grafo map[string]*GraphNode) int {
	fmt.Println(grafo["start"].Connections)
	paths := map[string]bool{}
	VisitNode(grafo, "start", "start", paths)
	return len(paths)
}

func main() {
	file, err := os.OpenFile("input", os.O_RDWR, 0644)

	if isError(err) {
		return
	}

	defer func(){
		err = file.Close()
		_ = isError(err)
	}()

	scanner := bufio.NewScanner(file)

	grafo := map[string]*GraphNode{}
	ln := 0
	for scanner.Scan() {
		ln = ln + 1
		line := scanner.Text()
		
		arrow := strings.Split(line, "-")
		if len(arrow) != 2 {
			isError(fmt.Errorf("Invalid arrow at line %s", ln))
			return
		}
		from := arrow[0]
		to := arrow[1]
		
		if grafo[from] == nil {
			AddNode(grafo, from)
		}

		if grafo[to] == nil {
			AddNode(grafo, to)
		}
		AddLink(grafo, from, to)
	}

	fmt.Println("Graph with", len(grafo), "nodes loaded.")
	fmt.Println(Walk(grafo),"path(s) found")
}