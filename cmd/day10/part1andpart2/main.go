package main

import (
	"bufio"
	"fmt"
	"os"
	"runtime"
	"sort"
)

func isError(err error) bool {
	if err != nil {
		_, filename, line, _ := runtime.Caller(1)
		fmt.Printf("%s:%d %s\n", filename, line, err.Error())
	}
	return (err != nil)
}

func closingChar(openChar int32) int32 {
	switch openChar {
	case '(':
		return ')'
	case '[':
		return ']'
	case '<':
		return '>'
	case '{':
		return '}'
	}
	return -1
}

func openChar(openChar int32) int32 {
	switch openChar {
	case ')':
		return '('
	case ']':
		return '['
	case '>':
		return '<'
	case '}':
		return '{'
	}
	return -1
}

func IsCorrupted(line string) (bool, int32) {
	pila := ""
	for _, c := range line {
		//fmt.Println("stack", pila)
		switch c {
		case '(':
			pila = pila + "("
		case '[':
			pila = pila + "["
		case '<':
			pila = pila + "<"
		case '{':
			pila = pila + "{"
		case ')':
			fallthrough
		case ']':
			fallthrough
		case '>':
			fallthrough
		case '}':
			if len(pila) == 0 {
				fmt.Printf("closing %s on empty stack\n", string(c))
				return true, c
			}
			if pila[len(pila)-1:] != string(openChar(c)) {
				fmt.Printf("expecting %s but got %s \n", string(closingChar(int32(pila[len(pila)-1:][0]))), string(c))
				return true, c
			}
			pila = pila[0:len(pila)-1]
		}
	}
	return false, 0
}

func CalcIllegalValue(c int32) int {
	const parenthesis = 3
	const bracket = 57
	const curlyBracket = 1197
	const greaterThan = 25137
	switch c {
	case ')':
		return parenthesis
	case ']':
		return bracket
	case '}':
		return curlyBracket
	case '>':
		return greaterThan
	}
	return 0
}

func CalcIncompleteValue(missing string) int {
	const parenthesis = 1
	const bracket = 2
	const curlyBracket = 3
	const greaterThan = 4
	res := 0
	for _, c := range missing {
		switch c {
		case ')':
			res = 5*res + parenthesis
		case ']':
			res = 5*res + bracket
		case '}':
			res = 5*res + curlyBracket
		case '>':
			res = 5*res + greaterThan
		}
	}
	return res
}

func IsCorruptedOrComplete(line string) (bool, int32, string) {
	pila := ""
	for _, c := range line {
		//fmt.Println("stack", pila)
		switch c {
		case '(':
			pila = pila + "("
		case '[':
			pila = pila + "["
		case '<':
			pila = pila + "<"
		case '{':
			pila = pila + "{"
		case ')':
			fallthrough
		case ']':
			fallthrough
		case '>':
			fallthrough
		case '}':
			if len(pila) == 0 {
				fmt.Printf("closing %s on empty stack\n", string(c))
				return true, c, ""
			}
			if pila[len(pila)-1:] != string(openChar(c)) {
				fmt.Printf("expecting %s but got %s \n", string(closingChar(int32(pila[len(pila)-1:][0]))), string(c))
				return true, c, ""
			}
			pila = pila[0:len(pila)-1]
		}
	}

	missing := ""
	if len(pila) != 0 {
		for x := len(pila) - 1; x >= 0; x-- {
			cc := closingChar(int32(pila[x]))
			missing = missing + string(cc)
		}
	}

	return false, 0, missing
}

func main1() {
	file, err := os.OpenFile("input", os.O_RDWR, 0644)
	if isError(err) {
		return
	}

	defer func(){
		err = file.Close()
		_ = isError(err)
	}()

	scanner := bufio.NewScanner(file)

	ln := 0
	corruptScore := 0
	for scanner.Scan() {
		ln = ln + 1
		line := scanner.Text()
		//fmt.Println("evaluating", line)
		isLineCorrupted, illChar := IsCorrupted(line)
		if isLineCorrupted {
			fmt.Println("line", ln, "is corrupted")
			corruptScore = corruptScore + CalcIllegalValue(illChar)
			continue
		}
		fmt.Println("line", ln, "is OK")
	}

	fmt.Println("illegal chars score is", corruptScore)
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

	ln := 0
	corruptScore := 0
	incompleteScoreList := []int{}
	for scanner.Scan() {
		ln = ln + 1
		line := scanner.Text()
		//fmt.Println("evaluating", line)
		isLineCorrupted, illChar, missing := IsCorruptedOrComplete(line)
		if isLineCorrupted {
			fmt.Println("line", ln, "is corrupted")
			corruptScore = corruptScore + CalcIllegalValue(illChar)
			continue
		}
		if len(missing) != 0 {
			fmt.Println("line", ln, "is incomplete, completed with", missing)
			incompleteScoreList = append(incompleteScoreList, CalcIncompleteValue(missing))
			continue
		}
		fmt.Println("line", ln, "is OK")
	}

	sort.Ints(incompleteScoreList)
	fmt.Println(incompleteScoreList)
	fmt.Println("Middlescore winner is", incompleteScoreList[int(len(incompleteScoreList)/2)])

}