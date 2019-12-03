package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
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
	w1 := strings.Split(s.Text(), ",")
	w1m := make(map[pos]int)

	//Get 2nd wire
	s.Scan()
	w2 := strings.Split(s.Text(), ",")
	w2m := make(map[pos]int)

	var sec []pos
	var currPos pos = pos{0, 0}
	steps := 0

	for _, s := range w1 {
		dir := s[:1]
		val, err := strconv.Atoi(s[1:])

		if err != nil {
			log.Fatalln(err)
		}

		if dir == "R" {
			for i := 1; i <= val; i++ {
				steps++
				currPos.x++
				w1m[currPos] = steps
			}
		} else if dir == "L" {
			for i := 1; i <= val; i++ {
				steps++
				currPos.x--
				w1m[currPos] = steps
			}
		} else if dir == "U" {
			for i := 1; i <= val; i++ {
				steps++
				currPos.y++
				w1m[currPos] = steps
			}
		} else if dir == "D" {
			for i := 1; i <= val; i++ {
				steps++
				currPos.y--
				w1m[currPos] = steps
			}
		}
	}

	currPos = pos{0, 0}
	steps = 0

	for _, s := range w2 {
		dir := s[:1]
		val, err := strconv.Atoi(s[1:])

		if err != nil {
			log.Fatalln(err)
		}

		if dir == "R" {
			for i := 1; i <= val; i++ {
				steps++
				currPos.x++
				w2m[currPos] = steps
			}
		} else if dir == "L" {
			for i := 1; i <= val; i++ {
				steps++
				currPos.x--
				w2m[currPos] = steps
			}
		} else if dir == "U" {
			for i := 1; i <= val; i++ {
				steps++
				currPos.y++
				w2m[currPos] = steps
			}
		} else if dir == "D" {
			for i := 1; i <= val; i++ {
				steps++
				currPos.y--
				w2m[currPos] = steps
			}
		}
	}

	for k := range w1m {
		if w2m[k] > 0 {
			sec = append(sec, k)
		}
	}

	var min int = math.MaxInt32

	for _, p := range sec {
		val := int(math.Abs(float64(p.x))) + int(math.Abs(float64(p.y)))
		if val < min {
			min = val
		}
	}

	fmt.Printf("The closest intersection if at distance : %d\n", min)

	var minSteps int = math.MaxInt32

	for k := range w1m {
		if w2m[k] > 0 {
			sumSteps := w1m[k] + w2m[k]
			if sumSteps < minSteps {
				minSteps = sumSteps
			}
		}
	}

	fmt.Printf("The minimum sum of steps to an intersection is : %d\n", minSteps)
}
