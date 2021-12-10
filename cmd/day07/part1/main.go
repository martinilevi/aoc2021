package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"runtime"
	"strconv"
	"strings"
)

func isError(err error) bool {
	if err != nil {
		_, filename, line, _ := runtime.Caller(1)
		fmt.Printf("%s:%d %s\n", filename, line, err.Error())
	}
	return (err != nil)
}

func Abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func findMinAndMax(a []int) (min int, max int) {
	min = a[0]
	max = a[0]
	for _, value := range a {
		if value < min {
			min = value
		}
		if value > max {
			max = value
		}
	}
	return min, max
}

func Align(fuelPos []int) (cost int, position int) {
	min, max := findMinAndMax(fuelPos)

	fmt.Printf("Searching minimum position between (%d, %d) on a set of %d crabs...\n",
		min, max, len(fuelPos))
	position = -1
	fuelUsed := math.MaxInt64
	for x:=min; x<=max; x++ {
		tmp := 0
		for _, v := range fuelPos {
			tmp = tmp + Abs(x-v)
		}
		if tmp < fuelUsed {
			position = x
			fuelUsed = tmp
		}
	}
	cost = fuelUsed
	return
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


	fuelPositions := []int{}

	scanner.Scan()
	line := scanner.Text()

	fuelVecStr := strings.Split(line, ",")

	if len(fuelVecStr) == 0 {
		isError(fmt.Errorf("no fuel!"))
		return
	}

	for _, v := range fuelVecStr {
		fuelPosition, err := strconv.ParseInt(v, 10, 64)
		if isError(err) {
			return
		}
		fuelPositions = append(fuelPositions, int(fuelPosition))
	}

	cost, position := Align(fuelPositions)

	if position == -1 {
		isError(fmt.Errorf("No minimum found"))
		return
	}

	fmt.Println("Minimum cost is", cost)
	fmt.Println("when aligned At position", position)
}