package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

func isError(err error) bool {
	if err != nil {
		fmt.Println(err.Error())
	}
	return (err != nil)
}

func main() {

	file, err := os.OpenFile("day1.input", os.O_RDWR, 0644)

	if isError(err) {
		return
	}

	defer func(){
		err = file.Close()
		_ = isError(err)
	}()

	scanner := bufio.NewScanner(file)
	previous := int64(0)
	increases := 0
	for scanner.Scan() {
		strDepth := scanner.Text()
		depth, err := strconv.ParseInt(strDepth, 0, 64)
		if isError(err) {
			return
		}
		if previous !=0 {
			if depth - previous > 0 {
				increases = increases + 1
			}
		}
		previous = depth
	}

	fmt.Printf("%d increases\n", increases)
}

