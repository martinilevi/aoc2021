package main

import (
	"bufio"
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"
)

//const screenSize = 10
const screenSize = 1000

type ScreenType [screenSize][screenSize]int64
type Point struct {
	X int64
	Y int64
}

func DrawPoints(screen *ScreenType, x1, y1, x2, y2 int64) {
	if x1 == x2 && y1 == y2 {
		fmt.Println("no segment")
		return
	}

	switch {
	case x1 == x2:
		if y2 >= y1 {
			for y:= y1; y <= y2; y++ {
				screen[x1][y] = screen[x1][y] + 1
			}
		} else {
			for y:= y1; y >= y2; y-- {
				screen[x1][y] = screen[x1][y] + 1
			}
		}
	case y1 == y2:
		if x2 >= x1 {
			for x:= x1; x <= x2; x++ {
				screen[x][y1] = screen[x][y1] + 1
			}
		} else {
			for x:= x1; x >= x2; x-- {
				screen[x][y1] = screen[x][y1] + 1
			}
		}
	default:
		incX := int64(-1)
		if x2 >= x1 {
			incX = 1
		}

		incY := int64(-1)
		if y2 >= y1 {
			incY = 1
		}
		//fmt.Println("from", x1,y1)
		//fmt.Println("to", x2,y2)
		for x, y := x1, y1; x != x2+incX || y != y2+incY; x, y=x+incX, y+incY {
			//fmt.Println(x,y)
			screen[x][y] = screen[x][y] + 1
	    }
	}


}

func Print(screen *ScreenType) {
	for y := 0; y < screenSize; y++ {
		for x := 0; x < screenSize; x++ {
			if screen[x][y] == 0 {
				fmt.Printf(".")
			} else {
				fmt.Printf("%d", screen[x][y])
			}
		}
		fmt.Printf("\n")
	}
}

func OverlPoints(screen *ScreenType) (cnt int64) {
	for y := 0; y < screenSize; y++ {
		for x := 0; x < screenSize; x++ {
			if screen[x][y] >= 2 {
				cnt = cnt + 1
			}
		}
	}
	return
}

func isError(err error) bool {
	if err != nil {
		_, filename, line, _ := runtime.Caller(1)
		fmt.Printf("%s:%d %s\n", filename, line, err.Error())
	}
	return (err != nil)
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
	sc := ScreenType{}
	for scanner.Scan() {
		ln = ln + 1
		line := scanner.Text()

		points := strings.Split(line, " -> ")
		if len(points) != 2 {
			isError(fmt.Errorf("point missing at line %d", ln))
			return
		}
		point1 := strings.Split(points[0],",")
		if len(point1) != 2 {
			isError(fmt.Errorf("point 1 missing coords at line %d", ln))
			return
		}
		point2 := strings.Split(points[1],",")
		if len(point2) != 2 {
			isError(fmt.Errorf("point 2 missing coords at line %d", ln))
			return
		}
		//fmt.Println(point1, point2)
		x1, err := strconv.ParseInt(point1[0], 0, 64)
		if isError(err) {
			return
		}
		y1, err := strconv.ParseInt(point1[1], 0, 64)
		if isError(err) {
			return
		}
		x2, err := strconv.ParseInt(point2[0], 0, 64)
		if isError(err) {
			return
		}
		y2, err := strconv.ParseInt(point2[1], 0, 64)
		if isError(err) {
			return
		}

		DrawPoints(&sc, x1, y1, x2, y2)
	}
	//Print(&sc)
	fmt.Printf("Has %d overlap points\n", OverlPoints(&sc))
}