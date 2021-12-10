package main

import (
	"bufio"
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"
)

const signalCount = 10

func isError(err error) bool {
	if err != nil {
		_, filename, line, _ := runtime.Caller(1)
		fmt.Printf("%s:%d %s\n", filename, line, err.Error())
	}
	return (err != nil)
}

func hasCandF(cf string, input string) bool {
	cnt := 0
	for _, w := range cf {
		for _, v := range input {
			if v == w {
				cnt = cnt + 1
			}
		}
	}
	return cnt == 2
}

func GetF(cf string, six string) int32 {
	for _, w := range cf {
		for _, v := range six {
			if v == w {
				return v
			}
		}
	}
	return -1
}

func GetC(cf string, f int32) int32 {
	for _, w := range cf {
		if w != f {
			return w
		}
	}
	return -1
}

func HasAllLettersFrom(candidate, group string) bool {
	found := 0
	for _, v := range group {
		for _, w := range candidate {
			if w == v {
				found = found + 1
			}
		}
	}
	return found == len(group)
}


func GuessSignals(input []string) (signal [signalCount]string) {
	for _, s := range input {
		switch len(s) {
		case 2:
			signal[1] = s //has real cf
		case 3:
			signal[7] = s //has real acf NOT b, d, e, g
		case 4:
			signal[4] = s //has real bcdf NOT a, e, g
		case 5:
			//may be 2 (aCdeg)
			//may be 3 (aCdFg) // en la segunda pasada sale
			//may be 5 (abdFg)
		case 6:
			//may be 0 (abCeFg)
			//may be 6 (abdeFg) // en la segunda pasada sale
			//may be 9 (abCdFg)
		case 7:
			signal[8] = s //has all
		}
	}

	F := int32(0)
	C := int32(0)

	candidate251 := ""
	candidate252 := ""

	candidate091 := ""
	candidate092 := ""

	cf := signal[1]
	for _, s := range input {
		switch len(s) {
		case 5:
			if hasCandF(cf, s) {
				signal[3] = s
				continue
			}

			//may be 2 (aCdeg)
			//may be 5 (abdFg)
			if candidate251 == "" {
				candidate251 = s
				continue
			} else {
				candidate252 = s
			}

		case 6:
			if !hasCandF(cf, s) {
				signal[6] = s
				F = GetF(cf, s)
				C = GetC(cf, F)
				continue
			}
			//may be 0 (abCeFg)
			//may be 9 (abCdFg)
			if candidate091 == "" {
				candidate091 = s
				continue
			} else {
				candidate092 = s
			}
		}
	}

	if strings.Contains(candidate251, string(C)) {
		signal[2] = candidate251
		signal[5] = candidate252
	} else {
		signal[2] = candidate252
		signal[5] = candidate251
	}

	if HasAllLettersFrom(candidate091,signal[5]) {
		signal[9] = candidate091
		signal[0] = candidate092
	} else {
		signal[9] = candidate092
		signal[0] = candidate091
	}
	return
}

func ReadOutput(guessed [signalCount]string, input []string) (output string) {
	for _, v := range input {
		for idx, w := range guessed {
			if len(v) != len(w) {
				continue
			}
			if HasAllLettersFrom(v, w) {
				switch idx {
				case 0:
					output = output + "0"
				case 1:
					output = output + "1"
				case 2:
					output = output + "2"
				case 3:
					output = output + "3"
				case 4:
					output = output + "4"
				case 5:
					output = output + "5"
				case 6:
					output = output + "6"
				case 7:
					output = output + "7"
				case 8:
					output = output + "8"
				case 9:
					output = output + "9"
				}

				break
			}
		}
	}
	return
}

func count1478(input string) (cnt int) {
	for _, v := range input {
		switch v {
		case '1':
			cnt++
		case '4':
			cnt++
		case '7':
			cnt++
		case '8':
			cnt++
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

	acc := int64(0)
	count := 0
	ln := 0
	for scanner.Scan() {
		ln := ln + 1
		line := scanner.Text()

		fuelVecStr := strings.Split(line, "|")

		if len(fuelVecStr) != 2 {
			isError(fmt.Errorf("line with invalid format, line %d", ln))
			return
		}

		//
		signalsStr := fuelVecStr[0]
		signals := strings.Split(signalsStr, " ")
		signals = signals[0:len(signals)-1]
		if len(signals) != 10 {
			fmt.Println(signals)
			isError(fmt.Errorf("invalid nr of signals %d, line %d", len(signals), ln))
			return
		}

		//4 digit output
		outputStr := fuelVecStr[1]
		output := strings.Split(outputStr, " ")
		output = output[1:]
		if len(output) != 4 {
			fmt.Println(">>>>>>>>>>>>>>>>>", output)
			isError(fmt.Errorf("invalid nr of outputs, line %d", ln))
			return
		}

		guessedSignals := GuessSignals(signals)
		//fmt.Println("UNGUESSED", signals)
		//fmt.Println("RAW READ", output)
		fmt.Println("GUESSED", guessedSignals)
		lineRead := ReadOutput(guessedSignals, output)
		fmt.Println("READ ", lineRead)

		count = count + count1478(lineRead)

		clean := strings.TrimLeft(lineRead, "0")

		if len(clean) == 0 {
			fmt.Println("0000 skipped")
			continue
		}

		x, err := strconv.ParseInt(clean, 0, 64)
		if isError(err) {
			return
		}
		acc = acc + x
	}

	fmt.Printf("part1: 1478 appear %d times\n", count)
	fmt.Printf("part2: totalSum is %d \n", acc)
}