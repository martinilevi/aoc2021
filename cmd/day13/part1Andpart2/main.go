package main

import (
	"bufio"
	"fmt"
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

type Fold struct {
	Coord rune
	Value int64
}

func CreateMatrix(points [][2]int64) (ret [][]rune) {
	maxX := int64(0)
	maxY := int64(0)

	for _, p := range points {
		x := p[0]
		y := p[1]
		if x > maxX {
			maxX = x
		}
		if y > maxY {
			maxY = y
		}
	}

	for y := int64(0); y <= maxY; y++ {
		ret = append(ret, []rune{})
		for x := int64(0); x <= maxX; x++ {
			ret[y] = append(ret[y], '.')
		}
	}

	for _, p := range points {
		ret[p[1]][p[0]] = '#'
	}
	return
}

func CreateEmptyMatrix(lineSize, colSize int) (ret [][]rune) {
	for y := 0; y < lineSize; y++ {
		ret = append(ret, []rune{})
		for x := 0; x < colSize; x++ {
			ret[y] = append(ret[y], '.')
		}
	}
	return
}

func DrawMatrix(in [][]rune) {
	for _, line := range in {
		for _, r := range line {
			fmt.Print(string(r))
		}
		fmt.Println()
	}
}

func CountPoints(in [][]rune) (cnt int) {
	for l, _ := range in {
		for c, _ := range in[l] {
			if in[l][c] == '#' {
				cnt++
			}
		}
	}
	return
}

func FoldLeftByColumn(in [][]rune, value int64) (ret [][]rune) {
	lineSize := len(in)
	colSize :=  len(in[0])

	//left
	left := [][]rune{}
	for l := 0; l < lineSize; l++ {
		left = append(left, []rune{})
		for col := int64(0); col<value; col++ {
			left[l] = append(left[l], in[l][col])
		}
	}

	//fmt.Println("lines x column")
	//fmt.Println("left", len(left), len(left[0]))
	//DrawMatrix(left)
	//value column is erased


	//right
	right := [][]rune{}
	for l := 0; l < lineSize; l++ {
		right = append(right, []rune{})
		for col := value+1; col < int64(colSize); col++ {
			right[l] = append(right[l], in[l][col])
		}
	}

	//fmt.Println("right", len(right), len(right[0]))
	//DrawMatrix(right)

	//fold right	
	rLineSize := len(right)
	rColSize := len(right[0])
	foldedRight := CreateEmptyMatrix(rLineSize, rColSize)
	for y := 0; y < rLineSize; y++ {
		foldedCol := 0
		for x := rColSize - 1; x >= 0; x--{
			foldedRight[y][foldedCol] = right[y][x]
			foldedCol++
		}		
	}

	//fmt.Println("Folded right...")
	//DrawMatrix(foldedRight)
	//fmt.Println("folded right", len(foldedRight), len(foldedRight[0]))

	lLineSize := len(left)
	lColSize := len(left[0])
	if rColSize >= lColSize {
		//fmt.Println("HEY", rColSize, lColSize)
		//fmt.Println("rColsize lColsize", rColSize, lColSize)

		for l := 0; l < lLineSize; l++ {
			leftCol := 0
			for c := rColSize - lColSize; c < rColSize; c++ {
				//fmt.Println("left ", l, leftCol, string(left[l][leftCol]), ">>> foldedRight ", l, c, string(foldedRight[l][c]))
				//merge only if destination is empty
				if foldedRight[l][c] == '.' {
					foldedRight[l][c] = left[l][leftCol]
				}
				leftCol++
			}
		}
		ret = foldedRight
	} else {
		//fmt.Println("HO", rColSize, lColSize)
		for l := 0; l < lLineSize; l++ {
			leftCol := 0
			for c := lColSize - rColSize; c < lColSize; c++ {
				//merge only if destination is empty
				//fmt.Println("folded right ", l, leftCol, ">>> left ", l, c)
				if left[l][c] == '.' {
					left[l][c] = foldedRight[l][leftCol]
				}
				leftCol++
			}
		}
		ret = left
	}
	return
}

func FoldBottomLineUp(in [][]rune, value int64) (ret [][]rune) {
	lineSize := len(in)
	colSize :=  len(in[0])

	//up
	up := [][]rune{}
	for l := int64(0); l < value; l++ {
		up = append(up, []rune{})
		for col := 0; col<colSize; col++ {
			up[l] = append(up[l], in[l][col])
		}
	}

	//value line is erased

	//down
	down := [][]rune{}
	for l := value+1; l < int64(lineSize); l++ {
		down = append(down, []rune{})
		for col := 0; col<colSize; col++ {
			down[l - (value + 1)] = append(down[l - (value + 1)], in[l][col])
		}
	}

	//fold down
	dLineSize := len(down)
	dColSize := len(down[0])
	foldedDown := CreateEmptyMatrix(dLineSize, dColSize)
	foldedLine := 0
	for y := dLineSize - 1; y >=0; y-- {
		for x := 0; x<dColSize; x++{
			foldedDown[foldedLine][x] = down[y][x]
		}
		foldedLine++
	}

	uLineSize := len(up)
	if dLineSize >= uLineSize {
		upLine := 0
		for l := dLineSize - uLineSize; l < dLineSize; l++{
			for c := 0; c<dColSize; c++ {
				//merge only if destination is empty
				if foldedDown[l][c] == '.' {
					foldedDown[l][c] = up[upLine][c]
				}
			}
			upLine++
		}
		ret = foldedDown
	} else {
		downLine := 0
		for l := uLineSize - dLineSize; l < uLineSize; l++{
			for c := 0; c<dColSize; c++ {
				//merge only if destination is empty
				if up[l][c] == '.' {
					up[l][c] = foldedDown[downLine][c]
				}
			}
			downLine++
		}
		ret = up
	}
	return
}

func FoldRunes(in [][]rune, fold *Fold) (ret [][]rune) {
	if fold == nil {
		fmt.Println("nil fold, nop")
		return
	}
	switch fold.Coord {
	case 'x':
		ret = FoldLeftByColumn(in, fold.Value)
	case 'y':
		ret = FoldBottomLineUp(in, fold.Value)
	default:
		fmt.Println("Invalid fold coord", fold.Coord)
	}
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

	ln := 0
	points := [][2]int64{}
	for scanner.Scan() {
		ln = ln + 1
		line := scanner.Text()

		if strings.TrimSpace(line) == "" {
			break
		}

		lineParts := strings.Split(line, ",")
		if len(lineParts) != 2 {
			fmt.Println("Invalid point at line", ln)
			return
		}

		x, err := strconv.ParseInt(lineParts[0], 10, 64)
		if isError(err) {
			return
		}

		y, err := strconv.ParseInt(lineParts[1], 10, 64)
		if isError(err) {
			return
		}

		points = append(points, [2]int64{x,y})
	}

	folds := []*Fold{}
	for scanner.Scan() {
		ln = ln + 1
		line := scanner.Text()

		if strings.TrimSpace(line) == "" {
			break
		}

		line = strings.TrimPrefix(line, "fold along ")

		lineParts := strings.Split(line, "=")
		if len(lineParts) != 2 {
			fmt.Println("Invalid fold at line", ln)
			return
		}

		v, err := strconv.ParseInt(lineParts[1], 10, 64)
		if isError(err) {
			return
		}

		folds = append(folds, &Fold{rune(lineParts[0][0]), v})
	}

	//fmt.Println(points)
	mat := CreateMatrix(points)
	//DrawMatrix(mat)
	c := 0
	for _, f := range folds{
		//fmt.Println("FOLD!!!")
		mat = FoldRunes(mat, f)
		//DrawMatrix(mat)
		if c == 0 {
			//NOTE: for part1
			fmt.Println("points counter ==",CountPoints(mat))
		}
		c++
	}
	DrawMatrix(mat)
}