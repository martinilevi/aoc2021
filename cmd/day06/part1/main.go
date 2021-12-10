package main

import (
	"bufio"
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"
)

const targetDays = 80

func isError(err error) bool {
	if err != nil {
		_, filename, line, _ := runtime.Caller(1)
		fmt.Printf("%s:%d %s\n", filename, line, err.Error())
	}
	return (err != nil)
}

func growFishes(fishBank []int64) []int64 {
	born := 0
	for idx, v := range fishBank {
		if v != 0 {
			fishBank[idx] = fishBank[idx] - 1
		} else {
			fishBank[idx] = 6
			born = born + 1
		}
	}
	for x := 0; x < born; x++ {
		fishBank = append(fishBank, 8)
	}
	return fishBank
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

	//read header
	ln := 0
	fishBank := []int64{}
	for scanner.Scan() {
		ln = ln + 1
		line := scanner.Text()

		fishVec := strings.Split(line, ",")
		if len(fishVec) == 0 {
			isError(fmt.Errorf("no fish!", ln))
			return
		}

		for _, v := range fishVec {
			newFish, err := strconv.ParseInt(v, 10, 64)
			if isError(err) {
				return
			}
			fishBank = append(fishBank, newFish)
		}
	}

	days := 0
	for ; days < targetDays; days++ {
		//fmt.Printf("day %d: %v\n", days, fishBank)
		fmt.Printf("day %d: %d fish\n", days, len(fishBank))
		fishBank = growFishes(fishBank)
	}
	//fmt.Printf("day %d: %v\n", days, fishBank)
	fmt.Printf("day %d: %d fish\n", days, len(fishBank))

}