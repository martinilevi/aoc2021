package main

import (
	"bufio"
	"fmt"
	"os"
	"runtime"
	"strconv"
)

const reportBitSize = 12

func isError(err error) bool {
	if err != nil {
		_, filename, line, _ := runtime.Caller(1)
		fmt.Printf("%s:%d %s\n", filename, line, err.Error())
	}
	return (err != nil)
}

func leastCommonBit(particular, total int64) bool {
	return !mostCommonBit(particular, total)
}

func mostCommonBit(particular, total int64) bool {
	ratio := float64(particular) / float64(total)
	switch {
	case ratio < 0.5 :
		//0 is the most common bit
		return false
	case ratio > 0.5:
		//1 is the most common bit
		return true
	}
	//assuming 1 for equally common case (in least common bit will assume 0)
	return true
}

func calc(cnt [reportBitSize]int64, lineCnt int64, bitFunc func(particular, total int64) bool) (int64, error) {
	tmp := calcStr(cnt, lineCnt, bitFunc)

	value, err := strconv.ParseInt(tmp, 2, 16)

	if err != nil {
		return -1, err
	}

	return value, nil
}

func calcStr(cnt [reportBitSize]int64, lineCnt int64, bitFunc func(particular, total int64) bool) (string) {
	tmp := ""
	for x := 0; x < reportBitSize; x++ {
		if bitFunc(cnt[x], lineCnt) {
			tmp = tmp + "1"
			continue
		}
		tmp = tmp + "0"
	}

	return tmp
}

func passesFilter(input string, filter string, bit int8) bool {
	if input[bit] == filter[bit] {
		return true
	}
	return false
}

func mapOnlyValue(m map[string]bool) string {
	if len(m) != 1 {return ""}
	k := ""
	for k, _ = range m {}
	return k
}

func searchForNumberWithFilter(bitSet map[string]bool, bitFunc func(particular, total int64) bool) (string, bool) {
	for x := int8(0); x < reportBitSize; x++ {
		cnt, lineCnt := getCntAndLineCnt(bitSet)
		filter := calcStr(cnt, int64(lineCnt), bitFunc)
		//fmt.Println("Will filter using ", filter)
		for bl, _ := range bitSet {
			if !passesFilter(bl, filter, x) {
				delete(bitSet, bl)
				//fmt.Println("filtering", bl)
				//fmt.Println("bitSet remains", bitSet)
			}
			switch {
			case len(bitSet) == 1:
				//fmt.Println("value found!", mapOnlyValue(bitSet))
				return mapOnlyValue(bitSet), true
			case len(bitSet) == 0:
				//fmt.Println("value NOT found! empty set")
				return "", false
			}
		}
		//fmt.Println("NEXT BIT")
	}
	//fmt.Println("value NOT found! many options (%d)", len(bitSet))
	return "", false
}

func getCntAndLineCnt(m map[string]bool) (cnt [reportBitSize]int64, lineCnt int) {
	for bitLine, _ := range m {
		for x := 0; x < reportBitSize; x++ {
			bit, err := strconv.ParseInt(string(bitLine[x]), 2, 2)
			if isError(err) {
				return
			}
			cnt[x] = cnt[x] + bit
		}
	}
	return cnt, len(m)
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

	cnt := [reportBitSize]int64{}
	lineCnt := int64(0)
	bitSetO2 := map[string]bool{}
	bitSetCO2 := map[string]bool{}

	for scanner.Scan() {
		lineCnt = lineCnt + 1
		bitLine := scanner.Text()
		bitSetO2[bitLine] = true
		bitSetCO2[bitLine] = true

		for x := 0; x < reportBitSize; x++ {
			bit, err := strconv.ParseInt(string(bitLine[x]), 2, 2)
			if isError(err) {
				return
			}
			cnt[x] = cnt[x] + bit
		}
	}

	O2ValueStr, found := searchForNumberWithFilter(bitSetO2, mostCommonBit)
	if !found {
		fmt.Println("O2 Value not found")
		return
	}

	CO2ValueStr, found := searchForNumberWithFilter(bitSetCO2, leastCommonBit)
	if !found {
		fmt.Println("CO2 Value not found")
		return
	}

	O2Value, err := strconv.ParseInt(O2ValueStr, 2, 16)
	if isError(err) {
		return
	}

	CO2Value, err := strconv.ParseInt(CO2ValueStr, 2, 16)
	if isError(err) {
		return
	}

	fmt.Printf("o2=%d co2=%d product=%d\n", O2Value, CO2Value, O2Value*CO2Value)
}