package main

import (
	"bufio"
	"fmt"
	"os"
	"runtime"
	"strconv"
)

func isError(err error) bool {
	if err != nil {
		_, filename, line, _ := runtime.Caller(1)
		fmt.Printf("%s:%d %s\n", filename, line, err.Error())
	}
	return (err != nil)
}

func HasUp(lineSize, colSize, linePos, colPos int) bool {
	if linePos == 0 {
		return false
	}

	return true
}

func HasDown(lineSize, colSize, linePos, colPos int) bool {
	if linePos == lineSize - 1 {
		return false
	}

	return true
}

func HasLeft(lineSize, colSize, linePos, colPos int) bool {
	if colPos == 0 {
		return false
	}

	return true
}

func HasRight(lineSize, colSize, linePos, colPos int) bool {

	if colPos == colSize - 1 {
		return false
	}

	return true
}

func HasUpLeft(lineSize, colSize, linePos, colPos int) bool {
	if HasUp(lineSize, colSize, linePos, colPos) && HasLeft(lineSize, colSize, linePos, colPos) {
		return true
	}

	return false
}

func HasUpRight(lineSize, colSize, linePos, colPos int) bool {
	if HasUp(lineSize, colSize, linePos, colPos) && HasRight(lineSize, colSize, linePos, colPos) {
		return true
	}

	return false
}

func HasDownLeft(lineSize, colSize, linePos, colPos int) bool {
	if HasDown(lineSize, colSize, linePos, colPos) && HasLeft(lineSize, colSize, linePos, colPos) {
		return true
	}

	return false
}

func HasDownRight(lineSize, colSize, linePos, colPos int) bool {
	if HasDown(lineSize, colSize, linePos, colPos) && HasRight(lineSize, colSize, linePos, colPos) {
		return true
	}

	return false
}

func incSurroundings(mat [][]Point, lineSize, colSize, linePos, colPos int) (flashes int) {

	if HasUpLeft(lineSize, colSize, linePos, colPos) {
		mat[linePos-1][colPos-1].Value++
	}

	if HasUp(lineSize, colSize, linePos, colPos) {
		mat[linePos-1][colPos].Value++
	}

	if HasUpRight(lineSize, colSize, linePos, colPos) {
		mat[linePos-1][colPos+1].Value++
	}
	if HasLeft(lineSize, colSize, linePos, colPos) {
		mat[linePos][colPos-1].Value++
	}

	if HasRight(lineSize, colSize, linePos, colPos) {
		mat[linePos][colPos+1].Value++
	}

	if HasDownLeft(lineSize, colSize, linePos, colPos) {
		mat[linePos+1][colPos-1].Value++
	}

	if HasDown(lineSize, colSize, linePos, colPos) {
		mat[linePos+1][colPos].Value++
	}

	if HasDownRight(lineSize, colSize, linePos, colPos) {
		mat[linePos+1][colPos+1].Value++
	}

	return
}

type Point struct {
	Value int64
	Flashed bool
}

func FlashStep(nrg [][]int64) (res int) {
	lineSize := len(nrg)
	colSize := len(nrg[0])

	//copy to memory matrix
	tmp := [][]Point{}
	for linePos := range nrg {
		tmp = append(tmp, []Point{})
		for _, x := range nrg[linePos] {
			tmp[linePos] = append(tmp[linePos], Point{x,false})
		}
	}

	//fmt.Println("inc all")

	for linePos := range tmp {
		for colPos, _ := range tmp[linePos] {
			tmp[linePos][colPos].Value++
		}
	}

	//PrintPointMatrix("\t", tmp)

	//flash
	seenFlash := true
	for seenFlash {
		seenFlash = false
		for linePos := range tmp {
			for colPos, x := range tmp[linePos] {
				//fmt.Println("position", linePos, colPos)
				if x.Value > 9 && !tmp[linePos][colPos].Flashed {
					seenFlash = true
					//fmt.Println("flash!", linePos, colPos)
					tmp[linePos][colPos].Flashed = true
					//PrintPointMatrix("\t", tmp)
					res = res + 1
					incSurroundings(tmp, lineSize, colSize, linePos, colPos)
				}
				//PrintPointMatrix("\t", tmp)
				//time.Sleep(time.Second)
			}
		}
	}


	//zeroes
	for linePos := range tmp {
		for colPos, x := range tmp[linePos] {
			if x.Value > 9 {
				tmp[linePos][colPos].Value = 0
			}
		}
	}

	//copy to result matrix
	for linePos := range nrg {
		for colPos, _ := range nrg[linePos] {
			nrg[linePos][colPos] = tmp[linePos][colPos].Value
		}
	}

	return
}

func PrintPointMatrix(prefix string, nrg [][]Point) {
	const HEADER = "\033[95m"
	const OKBLUE = "\033[94m"
	const OKCYAN = "\033[96m"
	const OKGREEN = "\033[92m"
	const WARNING = "\033[93m"
	const FAIL = "\033[91m"
	const ENDC = "\033[0m"
	const BOLD = "\033[1m"
	const UNDERLINE = "\033[4m"

	for linePos := range nrg {
		fmt.Print(prefix)
		for _, x := range nrg[linePos] {
			if x.Value == 0 {
				pre := OKBLUE
				pos := ENDC
				if x.Flashed {
					pre = UNDERLINE+OKBLUE
				}
				fmt.Printf(" %s%02d%s", pre, x.Value, pos)
				continue
			} else if x.Value == 9 {
				pre := OKGREEN
				pos := ENDC
				if x.Flashed {
					pre = UNDERLINE+OKGREEN
				}
				fmt.Printf(" %s%02d%s", pre, x.Value, pos)
				continue
			}
			pre := ""
			pos := ""
			if x.Flashed {
				pre = UNDERLINE
				pos = ENDC
			}
			fmt.Printf(" %s%02d%s", pre, x.Value, pos)
		}
		fmt.Println(prefix)
	}
	fmt.Println()
}

func PrintMatrix(prefix string, nrg [][]int64) {
	const HEADER = "\033[95m"
	const OKBLUE = "\033[94m"
	const OKCYAN = "\033[96m"
	const OKGREEN = "\033[92m"
	const WARNING = "\033[93m"
	const FAIL = "\033[91m"
	const ENDC = "\033[0m"
	const BOLD = "\033[1m"
	const UNDERLINE = "\033[4m"

	for linePos := range nrg {
		for _, x := range nrg[linePos] {
			if x == 0 {
				fmt.Print(prefix, OKBLUE, x ,ENDC)
				continue
			} else if x == 9 {
				fmt.Print(prefix, OKGREEN, x ,ENDC)
				continue
			}
			fmt.Print(prefix, x)
		}
		fmt.Println()
	}
	fmt.Println()
}

func IsZeroMatrix(nrg [][]int64) bool {
	for linePos := range nrg {
		for _, x := range nrg[linePos] {
			if x != 0 {
				return false
			}
		}
	}
	return true
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

	nrgMatrix := [][]int64{}

	ln := 0
	for scanner.Scan() {
		ln = ln + 1
		line := scanner.Text()
		for _, n := range line {
			x, err := strconv.ParseInt(string(n), 10, 64)
			if isError(err) {
				return
			}
			if len(nrgMatrix) < ln {
				nrgMatrix = append(nrgMatrix, []int64{})
			}
			nrgMatrix[ln-1] = append(nrgMatrix[ln-1], x)
		}
	}

	step := 0
	for step = 0; !IsZeroMatrix(nrgMatrix); step++ {
		FlashStep(nrgMatrix)
	}
	fmt.Println("Zero matrix found at step", step)
}