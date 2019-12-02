package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
)

func main() {

	sum1 := part1(os.Args[1])
	sum2 := part2(os.Args[1])

	fmt.Println("Sum1 = " + strconv.Itoa(sum1))
	fmt.Println("Sum2 = " + strconv.Itoa(sum2))
}

func part1(inputFile string) (sum int) {
	var moduleList []int

	f, err := os.Open(os.Args[1])
	defer f.Close()

	if err != nil {
		log.Fatalln(err)
	}

	s := bufio.NewScanner(f)
	for s.Scan() {
		m, err := strconv.Atoi(s.Text())

		if err != nil {
			log.Fatalln(err)
		}

		moduleList = append(moduleList, m)
	}

	for _, m := range moduleList {
		sum += m/3 - 2
	}

	return sum
}

func part2(inputFile string) (sum int) {
	var moduleList []int

	f, err := os.Open(inputFile)
	defer f.Close()

	if err != nil {
		log.Fatalln(err)
	}

	s := bufio.NewScanner(f)
	for s.Scan() {
		m, err := strconv.Atoi(s.Text())

		if err != nil {
			log.Fatalln(err)
		}

		moduleList = append(moduleList, m)
	}

	for _, m := range moduleList {
		modFuel := m/3 - 2
		sum += modFuel

		fuel := modFuel/3 - 2

		for fuel > 0 {
			sum += fuel
			fuel = fuel/3 - 2
		}
	}

	return sum
}
