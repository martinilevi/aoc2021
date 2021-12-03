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

	fmt.Println("Assuming 1 for undefined case")
	return true
}

func calc(cnt [reportBitSize]int64, lineCnt int64, bitFunc func(particular, total int64) bool) (int64, error) {
	tmp := ""
	for x := 0; x < reportBitSize; x++ {
		if bitFunc(cnt[x], lineCnt) {
			tmp = tmp + "1"
			continue
		}
		tmp = tmp + "0"
	}

	value, err := strconv.ParseInt(tmp, 2, 16)

	if err != nil {
		return -1, err
	}

	return value, nil
}

func calcGamma(cnt [reportBitSize]int64, lineCnt int64) (int64, error) {
	return calc(cnt, lineCnt, mostCommonBit)
}

func calcEpsilon(cnt [reportBitSize]int64, lineCnt int64) (int64, error) {
	return calc(cnt, lineCnt, leastCommonBit)
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
	for scanner.Scan() {
		lineCnt = lineCnt + 1
		bitLine := scanner.Text()
		for x := 0; x < reportBitSize; x++ {
			bit, err := strconv.ParseInt(string(bitLine[x]), 2, 2)
			if isError(err) {
				return
			}
			cnt[x] = cnt[x] + bit
		}
	}

	gamma, err := calcGamma(cnt, lineCnt)
	if isError(err) {
		return
	}

	epsilon, err := calcEpsilon(cnt, lineCnt)
	if isError(err) {
		return
	}

	fmt.Printf("gamma=%d epsilon=%d product=%d\n", gamma, epsilon, gamma*epsilon)
}

