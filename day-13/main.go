package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

type pos struct {
	lin int
	col int
}

var dir = [4]string{"up", "right", "down", "left"}

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

	extraMem := make([]int, 1000, 1000)

	program = append(program, extraMem...)
	baseProgram := copyProgram(program)

	//Part 1

	screen := make([][]int, 40)

	for i := 0; i < len(screen); i++ {
		screen[i] = make([]int, 40)
	}

	input := make(chan int, 100)
	output := make(chan int, 100)

	go func() {
		ret, _ := compute(program, input, output)

		output <- ret
	}()

	for {
		col := <-output

		if col == -99 {
			//End of game
			break
		}

		lin := <-output
		elmt := <-output

		screen[lin][col] = elmt
	}

	close(input)
	close(output)

	count := 0

	for _, l := range screen {
		for _, c := range l {
			if c == 0 {
				fmt.Printf(" ")
			} else if c == 1 {
				fmt.Printf("|")
			} else if c == 2 {
				count++
				fmt.Printf("\u2588")
			} else if c == 3 {
				fmt.Printf("_")
			} else if c == 4 {
				fmt.Printf("o")
			}
		}
		fmt.Printf("\n")
	}

	fmt.Printf("The number of blocks when the game exits is %v\n", count)

	//Part 2
	program = copyProgram(baseProgram)
	//Play for free (:
	program[0] = 2

	screen = make([][]int, 40)

	for i := 0; i < len(screen); i++ {
		screen[i] = make([]int, 40)
	}

	input = make(chan int, 100)
	output = make(chan int, 100)

	go func() {
		ret, _ := compute(program, input, output)

		output <- ret
	}()

	score := 0
	paddlePos := pos{0, 0}

	for {
		col := <-output

		if col == -99 {
			//End of game
			break
		}

		lin := <-output
		elmt := <-output

		if col == -1 && lin == 0 {
			score = elmt
		} else {
			screen[lin][col] = elmt
		}

		if elmt == 3 {
			paddlePos.col = col
			paddlePos.lin = lin
		}

		if elmt == 4 {
			if paddlePos.col < col {
				input <- 1
			} else if paddlePos.col > col {
				input <- -1
			} else {
				input <- 0
			}
			time.Sleep(100 * time.Millisecond)
		}

		// Displaying is broken but good enough!
		for i := 0; i < 22; i++ {
			for _, c := range screen[i] {
				if c == 0 {
					fmt.Printf(" ")
				} else if c == 1 {
					fmt.Printf("|")
				} else if c == 2 {
					count++
					fmt.Printf("\u2588")
				} else if c == 3 {
					fmt.Printf("_")
				} else if c == 4 {
					fmt.Printf("o")
				}
			}
			fmt.Printf("\n")
		}

		fmt.Printf("Score : %v\n", score)
	}

	close(input)
	close(output)

	fmt.Printf("The final score is %v\n", score)

}

func compute(program []int, input chan int, output chan int) (int, error) {
	i := 0
	// prevFunc := 0
	relativeBase := 0

	for i < len(program) {
		op := program[i] % 100
		mode1 := (program[i] % 1000) / 100
		mode2 := (program[i] % 10000) / 1000
		mode3 := (program[i] % 100000) / 10000

		// mode
		// 0 = position
		// 1 = immediate
		// 2 = relative

		// Parameters that an instruction writes to will never be in immediate mode
		// Thus op 1, 2 & 3 will never have a mode 3

		if op == 99 {
			return -99, nil
		}

		if op == 1 {
			//Add
			//prevFunc = program[i]

			if i+3 >= len(program) {
				//fmt.Println("Tried to access a value outside of the program!")
				return -1, errors.New("tried to access a value outside of the program")
			}

			val1 := program[i+1]
			val2 := program[i+2]
			res := program[i+3]

			if mode1 == 0 {
				val1 = program[val1]
			} else if mode1 == 2 {
				val1 = program[relativeBase+val1]
			}

			if mode2 == 0 {
				val2 = program[val2]
			} else if mode2 == 2 {
				val2 = program[relativeBase+val2]
			}

			if mode3 == 2 {
				res = res + relativeBase
			}

			program[res] = val1 + val2

			i += 4
		} else if op == 2 {
			//Multiply
			//prevFunc = program[i]

			if i+3 >= len(program) {
				//fmt.Println("Tried to access a value outside of the program!")
				return -1, errors.New("tried to access a value outside of the program")
			}

			val1 := program[i+1]
			val2 := program[i+2]
			res := program[i+3]

			if mode1 == 0 {
				val1 = program[val1]
			} else if mode1 == 2 {
				val1 = program[relativeBase+val1]
			}

			if mode2 == 0 {
				val2 = program[val2]
			} else if mode2 == 2 {
				val2 = program[relativeBase+val2]
			}

			if mode3 == 2 {
				res = res + relativeBase
			}

			program[res] = val1 * val2

			i += 4
		} else if op == 3 {
			//Store input
			// scanner := bufio.NewScanner(os.Stdin)
			// fmt.Print("Enter input: ")
			// scanner.Scan()

			// inVal, err := strconv.Atoi(scanner.Text())

			// if err != nil {
			// 	return -1, errors.New("invalid input")
			// }

			inVal := <-input

			res := program[i+1]

			if mode1 == 0 {
				program[res] = inVal

			} else if mode1 == 2 {
				program[relativeBase+res] = inVal
			}

			i += 2
		} else if op == 4 {
			//Output
			res := program[i+1]
			var out int

			if mode1 == 0 {
				out = program[res]
			} else if mode1 == 2 {
				out = program[res+relativeBase]
			} else {
				out = res
			}

			//fmt.Printf("Output : %d, previous function : %d\n", out, prevFunc)

			output <- out

			i += 2
		} else if op == 5 {
			// jump if true
			//prevFunc = program[i]

			val1 := program[i+1]
			res := program[i+2]

			if mode1 == 0 {
				val1 = program[val1]
			} else if mode1 == 2 {
				val1 = program[val1+relativeBase]
			}

			if mode2 == 0 {
				res = program[res]
			} else if mode2 == 2 {
				res = program[res+relativeBase]
			}

			if val1 != 0 {
				i = res
			} else {
				i += 3
			}
		} else if op == 6 {
			// jump if false
			//prevFunc = program[i]

			val1 := program[i+1]
			res := program[i+2]

			if mode1 == 0 {
				val1 = program[val1]
			} else if mode1 == 2 {
				val1 = program[val1+relativeBase]
			}

			if mode2 == 0 {
				res = program[res]
			} else if mode2 == 2 {
				res = program[res+relativeBase]
			}

			if val1 == 0 {
				i = res
			} else {
				i += 3
			}
		} else if op == 7 {
			// less than
			//prevFunc = program[i]

			val1 := program[i+1]
			val2 := program[i+2]
			res := program[i+3]

			if mode1 == 0 {
				val1 = program[val1]
			} else if mode1 == 2 {
				val1 = program[val1+relativeBase]
			}

			if mode2 == 0 {
				val2 = program[val2]
			} else if mode2 == 2 {
				val2 = program[val2+relativeBase]
			}

			if mode3 == 2 {
				res = res + relativeBase
			}

			if val1 < val2 {
				program[res] = 1
			} else {
				program[res] = 0
			}

			i += 4
		} else if op == 8 {
			// equals
			//prevFunc = program[i]

			val1 := program[i+1]
			val2 := program[i+2]
			res := program[i+3]

			if mode1 == 0 {
				val1 = program[val1]
			} else if mode1 == 2 {
				val1 = program[val1+relativeBase]
			}

			if mode2 == 0 {
				val2 = program[val2]
			} else if mode2 == 2 {
				val2 = program[val2+relativeBase]
			}

			if mode3 == 2 {
				res = res + relativeBase
			}

			if val1 == val2 {
				program[res] = 1
			} else {
				program[res] = 0
			}

			i += 4
		} else if op == 9 {
			//Output
			res := program[i+1]

			if mode1 == 0 {
				relativeBase = relativeBase + program[res]
			} else if mode1 == 1 {
				relativeBase = relativeBase + res
			} else if mode1 == 2 {
				relativeBase = relativeBase + program[relativeBase+res]
			}

			i += 2
		} else {
			return -1, errors.New("unknown op code")
		}
	}

	return program[0], nil
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
