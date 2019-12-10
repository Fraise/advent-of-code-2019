package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
)

type pair struct {
	lin int
	col int
}

func main() {
	f, err := os.Open(os.Args[1])
	defer f.Close()

	if err != nil {
		log.Fatalln(err)
	}

	s := bufio.NewScanner(f)
	var asteroids [][]bool
	num := 0

	for s.Scan() {
		line := s.Text()
		asteroids = append(asteroids, []bool{})

		for _, c := range line {
			if c == rune('.') {
				asteroids[num] = append(asteroids[num], false)
			} else {
				asteroids[num] = append(asteroids[num], true)
			}
		}

		num++
	}

	maxAst := 0

	for i := 0; i < len(asteroids); i++ {
		for j := 0; j < len(asteroids[0]); j++ {
			if !asteroids[i][j] {
				continue
			}

			curr := 0
			c := makeCopy(asteroids)
			v := makeEmpty(asteroids)

			queue := make(chan pair, 1000)
			queue <- pair{i, j}
			empty := false

			v[i][j] = true
			c[i][j] = false

			count := 0

			for !empty {
				count++
				select {
				case p := <-queue:
					//Add the adjacent asteroids
					if p.lin-1 >= 0 {
						if !v[p.lin-1][p.col] {
							v[p.lin-1][p.col] = true
							queue <- pair{p.lin - 1, p.col}
						}

						if p.col-1 >= 0 {
							if !v[p.lin-1][p.col-1] {
								v[p.lin-1][p.col-1] = true
								queue <- pair{p.lin - 1, p.col - 1}
							}
						}
						if p.col+1 < len(asteroids[0]) {
							if !v[p.lin-1][p.col+1] {
								v[p.lin-1][p.col+1] = true
								queue <- pair{p.lin - 1, p.col + 1}
							}
						}
					}
					if p.lin+1 < len(asteroids) {
						if !v[p.lin+1][p.col] {
							v[p.lin+1][p.col] = true
							queue <- pair{p.lin + 1, p.col}
						}

						if p.col-1 >= 0 {
							if !v[p.lin+1][p.col-1] {
								v[p.lin+1][p.col-1] = true
								queue <- pair{p.lin + 1, p.col - 1}
							}
						}
						if p.col+1 < len(asteroids[0]) {
							if !v[p.lin+1][p.col+1] {
								v[p.lin+1][p.col+1] = true
								queue <- pair{p.lin + 1, p.col + 1}
							}
						}
					}
					if p.col-1 >= 0 {
						if !v[p.lin][p.col-1] {
							v[p.lin][p.col-1] = true
							queue <- pair{p.lin, p.col - 1}
						}
					}
					if p.col+1 < len(asteroids[0]) {
						if !v[p.lin][p.col+1] {
							v[p.lin][p.col+1] = true
							queue <- pair{p.lin, p.col + 1}
						}
					}

					//Add a seen asteroid
					if c[p.lin][p.col] {
						curr++

						ldif := p.lin - i
						cdif := p.col - j

						pos := pair{p.lin, p.col}
						seenPos := pair{p.lin, p.col}

						//Remove out of sight asteroids
						for pos.lin < len(asteroids) && pos.lin >= 0 && pos.col < len(asteroids[0]) && pos.col >= 0 {
							c[pos.lin][pos.col] = false
							if ldif == 0 {
								if cdif > 0 {
									pos.col++
								} else {
									pos.col--
								}
							} else if cdif == 0 {
								if ldif > 0 {
									pos.lin++
								} else {
									pos.lin--
								}
							} else if int(math.Abs(float64(ldif))) == int(math.Abs(float64(cdif))) {
								if ldif > 0 {
									pos.lin++
								} else {
									pos.lin--
								}

								if cdif > 0 {
									pos.col++
								} else {
									pos.col--
								}
							} else {
								// Not quick but definitly dirty trick...
								if ldif%10 == 0 && cdif%10 == 0 {
									pos.lin += ldif / 10
									pos.col += cdif / 10
								} else if ldif%9 == 0 && cdif%9 == 0 {
									pos.lin += ldif / 9
									pos.col += cdif / 9
								} else if ldif%8 == 0 && cdif%8 == 0 {
									pos.lin += ldif / 8
									pos.col += cdif / 8
								} else if ldif%7 == 0 && cdif%7 == 0 {
									pos.lin += ldif / 7
									pos.col += cdif / 7
								} else if ldif%6 == 0 && cdif%6 == 0 {
									pos.lin += ldif / 6
									pos.col += cdif / 6
								} else if ldif%5 == 0 && cdif%5 == 0 {
									pos.lin += ldif / 5
									pos.col += cdif / 5
								} else if ldif%4 == 0 && cdif%4 == 0 {
									pos.lin += ldif / 4
									pos.col += cdif / 4
								} else if ldif%3 == 0 && cdif%3 == 0 {
									pos.lin += ldif / 3
									pos.col += cdif / 3
								} else if ldif%2 == 0 && cdif%2 == 0 {
									pos.lin += ldif / 2
									pos.col += cdif / 2
								} else {
									pos.lin += ldif
									pos.col += cdif
								}
							}
						}

						c[seenPos.lin][seenPos.col] = true
					}
				default:
					empty = true
				}
			}

			if curr > maxAst {
				maxAst = curr
			}

			close(queue)
		}
	}

	fmt.Printf("The maximum number of asteroids detected is : %v\n", maxAst)
}

func makeCopy(asteroids [][]bool) [][]bool {
	c := make([][]bool, len(asteroids))

	for i := range c {
		c[i] = make([]bool, len(asteroids[i]))
		copy(c[i], asteroids[i])
	}

	return c
}

func makeEmpty(asteroids [][]bool) [][]bool {
	c := make([][]bool, len(asteroids))

	for i := range c {
		c[i] = make([]bool, len(asteroids[i]))
	}

	return c
}
