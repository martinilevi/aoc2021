package main

import (
	"bufio"
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"
)

const cardSize = 5

func isError(err error) bool {
	if err != nil {
		_, filename, line, _ := runtime.Caller(1)
		fmt.Printf("%s:%d %s\n", filename, line, err.Error())
	}
	return (err != nil)
}

type Card [cardSize][cardSize]int64 // row column
type Memory [cardSize][cardSize]bool // row column
type MemoryCard struct {
	Kard Card
	Seen Memory
	Won bool
}

func WinningRow(number int64, c *MemoryCard) bool {
	for row := 0; row < cardSize; row++ {
		matches := 0
		for x := 0; x < cardSize; x++{
			if c.Seen[row][x] {
				//number already seen, count
				matches = matches + 1
			} else if c.Kard[row][x] == number {
				//new match, mark and count
				c.Seen[row][x] = true
				matches = matches + 1
			}
		}
		if matches == cardSize {
			return true
		}
	}
	return false
}

func WinningColumn(number int64, c* MemoryCard) bool {
	for col := 0; col < cardSize; col++ {
		matches := 0
		for x := 0; x < cardSize; x++{
			if c.Seen[x][col] {
				//number already seen, count
				matches = matches + 1
			} else if c.Kard[x][col] == number {
				//new match, mark and count
				c.Seen[x][col] = true
				matches = matches + 1
			}
		}
		if matches == cardSize {
			return true
		}
	}
	return false
}

func PrintCard(card *MemoryCard) {
	for x:=0; x < cardSize; x++{
		for y:=0; y < cardSize; y++{
			fmt.Printf("%d ",card.Kard[x][y])
		}
		fmt.Printf("\n")
	}
	for x:=0; x < cardSize; x++{
		for y:=0; y < cardSize; y++{
			fmt.Printf("%t ",card.Seen[x][y])
		}
		fmt.Printf("\n")
	}
}

func WinningScore(winNr int64, card *MemoryCard) (score int64) {
	score = int64(0)
	for x:=0; x < cardSize; x++{
		for y:=0; y < cardSize; y++{
			if ! card.Seen[x][y] {
				score = score + card.Kard[x][y]
			}
		}
	}
	score = score * winNr
	return
}

func Bingo(numbers []int64, cards []*MemoryCard) (winner bool, winNr int64, winCard *MemoryCard) {
	for order, n := range numbers {
		for i, c := range cards {
			if WinningRow(n, c) {
				fmt.Printf("On number %d (order %d): Winning row on card %d\n", n, order+1, i )
				PrintCard(c)
				return true, n, c
			}
			if WinningColumn(n, c) {
				fmt.Printf("On number %d (order %d): Winning column on card %d\n", n, order+1, i )
				PrintCard(c)
				return true, n, c
			}
		}
	}
	return false, -1, nil
}

func LastCardBingo(numbers []int64, cards []*MemoryCard) (winner bool, winNr int64, winCard *MemoryCard) {
	winNr = -1
	for order, n := range numbers {
		for i, c := range cards {
			if c.Won {
				//a card only wins once
				continue
			}
			if WinningRow(n, c) {
				fmt.Printf("On number %d (order %d): Winning row on card %d\n", n, order+1, i )
				c.Won = true
				PrintCard(c)
				winner = true
				winNr = n
				winCard = c
				continue
			}
			if WinningColumn(n, c) {
				fmt.Printf("On number %d (order %d): Winning column on card %d\n", n, order+1, i )
				c.Won = true
				PrintCard(c)
				winner = true
				winNr = n
				winCard = c
				continue
			}
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

	//read header
	ln := 0
	numbers := []int64{}
	partial := 0
	complete := 0
	cards := []*MemoryCard{}
	current := Card{}

	for scanner.Scan() {
		ln = ln + 1
		line := scanner.Text()
		fmt.Printf("line(%d) partial(%d) '%s'\n",ln, partial, line)
		if ln == 1 {
			numberStr := strings.Split(line, ",")
			for x := 0; x < len(numberStr); x++ {
				nr, err := strconv.ParseInt(numberStr[x], 10, 64)
				if isError(err) { return }
				numbers = append(numbers, nr)
			}
			continue
		}

		if line == "" {
			if ln > 2 {
				complete = complete + 1
				if partial != cardSize {
					isError(fmt.Errorf("card with invalid lines (%d) at line %d", partial, ln))
					return
				}
				mc := MemoryCard{ Kard: current, Seen: [cardSize][cardSize]bool{}}
				cards = append(cards, &mc)
			}
			partial = 0
			continue
		}

		if partial == cardSize {
			isError(fmt.Errorf("card with over 5 lines at line %d", ln))
			return
		}

		line = strings.TrimSpace(line)
		line = strings.ReplaceAll(line, "  ", " ")
		numberStr := strings.Split(line, " ")
		for x := 0; x < len(numberStr); x++ {
			nr, err := strconv.ParseInt(numberStr[x], 10, 64)
			if isError(err) { return }
			current[partial][x] = nr
		}

		partial = partial + 1
	}

	complete = complete + 1
	if partial != cardSize {
		isError(fmt.Errorf("card with invalid lines (%d) at line %d", partial, ln))
		return
	}
	mc := MemoryCard{ Kard: current, Seen: [cardSize][cardSize]bool{}}
	cards = append(cards, &mc)

	winner, winNr, winCard := LastCardBingo(numbers, cards)
	if winner {
		fmt.Println("bingo had a winner!")
		fmt.Println("SCORE: ", WinningScore(winNr, winCard))
	}
}