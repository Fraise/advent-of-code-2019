package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {
	f, err := os.Open(os.Args[1])
	defer f.Close()

	if err != nil {
		log.Fatalln(err)
	}

	var planetMap map[string][]string = make(map[string][]string)

	s := bufio.NewScanner(f)

	for s.Scan() {
		planets := strings.Split(s.Text(), ")")

		orbits := planetMap[planets[0]]

		orbits = append(orbits, planets[1])

		planetMap[planets[0]] = orbits
	}

	//Part 1
	orbitsNum := helper(planetMap["COM"], 0, planetMap)

	fmt.Printf("The total number of direct and indirect orbits is : %d\n", orbitsNum)

	//Part 2
	ances := commonAncestor("COM", "YOU", "SAN", planetMap)

	path1 := shortestPath(ances, "YOU", planetMap)
	path2 := shortestPath(ances, "SAN", planetMap)

	fmt.Printf("The minimun of orbital transfers required is : %v\n", path1+path2-2)
}

func helper(orbits []string, num int, planetMap map[string][]string) int {
	if orbits == nil {
		return num
	}

	total := num

	for _, o := range orbits {
		orbs := helper(planetMap[o], num+1, planetMap)
		total += orbs
	}

	return total
}

func commonAncestor(curr string, o1 string, o2 string, planetMap map[string][]string) string {
	if curr == o1 || curr == o2 {
		return curr
	}

	if planetMap[curr] == nil {
		return ""
	}

	var o1found string
	var o2found string

	for _, o := range planetMap[curr] {
		found := commonAncestor(o, o1, o2, planetMap)

		if o1found == "" && found != "" {
			o1found = found
			continue
		}

		if o1found != "" && o2found == "" && found != "" {
			o2found = found
		}
	}

	if o1found != "" && o2found != "" {
		return curr
	}

	if o1found != "" {
		return o1found
	}

	return o2found
}

func shortestPath(curr string, o string, planetMap map[string][]string) int {
	if curr == o {
		return 0
	}

	if planetMap[curr] == nil {
		return -1
	}

	for _, orb := range planetMap[curr] {
		val := shortestPath(orb, o, planetMap)

		if val >= 0 {
			return val + 1
		}
	}

	return -1
}
