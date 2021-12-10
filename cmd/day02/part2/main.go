package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func isError(err error) bool {
	if err != nil {
		fmt.Println(err.Error())
	}
	return (err != nil)
}

const (
	up = iota
	down
	forward
)


func GetPosition(strPositionChange string) (where int8, howmuch int64, err error) {
	res := strings.Split(strPositionChange, " ")
	if len(res) != 2 {
		return -1, -1, fmt.Errorf("invalid input line: '%s'", strPositionChange)
	}
	whereStr := res[0]
	howmuchStr := res[1]

	switch whereStr {
	case "up":
		where = up
	case "down":
		where = down
	case "forward":
		where = forward
	}

	howmuch, err = strconv.ParseInt(howmuchStr, 0, 64)

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

	x := int64(0)
	y := int64(0)
	aim := int64(0)
	for scanner.Scan() {
		strPositionChange := scanner.Text()
		where, howmuch, err := GetPosition(strPositionChange)

		if isError(err) {
			return
		}

		switch where {
		case up:
			aim = aim - howmuch
		case down:
			aim = aim + howmuch
		case forward:
			x = x + howmuch
			y = y - aim * howmuch
		}
	}

	depth := -y

	fmt.Printf("forward: %d, depth: %d, product: %d\n", x, depth, x*depth)
}

