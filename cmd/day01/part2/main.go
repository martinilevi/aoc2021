package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

const windowSize = 3

func isError(err error) bool {
	if err != nil {
		fmt.Println(err.Error())
	}
	return (err != nil)
}

// hasIncreased: tells if a window has increased, resets oldest window.
func hasIncreased(it []int64, add []int64) bool {
	first := 0
	second := 0

	fmt.Println(it)
	for x:= 0; x < windowSize; x++ {
		if it[x] == windowSize + 1 {
			first = x
			fmt.Printf("first index is %d\n", x)
		} else if it[x] == windowSize {
			second = x
			fmt.Printf("second index is %d\n", x)
		}
	}

	if first == second {
		//protect against incomplete/undefined dataset
		return false
	}

	defer func(){
		fmt.Printf("will reset idx %d\n", first)
		it[first] = 0
		add[first] = 0
	}()

	fmt.Printf("add1st == %d, add2nd == %d \n", add[first], add[second])

	if add[second] > add[first] {
		fmt.Println("increase!\n")
		return true
	}
	fmt.Println("not an increase\n")
	return false
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
	increases := 0

	//array de int64 ... TODO finish for 3 sliding windows
	add := []int64{}
	it := []int64{}

	for x := 0; x < windowSize; x++ {
		add = append(add, int64(0))
		it = append(it, int64(0))
	}

	line := 0
	for scanner.Scan() {
		strDepth := scanner.Text()
		depth, err := strconv.ParseInt(strDepth, 0, 64)

		if isError(err) {
			return
		}

		limit := windowSize - 1
		if line < windowSize {
			limit = line
		}

		resetIdx := 0
		resetValue := int64(0)

		for x := 0; x <= limit; x++ {
			it[x] = it[x] + 1
			if it[x] == windowSize + 1 {
				resetIdx = x
				resetValue = depth
				continue
			}
			add[x] = add[x] + depth
		}

		if line < windowSize {
			line = line + 1
			continue
		}
		line = line + 1

		if hasIncreased(it, add) {
			increases = increases + 1
		}

		it[resetIdx] = 1
		add[resetIdx] = resetValue

	}

	fmt.Printf("%d increases\n", increases)
}

