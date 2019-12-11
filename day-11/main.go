package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
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
	hull := make([][]int, 100)

	for i := range hull {
		hull[i] = make([]int, 100)
	}

	panPaint := make([][]bool, 100)
	for i := range panPaint {
		panPaint[i] = make([]bool, 100)
	}

	initPos := pos{len(hull) / 2, len(hull[0]) / 2}
	initDir := 0

	input := make(chan int, 100)
	output := make(chan int, 100)
	retCode := 0

	go func() {
		retCode, _ = compute(program, input, output)
	}()

	numPainted := 0

	for retCode != 99 {
		input <- hull[initPos.lin][initPos.col]
		col := <-output
		dir := <-output

		hull[initPos.lin][initPos.col] = col

		if !panPaint[initPos.lin][initPos.col] {
			panPaint[initPos.lin][initPos.col] = true
			numPainted++
		}

		if dir == 0 {
			initDir--

			if initDir == -1 {
				initDir = 3
			}
		} else {
			initDir = (initDir + 1) % 4
		}

		if initDir == 0 {
			initPos.lin--
		} else if initDir == 1 {
			initPos.col++
		} else if initDir == 2 {
			initPos.lin++
		} else {
			initPos.col--
		}
	}

	close(input)
	close(output)

	fmt.Printf("The number of panels the robot paints at least one is %v\n", numPainted)

	//Part 2
	resetProgram(program, baseProgram)

	hull = make([][]int, 100)

	for i := range hull {
		hull[i] = make([]int, 100)
	}

	initPos = pos{len(hull) / 2, len(hull[0]) / 2}
	initDir = 0

	//Start on a white panel
	hull[initPos.lin][initPos.col] = 1

	input = make(chan int, 100)
	output = make(chan int, 100)
	retCode = 0

	go func() {
		retCode, _ = compute(program, input, output)
	}()

	for retCode != 99 {
		input <- hull[initPos.lin][initPos.col]
		col := <-output
		dir := <-output

		hull[initPos.lin][initPos.col] = col

		if !panPaint[initPos.lin][initPos.col] {
			panPaint[initPos.lin][initPos.col] = true
			numPainted++
		}

		if dir == 0 {
			initDir--

			if initDir == -1 {
				initDir = 3
			}
		} else {
			initDir = (initDir + 1) % 4
		}

		if initDir == 0 {
			initPos.lin--
		} else if initDir == 1 {
			initPos.col++
		} else if initDir == 2 {
			initPos.lin++
		} else {
			initPos.col--
		}
	}

	close(input)
	close(output)

	for _, l := range hull {
		for _, c := range l {
			if c == 1 {
				fmt.Printf("%v", "\u2588")
			} else {
				fmt.Printf(" ")
			}
		}
		fmt.Printf("\n")
	}

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
			return 99, nil
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
