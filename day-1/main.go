package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
)

func main() {
	var moduleList []int
	var sum int = 0

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

	fmt.Println("Sum = " + strconv.Itoa(sum))
}
