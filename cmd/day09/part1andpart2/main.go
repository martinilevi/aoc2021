package main

import (
	"bufio"
	"fmt"
	"os"
	"runtime"
	"sort"
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

func checkIfMin(mat [][]int64, lineSize, colSize, linePos, colPos int) bool {
	if HasUp(lineSize, colSize, linePos, colPos) {
		if mat[linePos][colPos] >= mat[linePos-1][colPos] {
			return false
		}
	}
	if HasDown(lineSize, colSize, linePos, colPos) {
		if mat[linePos][colPos] >= mat[linePos+1][colPos] {
			return false
		}
	}
	if HasLeft(lineSize, colSize, linePos, colPos) {
		if mat[linePos][colPos] >= mat[linePos][colPos-1] {
			return false
		}
	}
	if HasRight(lineSize, colSize, linePos, colPos) {
		if mat[linePos][colPos] >= mat[linePos][colPos+1] {
			return false
		}
	}
	return true
}

func findLocalMins(heights [][]int64) (res []int64) {
	lineSize := len(heights)
	colSize := len(heights[0])
	fmt.Printf("matrix %d x %d\n", lineSize, colSize)
	for linePos := range heights {
		for colPos, x := range heights[linePos] {
			//fmt.Printf("%d ", x)
			if checkIfMin(heights, lineSize, colSize, linePos, colPos) {
				res = append(res, x)
			}
		}
		//fmt.Printf("\n")
	}
	return
}

func findLocalMinsPositions(heights [][]int64) (res [][2]int) {
	lineSize := len(heights)
	colSize := len(heights[0])
	fmt.Printf("matrix %d x %d\n", lineSize, colSize)
	for linePos := range heights {
		for colPos, _ := range heights[linePos] {
			//fmt.Printf("%d ", x)
			if checkIfMin(heights, lineSize, colSize, linePos, colPos) {
				res = append(res,[2]int{linePos, colPos})
			}
		}
		//fmt.Printf("\n")
	}
	return
}

func markUpBasin(mat [][]int64, mark [][]PointData, lineSize, colSize, linePos, colPos int) {
	//fmt.Println("marking surroundings of ", linePos, colPos)
	if HasUp(lineSize, colSize, linePos, colPos) {
		if mat[linePos-1][colPos] != 9 {
			mark[linePos-1][colPos].PartOfBasin = true
		}
	}
}

func markDownBasin(mat [][]int64, mark [][]PointData, lineSize, colSize, linePos, colPos int) {
	//fmt.Println("marking surroundings of ", linePos, colPos)
	if HasDown(lineSize, colSize, linePos, colPos) {
		if mat[linePos+1][colPos] != 9 {
			mark[linePos+1][colPos].PartOfBasin = true
		}
	}
}

func markLeftBasin(mat [][]int64, mark [][]PointData, lineSize, colSize, linePos, colPos int) {
	if HasLeft(lineSize, colSize, linePos, colPos) {
		if mat[linePos][colPos-1] != 9 {
			mark[linePos][colPos-1].PartOfBasin = true
		}
	}
}

func markRightBasin(mat [][]int64, mark [][]PointData, lineSize, colSize, linePos, colPos int) {
	if HasRight(lineSize, colSize, linePos, colPos) {
		if mat[linePos][colPos+1] != 9 {
			mark[linePos][colPos+1].PartOfBasin = true
		}
	}
}

func markBasin(mat [][]int64, mark [][]PointData, lineSize, colSize, xMin, yMin int) {
	markDownBasin(mat, mark, lineSize, colSize, xMin, yMin)
	markUpBasin(mat, mark, lineSize, colSize, xMin, yMin)
	markLeftBasin(mat, mark, lineSize, colSize, xMin, yMin)
	markRightBasin(mat, mark, lineSize, colSize, xMin, yMin)
}

func getSurrounding(mark [][]PointData, lineSize, colSize, x, y int) (surr [][2]int) {

	if HasUp(lineSize, colSize, x, y) && !mark[x-1][y].Checked {
		surr = append(surr, [2]int{x-1,y})
	}

	if HasDown(lineSize, colSize, x, y) && !mark[x+1][y].Checked {
		surr = append(surr, [2]int{x+1,y})
	}

	if HasLeft(lineSize, colSize, x, y) && !mark[x][y-1].Checked {
		surr = append(surr, [2]int{x,y-1})
	}

	if HasRight(lineSize, colSize, x, y) && !mark[x][y+1].Checked {
		surr = append(surr, [2]int{x,y+1})
	}

	return
}

type PointData struct {
	Checked bool
	PartOfBasin bool
}

func findBasinSize(xMin, yMin, lineSize, colSize int, mark [][]PointData, mat [][]int64) (size int) {
	//zeroing mark matrix
	for x:=0; x<lineSize; x++ {
		mark = append(mark, []PointData{})
		for y:=0; y<colSize; y++ {
			mark[x] = append(mark[x], PointData{false, false})
		}
	}
	//mark basin center
	mark[xMin][yMin]=PointData{true, true}

	points := map [[2]int]bool{{xMin,yMin}:true}
	for len(points) != 0 {
		//fmt.Println("going for ", points)
		newpoints := map [[2]int]bool{}

		//traverse current points
		for p, _ := range points {
			if mat[p[0]][p[1]] != 9 {
				markBasin(mat, mark, lineSize, colSize, p[0], p[1])
			}
			mark[p[0]][p[1]].Checked = true
			delete(points, p)
			if mat[p[0]][p[1]] != 9 {
				surr := getSurrounding(mark, lineSize, colSize, p[0], p[1])
				for _, pp := range surr {
					newpoints[pp] = true
				}
			}
		}

		//add new
		for p, _ := range newpoints {
			points[p] = true
		}
	}

	return calcBasinSize(mark)
}

func calcRiskLevel(mins []int64) (risk int64) {
	for _, v := range mins {
		risk = risk + v + 1
	}
	return
}

func calcBasinSize(mat [][]PointData) (res int) {
	for x := range mat {
		for y := range mat[x] {
			if mat[x][y].PartOfBasin {
				//fmt.Printf("*")
				res = res + 1
			} else {
				//fmt.Printf(".")
			}
		}
		//fmt.Printf("\n")
	}
	return
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

	heightsMatrix := [][]int64{}
	ln := 0
	for scanner.Scan() {
		ln = ln + 1
		line := scanner.Text()
		for _, n := range line {
			x, err := strconv.ParseInt(string(n), 10, 64)
			if isError(err) {
				return
			}
			if len(heightsMatrix) < ln {
				heightsMatrix = append(heightsMatrix, []int64{})
			}
			heightsMatrix[ln-1] = append(heightsMatrix[ln-1], x)
		}
	}
	//fmt.Println("matrix loaded ", heightsMatrix)
	mins := findLocalMins(heightsMatrix)
	fmt.Println("local mins", mins)
	fmt.Println("risk is", calcRiskLevel(mins))
}

func main() {
	//part2
	file, err := os.OpenFile("input", os.O_RDWR, 0644)

	if isError(err) {
		return
	}

	defer func(){
		err = file.Close()
		_ = isError(err)
	}()

	scanner := bufio.NewScanner(file)

	heightsMatrix := [][]int64{}

	ln := 0
	for scanner.Scan() {
		ln = ln + 1
		line := scanner.Text()
		for _, n := range line {
			x, err := strconv.ParseInt(string(n), 10, 64)
			if isError(err) {
				return
			}
			if len(heightsMatrix) < ln {
				heightsMatrix = append(heightsMatrix, []int64{})
			}
			heightsMatrix[ln-1] = append(heightsMatrix[ln-1], x)
		}
	}
	//fmt.Println("matrix loaded ", heightsMatrix)
	mins := findLocalMinsPositions(heightsMatrix)
	//fmt.Println("local mins", mins)

	basinNr := 0
	sizes := []int{}
	for x := range mins {
		basinNr = basinNr + 1
		mark := [][]PointData{}
		sz := findBasinSize(mins[x][0], mins[x][1], len(heightsMatrix), len(heightsMatrix[0]), mark, heightsMatrix)
		fmt.Printf("basinNr %d with center in %v is %d\n", basinNr, x, sz)
		sizes = append(sizes, sz)
	}

	if len(sizes) < 3 {
		isError(fmt.Errorf("need more basins to calculate 3 best got %d", len(sizes)))
		return
	}
	sort.Ints(sizes)
	acum := 1
	for x := len(sizes) - 1; x > len(sizes) - 4; x-- {
		fmt.Println(sizes[x])
		acum = acum * sizes[x]
	}

	fmt.Println("magic number is", acum)
}