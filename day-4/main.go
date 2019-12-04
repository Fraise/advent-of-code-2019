package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type pos struct {
	x int
	y int
}

func main() {
	f, err := os.Open(os.Args[1])
	defer f.Close()

	if err != nil {
		log.Fatalln(err)
	}

	s := bufio.NewScanner(f)

	//Get 1st wire
	s.Scan()
	r := strings.Split(s.Text(), "-")

	if len(r) != 2 {
		log.Fatalln("Error parsing input!")
	}

	min, err := strconv.Atoi(r[0])

	if err != nil {
		log.Fatalln(err)
	}

	max, err := strconv.Atoi(r[1])

	if err != nil {
		log.Fatalln(err)
	}

	passCount := 0

	for i := min; i <= max; i++ {
		if isValid(i) {
			passCount++
		}
	}

	fmt.Printf("The total number of possible passwords for part 1 is : %d\n", passCount)

	passCount = 0

	for i := min; i <= max; i++ {
		if isValid2(i) {
			passCount++
		}
	}

	fmt.Printf("The total number of possible passwords for part 2 is : %d\n", passCount)
}

func isValid(pass int) bool {
	hasDouble := false
	prev := pass % 10
	pass /= 10

	for pass != 0 {
		digit := pass % 10
		if digit == prev {
			hasDouble = true
		}

		if digit > prev {
			return false
		}

		prev = digit
		pass /= 10
	}

	return hasDouble
}

func isValid2(pass int) bool {
	hasDouble := false
	group := 1
	prev := pass % 10
	pass /= 10

	for pass != 0 {
		digit := pass % 10
		if digit == prev {
			group++
		} else if group == 2 {
			hasDouble = true
			group = 1
		} else {
			group = 1
		}

		if digit > prev {
			return false
		}

		prev = digit
		pass /= 10
	}

	//When the group is the 2 first digits
	if group == 2 {
		hasDouble = true
	}

	return hasDouble
}
