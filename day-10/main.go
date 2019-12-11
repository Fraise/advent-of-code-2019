package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"sort"
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

	max, pos := detectAstAngle(asteroids)

	fmt.Printf("The maximum number of asteroids detected is %v at position %v\n", max, pos)

	posVap := vaporiseAstAngle(asteroids, pos.lin, pos.col, 200)

	fmt.Printf("The 200th vaporised asteroid is as pos %v, thus the answer is %v\n", posVap, posVap.col*100+posVap.lin)
}

func detectAstAngle(asteroids [][]bool) (maxAst int, p pair) {
	for i := 0; i < len(asteroids); i++ {
		for j := 0; j < len(asteroids[0]); j++ {
			if !asteroids[i][j] {
				continue
			}

			amap := map[float64]bool{}

			for k := 0; k < len(asteroids); k++ {
				for l := 0; l < len(asteroids[0]); l++ {
					if !asteroids[k][l] {
						continue
					}

					angle := math.Atan2(float64(k-i), float64(l-j))

					amap[angle] = true
				}
			}

			if len(amap) > maxAst {
				maxAst = len(amap)
				p = pair{i, j}
			}
		}
	}

	return maxAst, p
}

func vaporiseAstAngle(asteroids [][]bool, lin int, col int, numDest int) (coord pair) {
	amap := map[float64][]pair{}
	var start float64

	for i := 0; i < len(asteroids); i++ {
		for j := 0; j < len(asteroids[0]); j++ {
			if !asteroids[i][j] {
				continue
			}

			angle := math.Atan2(float64(i-lin), float64(j-col))

			if i-lin < 0 && j-col == 0 {
				start = angle
			}

			if amap[angle] == nil {
				amap[angle] = []pair{pair{i, j}}
			} else {
				arr := amap[angle]

				arr = append(arr, pair{i, j})

				amap[angle] = arr
			}
		}
	}

	//Sort angles
	var list []float64

	for k := range amap {
		list = append(list, k)
	}

	sort.Float64s(list)

	var startPos int

	for i, v := range list {
		if v == start {
			startPos = i
			break
		}
	}

	count := 0

	for count < numDest {
		pos := list[startPos]
		min := 1000000

		for _, p := range amap[pos] {
			currMin := int(math.Abs(float64(p.col-col)) + math.Abs(float64(p.lin-lin)))
			if currMin < min {
				min = currMin
			}
		}

		for i, p := range amap[pos] {
			currMin := int(math.Abs(float64(p.col-col)) + math.Abs(float64(p.lin-lin)))
			if currMin == min {
				if count == numDest-1 {
					coord = p
				} else {
					copy(amap[pos][i:], amap[pos][i+1:])
				}
			}
		}

		startPos = (startPos + 1) % len(list)
		count++
	}

	return coord
}

func detectAst(asteroids [][]bool) (maxAst int) {
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

	return maxAst
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
