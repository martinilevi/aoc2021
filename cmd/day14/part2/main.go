package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"runtime"
	"strings"
)

func isError(err error) bool {
	if err != nil {
		_, filename, line, _ := runtime.Caller(1)
		fmt.Printf("%s:%d %s\n", filename, line, err.Error())
	}
	return (err != nil)
}

type Pair struct {
	Adjacent string
	Insert string
}

func SearchInsertion(pairs []*Pair, pair string) string {
	for x := 0; x < len(pairs); x++ {
		if pairs[x].Adjacent == pair {
			return pairs[x].Insert
		}
	}
	return ""
}

func Step(pairs []*Pair, polymer string) string {
	tmp := ""
	//fmt.Println("'"+polymer+"'")
	for x := 0; x+1 < len(polymer); x++ {
		work := polymer[x:x+2]
		//fmt.Printf("\t work with [%d:%d] %s\n", x, x+2, work)
		insert := SearchInsertion(pairs, work)
		if insert != "" {
			//fmt.Println("inserting "+insert)
			if tmp == "" {
				work = string(work[0]) + string(insert[0]) + string(work[1])
			} else {
				work = string(insert[0]) + string(work[1])
			}

		}
		tmp = tmp + work
		//fmt.Println("\t"+tmp)
	}
	return tmp
}

func CountChars(polymer string) (cnt map[int32]int64) {
	cnt = map[int32]int64{}
	for _, c := range(polymer) {

		if _, ok := cnt[c]; !ok {
			cnt[c] = 1
			continue
		}
		cnt[c]++
	}
	return
}

func PrintCharMap(cnt map[int32]int64) () {
	for k, v := range cnt  {
		fmt.Printf("%s : %d\n", string(k), v)
	}
	return
}

func FindMaxAndMin(cnt map[int32]int64) (max, min int64, maxChar, minChar int32) {
	max = math.MinInt64
	maxChar = int32(0)

	min = math.MaxInt64
	minChar = int32(0)

	for k, v := range cnt {
		if v > max {
			max = v
			maxChar = k
		}
		if v < min {
			min = v
			minChar = k
		}
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
	polymer := ""
	for scanner.Scan() {
		ln = ln + 1
		line := scanner.Text()

		if strings.TrimSpace(line) == "" {
			break
		}

		polymer = line
	}

	fmt.Println("Going for...", polymer)

	blocks := map[string]bool{}
	pairs := map[string]string{}
	for scanner.Scan() {
		ln = ln + 1
		line := scanner.Text()
		lineParts := strings.Split(line, " -> ")
		if len(lineParts) != 2 {
			fmt.Println("Invalid pair mutation at line", ln)
			return
		}
		if len(lineParts[0]) != 2 {
			fmt.Println("Invalid pair mutation source length at line", ln)
			return
		}
		if len(lineParts[1]) != 1 {
			fmt.Println("Invalid pair mutation destination length at line", ln)
			return
		}
		fmt.Println(lineParts[0], " -> ", lineParts[1])
		pairs[lineParts[0]] = lineParts[1]
		blocks[lineParts[1]] = true
	}

	fmt.Println("blocks...")
	for k, _ := range blocks{
		fmt.Print(k, " ")
	}
	fmt.Println()

	count := map[string]int64{}
	for x := 0; x + 1 < len(polymer); x++ {
		count[polymer[x:x+2]]++
	}

	for step := 0; step < 40; step++ {
		recount := map[string]int64{}
		//s = Step(pairs, s)
		//length := len(s)
		fmt.Printf("STEP %d\n", step+1)
		for p, _ := range count {
			recount[string(p[0])+pairs[p]] += count[p]
			recount[pairs[p]+string(p[1])] += count[p]
		}
		count = recount
		if step == 39 {
			charCount := map[uint8]int64{}
			for k, v := range recount {
				//let's count just the first letter of each pair as it is repeated
				//last letter will always be the last letter of the original polymer
				//we add it in the end
				charCount[k[0]] += v
			}
			charCount[polymer[len(polymer)-1]] += 1

			max := int64(0)
			maxC := uint8(0)
			min := int64(math.MaxInt64)
			minC := uint8(0)
			for c, v := range charCount {
				if int64(v) < min {
					min = v
					minC = c
				}
				if v > max {
					max = v
					maxC = c
				}
			}

			fmt.Print("Char count \n")
			for k, v := range charCount {
				fmt.Printf("%s %d\n", string(k), v)
			}
			fmt.Printf("max is %s with %d\n", string(maxC), max)
			fmt.Printf("min is %s with %d\n", string(minC), min)
			fmt.Printf("max-min = %d\n", max - min)
		}
	}


}