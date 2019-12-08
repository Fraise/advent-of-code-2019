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

	vals := []int{0, 1, 2, 3, 4}
	perms := make([][]int, 0)

	perms = permutations(vals)

	inputsChan := make(chan int, 5)
	defer close(inputsChan)

	outputsChan := make(chan int, 5)
	defer close(outputsChan)

	baseProgram := copyProgram(program)

	var max int

	for _, p := range perms {
		for j, v := range p {
			if j == 0 {
				inputsChan <- v
				inputsChan <- 0

				_, err = compute(program, inputsChan, outputsChan)
			} else {
				inputsChan <- v
				inputsChan <- <-outputsChan

				_, err = compute(program, inputsChan, outputsChan)
			}

			//Restore original program
			resetProgram(program, baseProgram)

			if err != nil {
				log.Printf("Error : %v", err)
			}
		}

		curr := <-outputsChan

		if curr > max {
			max = curr
		}
	}

	fmt.Printf("The highest signal that can be sent to the thursters is : %v\n", max)

	vals = []int{9, 7, 8, 5, 6}

	perms = permutations(vals)

	maxOut := 0

	for _, p := range perms {
		ptrs := []int{0, 0, 0, 0, 0}

		count := 0
		amp := 0
		lastOut := 0

		programs := [][]int{
			copyProgram(program),
			copyProgram(program),
			copyProgram(program),
			copyProgram(program),
			copyProgram(program),
		}

		inChan := make(chan int, 5)

		outChan := make(chan int, 5)

		for ptrs[0] != -99 && ptrs[1] != -99 && ptrs[2] != -99 && ptrs[3] != -99 && ptrs[4] != -99 {
			if count == 0 {
				for j, v := range p {

					inChan <- v
					if j == 0 {
						inChan <- 0
					} else {
						inChan <- <-outChan
					}

					ptrs[amp], err = computeFeedback(programs[amp], inChan, outChan, ptrs[amp])

					if err != nil {
						log.Printf("Error : %v", err)
					}

					amp = (amp + 1) % 5
				}
				count++
			} else {
				currOut := <-outChan

				if amp == 0 {
					lastOut = currOut
				}

				inChan <- currOut

				ptrs[amp], err = computeFeedback(programs[amp], inChan, outChan, ptrs[amp])

				if err != nil {
					log.Printf("Error : %v", err)
				}

				amp = (amp + 1) % 5
			}
		}

		if lastOut > maxOut {
			maxOut = lastOut
		}

		close(inChan)
		close(outChan)
	}

	fmt.Printf("The highest signal that can be sent to the thursters through a feedback loop is : %v\n", maxOut)
}

func permutations(arr []int) [][]int {
	var helper func([]int, int)
	res := [][]int{}

	helper = func(arr []int, n int) {
		if n == 1 {
			tmp := make([]int, len(arr))
			copy(tmp, arr)
			res = append(res, tmp)
		} else {
			for i := 0; i < n; i++ {
				helper(arr, n-1)
				if n%2 == 1 {
					tmp := arr[i]
					arr[i] = arr[n-1]
					arr[n-1] = tmp
				} else {
					tmp := arr[0]
					arr[0] = arr[n-1]
					arr[n-1] = tmp
				}
			}
		}
	}
	helper(arr, len(arr))
	return res
}

func compute(program []int, input <-chan int, output chan<- int) (int, error) {
	i := 0
	//prevFunc := 0

	for i < len(program) {
		op := program[i] % 100
		mode1 := (program[i] % 1000) / 100
		mode2 := (program[i] % 10000) / 1000
		//mode3 := (program[i] % 100000) / 10000

		// mode
		// 0 = position
		// 1 = immediate

		// Parameters that an instruction writes to will never be in immediate mode
		// Thus op 1, 2 & 3 will never have a mode 3

		if op == 99 {
			//fmt.Println("op 99 : Program halted.")
			break
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

			switch {
			case mode1 == 0 && mode2 == 0:
				program[res] = program[val1] + program[val2]
			case mode1 == 1 && mode2 == 0:
				program[res] = val1 + program[val2]
			case mode1 == 0 && mode2 == 1:
				program[res] = program[val1] + val2
			case mode1 == 1 && mode2 == 1:
				program[res] = val1 + val2
			default:
				return -1, errors.New("error processing op code : " + strconv.Itoa(program[i]))
			}

			i += 4
		} else if op == 2 {
			//Multiply
			//prevFunc = program[i]

			val1 := program[i+1]
			val2 := program[i+2]
			res := program[i+3]

			switch {
			case mode1 == 0 && mode2 == 0:
				program[res] = program[val1] * program[val2]
			case mode1 == 1 && mode2 == 0:
				program[res] = val1 * program[val2]
			case mode1 == 0 && mode2 == 1:
				program[res] = program[val1] * val2
			case mode1 == 1 && mode2 == 1:
				program[res] = val1 * val2
			default:
				return -1, errors.New("error processing op code : " + strconv.Itoa(program[i]))
			}

			i += 4
		} else if op == 3 {
			// //Store input
			// scanner := bufio.NewScanner(os.Stdin)
			// fmt.Print("Enter input: ")
			// scanner.Scan()

			// inVal, err := strconv.Atoi(scanner.Text())

			// if err != nil {
			// 	return -1, errors.New("invalid input")
			// }

			// res := program[i+1]

			// program[res] = inVal

			// i += 2

			var inVal int

			inVal = <-input

			res := program[i+1]
			program[res] = inVal
			i += 2
		} else if op == 4 {
			//Output
			res := program[i+1]
			var out int

			if mode1 == 0 {
				out = program[res]
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
			}

			if mode2 == 0 {
				res = program[res]
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
			}

			if mode2 == 0 {
				res = program[res]
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
			}

			if mode2 == 0 {
				val2 = program[val2]
			}

			if val1 < val2 {
				program[res] = 1
			} else {
				program[res] = 0
			}

			i += 4
		} else if op == 8 {
			// less than
			//prevFunc = program[i]

			val1 := program[i+1]
			val2 := program[i+2]
			res := program[i+3]

			if mode1 == 0 {
				val1 = program[val1]
			}

			if mode2 == 0 {
				val2 = program[val2]
			}

			if val1 == val2 {
				program[res] = 1
			} else {
				program[res] = 0
			}

			i += 4
		} else {
			return -1, errors.New("unknown op code")
		}
	}

	return program[0], nil
}

func computeFeedback(program []int, input <-chan int, output chan<- int, ptr int) (int, error) {
	i := ptr
	//prevFunc := 0

	for i < len(program) {
		op := program[i] % 100
		mode1 := (program[i] % 1000) / 100
		mode2 := (program[i] % 10000) / 1000
		//mode3 := (program[i] % 100000) / 10000

		// mode
		// 0 = position
		// 1 = immediate

		// Parameters that an instruction writes to will never be in immediate mode
		// Thus op 1, 2 & 3 will never have a mode 3

		if op == 99 {
			//fmt.Println("op 99 : Program halted.")
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

			switch {
			case mode1 == 0 && mode2 == 0:
				program[res] = program[val1] + program[val2]
			case mode1 == 1 && mode2 == 0:
				program[res] = val1 + program[val2]
			case mode1 == 0 && mode2 == 1:
				program[res] = program[val1] + val2
			case mode1 == 1 && mode2 == 1:
				program[res] = val1 + val2
			default:
				return -1, errors.New("error processing op code : " + strconv.Itoa(program[i]))
			}

			i += 4
		} else if op == 2 {
			//Multiply
			//prevFunc = program[i]

			val1 := program[i+1]
			val2 := program[i+2]
			res := program[i+3]

			switch {
			case mode1 == 0 && mode2 == 0:
				program[res] = program[val1] * program[val2]
			case mode1 == 1 && mode2 == 0:
				program[res] = val1 * program[val2]
			case mode1 == 0 && mode2 == 1:
				program[res] = program[val1] * val2
			case mode1 == 1 && mode2 == 1:
				program[res] = val1 * val2
			default:
				return -1, errors.New("error processing op code : " + strconv.Itoa(program[i]))
			}

			i += 4
		} else if op == 3 {
			// //Store input
			// scanner := bufio.NewScanner(os.Stdin)
			// fmt.Print("Enter input: ")
			// scanner.Scan()

			// inVal, err := strconv.Atoi(scanner.Text())

			// if err != nil {
			// 	return -1, errors.New("invalid input")
			// }

			// res := program[i+1]

			// program[res] = inVal

			// i += 2

			var inVal int

			inVal = <-input

			res := program[i+1]
			program[res] = inVal
			i += 2
		} else if op == 4 {
			//Output
			res := program[i+1]
			var out int

			if mode1 == 0 {
				out = program[res]
			} else {
				out = res
			}

			//fmt.Printf("Output : %d, previous function : %d\n", out, prevFunc)
			output <- out

			i += 2
			return i, nil
		} else if op == 5 {
			// jump if true
			//prevFunc = program[i]

			val1 := program[i+1]
			res := program[i+2]

			if mode1 == 0 {
				val1 = program[val1]
			}

			if mode2 == 0 {
				if res >= len(program) {
					return -1, errors.New("tried to access a value outside of the program")
				}

				res = program[res]
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
			}

			if mode2 == 0 {
				res = program[res]
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
			}

			if mode2 == 0 {
				val2 = program[val2]
			}

			if val1 < val2 {
				program[res] = 1
			} else {
				program[res] = 0
			}

			i += 4
		} else if op == 8 {
			// less than
			//prevFunc = program[i]

			val1 := program[i+1]
			val2 := program[i+2]
			res := program[i+3]

			if mode1 == 0 {
				val1 = program[val1]
			}

			if mode2 == 0 {
				val2 = program[val2]
			}

			if val1 == val2 {
				program[res] = 1
			} else {
				program[res] = 0
			}

			i += 4
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
