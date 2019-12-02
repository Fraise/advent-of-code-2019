package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {
	f, err := os.Open(os.Args[1])
	defer f.Close()

	if err != nil {
		log.Fatalln(err)
	}

	s := bufio.NewScanner(f)

	s.Scan()
	strProgram := strings.Split(s.Text(), ",")

	var program []int

	for _, s := range strProgram {
		i, err := strconv.Atoi(s)

		if err != nil {
			log.Fatalln(err)
		}

		program = append(program, i)
	}

	//Part 1

	//Restore the gravity assist program to the "1202 program alarm" state
	program[1] = 12
	program[2] = 2

	baseProgram := copyProgram(program)

	fmt.Println("Part 1 :")
	fmt.Printf("Program value @ pos 0 : %d\n", compute(program))

	//Part 2
	//Bruteforcing the solution...
	for i := 0; i < 1000; i++ {
		for j := 0; j < 1000; j++ {
			resetProgram(program, baseProgram)
			program[1] = i
			program[2] = j
			if compute(program) == 19690720 {
				fmt.Println("Part 2 :")
				fmt.Printf("noun : %d\nverb : %d\n", i, j)
				fmt.Printf("Answer : %d\n", 100*i+j)
				break
			}
		}
	}
}

func compute(program []int) int {
	i := 0

	for i < len(program) {
		if i+3 >= len(program) {
			//fmt.Println("Tried to access a value outside of the program!")
			return -1
		}

		op := program[i]

		if op == 99 {
			//fmt.Println("op 99 : Program halted.")
			break
		}

		val1 := program[i+1]
		val2 := program[i+2]
		res := program[i+3]

		if val1 >= len(program) || val2 >= len(program) || res >= len(program) {
			//fmt.Println("Tried to access a value outside of the program!")
			return -1
		}

		if op == 1 {
			program[res] = program[val1] + program[val2]
		} else if op == 2 {
			program[res] = program[val1] * program[val2]
		} else {
			fmt.Println("Unknown op : Program crashed.")
		}

		i += 4
	}

	return program[0]
}

func resetProgram(program []int, baseProgram []int) []int {
	for i := 0; i < len(program) || i < len(baseProgram); i++ {
		program[i] = baseProgram[i]
	}

	return program
}

func copyProgram(program []int) (copy []int) {
	for _, v := range program {
		copy = append(copy, v)
	}

	return copy
}
