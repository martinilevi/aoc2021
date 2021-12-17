package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"runtime"
	"sort"
	"strings"
	"time"
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

type FastPair struct {
	Adjacent uint8
	Insert uint8
}

func SearchInsertion(pairs []*Pair, pair string) string {
	for x := 0; x < len(pairs); x++ {
		if pairs[x].Adjacent == pair {
			//fmt.Printf("%s leads to insert %s\n", pair, pairs[x].Insert)
			return pairs[x].Insert
		}
	}
	return ""
}

func FastSearchInsertion(pairs []*FastPair, pair uint8, blockDecode map[uint8]string) uint8 {
	for x := 0; x < len(pairs); x++ {
		if pairs[x].Adjacent == pair {
			//fmt.Println("found pair", blockDecode[pair>>4], blockDecode[pair&0x0f])
			//fmt.Println("Will insert ", blockDecode[pairs[x].Insert>>4], blockDecode[pairs[x].Insert&0x0f] )
			return pairs[x].Insert
		}
	}
	fmt.Println("pair", blockDecode[pair>>4], blockDecode[pair&0x0f], "not found",)
	panic("fast insertion not found")
	return 0xff
}

func Step(pairs []*Pair, qpairs []*Pair, polymer string) string {
	var tmp strings.Builder
	//fmt.Println("'"+polymer+"'")
	start := time.Now()
	totalSearchTime := time.Duration(0)

	totalConcatTime := time.Duration(0)
	searches := 0

	work := polymer[0:2]
	//fmt.Printf("\t work with [%d:%d] %s\n", 0, 2, work)
	insert := SearchInsertion(pairs, work)
	work = string(work[0]) + string(insert[0]) + string(work[1])
	tmp.WriteString(work)
	//fmt.Println("\t"+tmp.String())

	for x := 1; x+1 < len(polymer); x++ {
		//fmt.Printf("\t work with [%d:%d] %s\n", x, x+2, polymer[x:x+2])
		startSearch := time.Now()
		searches++
		insert := SearchInsertion(qpairs, polymer[x:x+2])
		deltaSearch := time.Since(startSearch)
		totalSearchTime = totalSearchTime + deltaSearch
		startConcat := time.Now()
		tmp.WriteString(insert)
		deltaConcat := time.Since(startConcat)
		totalConcatTime = totalConcatTime + deltaConcat
		//fmt.Println("\t"+tmp.String())
	}

	fmt.Printf("TOTAL iterations in this step %d\n", searches)
	fmt.Printf("TOTAL time used in searchs %s\n", totalSearchTime)
	fmt.Printf("TOTAL time used in concats %s\n", totalConcatTime)
	fmt.Printf("TOTAL STEP TIME: %s\n", time.Since(start).Round(time.Second))
	return tmp.String()
}

func HighPosition(b uint8) uint8 {
	return (b & 0xf0) >> 4
}

func LowPosition(b uint8) uint8 {
	return b & 0x0f
}

func GetMidArray(idx int, arr []uint8) uint8 {
	//               H L  H L
	//MID POSITION   0 1  2 3 4 5 6 7
	//BYTE POSITION |  0||  1|  2|  3|

	if idx == 0 {
		return HighPosition(arr[0])
	}

	isEven := (idx % 2 == 0)

	if ! isEven {
		return LowPosition(arr[int(idx/2)])
	}

	return HighPosition(arr[idx/2])
}

func LenMidArray(arr []uint8) int {
	if len(arr) == 0 {
		return 0
	}
	last := arr[len(arr)-1]
	if last&0x0f == 0 {
		return 2*len(arr)-1
	}
	return 2*len(arr)
}

func AppendMidArr(arr []uint8, midByte uint8) []uint8 {
	//high part of midbyte is not used
	//just lower 4 bits are used

	if len(arr) == 0 {
		arr = append(arr, midByte<<4)
		return arr
	}

	last := arr[len(arr)-1]

	//if low not empty
	if last & 0x0F != 0 {
		arr = append(arr, midByte<<4)
		return arr
	}

	//low empty!, just add as low
	arr[len(arr)-1] = last + (midByte&0x0f)

	return arr
}


func FastStep(pairs []*FastPair, qpairs []*FastPair, polymer []uint8, blocksDecode map[uint8]string) []uint8 {
	//fmt.Println("'"+polymer+"'")
	start := time.Now()
	totalSearchTime := time.Duration(0)

	totalConcatTime := time.Duration(0)
	searches := 0

	work := polymer[0]
	//fmt.Printf("\t work with [%d:%d] %s\n", 0, 2, work)
	insert := FastSearchInsertion(pairs, work, blocksDecode)

	head := work&0xF0
	tail := work&0x0F
	work = head + insert
	tmp := []uint8{work}
	tmp = AppendMidArr(tmp, tail)
	/*
	fmt.Println("head is",   blocksDecode[head>>4])
	fmt.Println("insert is", blocksDecode[insert])
	fmt.Println("tail is",   blocksDecode[tail])
	for x := 0; x < LenMidArray(tmp); x++ {
		fmt.Print(blocksDecode[GetMidArray(x,tmp)])
	}
	fmt.Println()*/

	//fmt.Println("polymer len is", LenMidArray(polymer))
	for x := 1; x + 1 < LenMidArray(polymer); x++ {
		startSearch := time.Now()
		searches++
		insert := FastSearchInsertion(qpairs, GetMidArray(x, polymer) * 16 + GetMidArray(x+1, polymer), blocksDecode)
		deltaSearch := time.Since(startSearch)
		totalSearchTime = totalSearchTime + deltaSearch
		startConcat := time.Now()
		tmp = AppendMidArr(tmp, insert>>4)
		tmp = AppendMidArr(tmp, insert)
		/*
		for x := 0; x < LenMidArray(tmp); x++ {
			fmt.Print(blocksDecode[GetMidArray(x,tmp)])
		}
		fmt.Println()*/
		deltaConcat := time.Since(startConcat)
		totalConcatTime = totalConcatTime + deltaConcat
		//fmt.Println("\t"+tmp.String())
	}

	fmt.Printf("TOTAL iterations in this step %d\n", searches)
	fmt.Printf("TOTAL time used in searchs %s\n", totalSearchTime)
	fmt.Printf("TOTAL time used in concats %s\n", totalConcatTime)
	fmt.Printf("TOTAL STEP TIME: %s\n", time.Since(start).Round(time.Second))

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

func PrintByteArr(arr []uint8) {
	for _, v := range arr {
		fmt.Printf("0x%02x ", v)
	}
	fmt.Println()
}

func DecodeByteArr(arr []uint8, blocksDecode map[uint8]string) (ret string){

	for x := 0; x < LenMidArray(arr); x++ {
		b := GetMidArray(x, arr)
		ret = ret + blocksDecode[b]
	}

	return
}

func MidArrayTest() {
	a := []uint8{}
	PrintByteArr(a)
	a = AppendMidArr(a, 0xFF)
	PrintByteArr(a)
	a = AppendMidArr(a, 0xFF)
	PrintByteArr(a)
	a = AppendMidArr(a, 0xFF)
	PrintByteArr(a)
	a = AppendMidArr(a, 0x01)
	PrintByteArr(a)
	a = AppendMidArr(a, 0x01)
	PrintByteArr(a)
	a = AppendMidArr(a, 0x0F)
	PrintByteArr(a)

	for x := 0; x < LenMidArray(a); x++ {
		fmt.Printf("0x%02x ", GetMidArray(x, a))
	}
	fmt.Println()
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

	qpairs := []*Pair{}
	pairs := []*Pair{}
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
		pairs = append(pairs, &Pair{lineParts[0], lineParts[1]})
		qpairs = append(qpairs, &Pair{lineParts[0], lineParts[1]+string(lineParts[0][1])})
	}

	blocks := map[string]bool{}
	for x := range pairs {
		//fmt.Println(*pairs[x])
		blocks[pairs[x].Insert] = true
	}

	if len(blocks) > 15 {
		fmt.Println("more than 15 blocks, need other approach")
	}

	blocksOrdered := []string{}
	for x, _ := range blocks {
		blocksOrdered = append(blocksOrdered, x)
	}
	sort.Strings(blocksOrdered)

	blocksEncode := map[string]uint8{}
	blocksDecode := map[uint8]string{}
	count := uint8(1)
	for _, x := range blocksOrdered {
		blocksEncode[x] = count
		blocksDecode[count] = x
		count ++
		//fmt.Printf("%s %x\n", x, blocksEncode[x])
	}

	pairsEncode := map[string]uint8{}

	for x, _ := range pairs {
		pairsEncode[pairs[x].Adjacent] = blocksEncode[string(pairs[x].Adjacent[0])] * 16 + blocksEncode[string(pairs[x].Adjacent[1])]
		//fmt.Printf("%s 0x%02x(%d)\n", pairs[x].Adjacent, pairsEncode[pairs[x].Adjacent], pairsEncode[pairs[x].Adjacent])
	}

	qpairsEncode := map[string]uint8{}
	for x, _ := range qpairs {
		qpairsEncode[qpairs[x].Adjacent] = blocksEncode[string(qpairs[x].Adjacent[0])] * 16 + blocksEncode[string(qpairs[x].Adjacent[1])]
		//fmt.Printf("%s 0x%02x(%d)\n", qpairs[x].Adjacent, qpairsEncode[pairs[x].Adjacent], qpairsEncode[pairs[x].Adjacent])
	}

	fastPairs := []*FastPair{}
	for _, p := range pairs {
		fastPairs = append(fastPairs, &FastPair{pairsEncode[p.Adjacent], blocksEncode[p.Insert]})
	}

	/*
	for _, p := range fastPairs {
		//fmt.Printf("%02x leads to %02x\n", p.Adjacent, p.Insert)
		fmt.Printf("%s leads to %s\n", blocksDecode[p.Adjacent>>4] + blocksDecode[p.Adjacent&0x0f] , blocksDecode[p.Insert])
	}*/

	fastQPairs := []*FastPair{}
	for _, p := range qpairs {
		fastQPairs = append(fastQPairs, &FastPair{pairsEncode[p.Adjacent], blocksEncode[string(p.Insert[0])]*16+blocksEncode[string(p.Insert[1])]})
	}

	/*for _, p := range fastQPairs {
		//fmt.Printf("%02x leads to %02x\n", p.Adjacent, p.Insert)
		fmt.Printf("%s leads to %s\n",
			blocksDecode[p.Adjacent>>4] + blocksDecode[p.Adjacent&0x0f] ,
			blocksDecode[p.Insert>>4] + blocksDecode[p.Insert&0x0f])
	}*/

	fastPolymer := []uint8{}
	for x := 0; x < len(polymer); x = x + 2 {
		b := blocksEncode[string(polymer[x])] * 16 + blocksEncode[string(polymer[x+1])]
		fastPolymer = append(fastPolymer, b )
		//fmt.Printf("polymer 0x%02x\n", fastPolymer)
	}

	s := fastPolymer
	fmt.Println(DecodeByteArr(s, blocksDecode))
	for step := 0; step < 40; step++ {
		s = FastStep(fastPairs, fastQPairs, s, blocksDecode)
		length := LenMidArray(s)
		fmt.Printf("STEP %d - len[%d] (%d Mbytes)\n", step+1, length, int(len(s)/1024/1024))
		if length < 100 {
			fmt.Println(DecodeByteArr(s, blocksDecode))
		}
	}

}