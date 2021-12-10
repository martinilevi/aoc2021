package main

import (
	"bufio"
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"
)

const targetDays = 256

func isError(err error) bool {
	if err != nil {
		_, filename, line, _ := runtime.Caller(1)
		fmt.Printf("%s:%d %s\n", filename, line, err.Error())
	}
	return (err != nil)
}

func fishCount(age [9]int64) (c int64) {
	for _, v := range age {
		c = c + v
	}
	return
}

func growFish(fishBank *[9]int64) {

	cpy := [9]int64{}
	for x := 0; x < 9; x++ {
		cpy[x] = fishBank[x]
	}

	//one day less for everyone
	for x:= 8; x > 0; x-- {
		fishBank[x-1] = cpy[x]
	}

	//old zeroes now have 6 days to bring a new fish
	fishBank[6] = fishBank[6] + cpy[0]

	//old zeroes give birth to new fishes
	fishBank[8] = cpy[0]

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
	age := [9]int64{}
	for scanner.Scan() {
		ln = ln + 1
		line := scanner.Text()

		fishVec := strings.Split(line, ",")
		if len(fishVec) == 0 {
			isError(fmt.Errorf("no fish!", ln))
			return
		}

		for _, v := range fishVec {
			newFish, err := strconv.ParseUint(v, 10, 8)
			if isError(err) {
				return
			}
			if newFish > 8 {
				isError(fmt.Errorf("invalid initial age %d", newFish))
				return
			}
			age[newFish] = age[newFish] + 1
		}
	}

	days := 0
	for ; days < targetDays; days++ {
		//fmt.Printf("day %d: %v\n", days, fishBank)
		fmt.Printf("day %d: %d fish\n", days, fishCount(age))
		//fmt.Println(age)
		growFish(&age)
	}
	//fmt.Printf("day %d: %v\n", days, fishBank)
	fmt.Printf("day %d: %d fish\n", days, fishCount(age))
}