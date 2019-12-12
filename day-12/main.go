package main

import (
	"fmt"
	"math"
)

type moon struct {
	xp int
	yp int
	zp int

	xv int
	yv int
	zv int
}

type pair struct {
	p int
	v int
}

func main() {
	//input
	moons := []moon{
		moon{14, 4, 5, 0, 0, 0},
		moon{12, 10, 8, 0, 0, 0},
		moon{1, 7, -10, 0, 0, 0},
		moon{16, -5, 3, 0, 0, 0},
	}

	for i := 0; i < 1000; i++ {
		for j := 0; j < len(moons)-1; j++ {
			for k := j + 1; k < len(moons); k++ {
				moons[j], moons[k] = updateVelocity(moons[j], moons[k])
			}
		}

		for l := 0; l < len(moons); l++ {
			moons[l] = updatePosition(moons[l])
		}
	}

	energy := getTotEnergy(moons)

	fmt.Printf("The total energy after 1000 steps is %v\n", energy)

	moonsSim := []moon{
		moon{14, 4, 5, 0, 0, 0},
		moon{12, 10, 8, 0, 0, 0},
		moon{1, 7, -10, 0, 0, 0},
		moon{16, -5, 3, 0, 0, 0},
	}

	moonsSimInit := []moon{
		moon{14, 4, 5, 0, 0, 0},
		moon{12, 10, 8, 0, 0, 0},
		moon{1, 7, -10, 0, 0, 0},
		moon{16, -5, 3, 0, 0, 0},
	}

	cycles := []int{0, 0, 0}

	count := 0

	for {
		for j := 0; j < len(moonsSim)-1; j++ {
			for k := j + 1; k < len(moonsSim); k++ {
				moonsSim[j], moonsSim[k] = updateVelocity(moonsSim[j], moonsSim[k])
			}
		}

		for l := 0; l < len(moons); l++ {
			moonsSim[l] = updatePosition(moonsSim[l])
		}
		count++

		if cycles[0] == 0 {
			xLooped := true

			for i := 0; i < len(moonsSim); i++ {
				if moonsSim[i].xv == 0 {
					xLooped = xLooped && (moonsSim[i].xp == moonsSimInit[i].xp)
				} else {
					xLooped = false
					break
				}
			}

			if xLooped {
				cycles[0] = count
			}
		}

		if cycles[1] == 0 {
			yLooped := true

			for i := 0; i < len(moonsSim); i++ {
				if moonsSim[i].yv == 0 {
					yLooped = yLooped && (moonsSim[i].yp == moonsSimInit[i].yp)
				} else {
					yLooped = false
					break
				}
			}

			if yLooped {
				cycles[1] = count
			}
		}

		if cycles[2] == 0 {
			zLooped := true

			for i := 0; i < len(moonsSim); i++ {
				if moonsSim[i].zv == 0 {
					zLooped = zLooped && (moonsSim[i].zp == moonsSimInit[i].zp)
				} else {
					zLooped = false
					break
				}
			}

			if zLooped {
				cycles[2] = count
			}
		}

		if cycles[0] > 0 && cycles[1] > 0 && cycles[2] > 0 {
			break
		}
	}

	fmt.Printf("We need %v steps to first match a previous state!\n", LCM(cycles[0], cycles[1], cycles[2]))

}

func updateVelocity(m1 moon, m2 moon) (moon, moon) {
	if m1.xp > m2.xp {
		m1.xv--
		m2.xv++
	} else if m1.xp < m2.xp {
		m1.xv++
		m2.xv--
	}

	if m1.yp > m2.yp {
		m1.yv--
		m2.yv++
	} else if m1.yp < m2.yp {
		m1.yv++
		m2.yv--
	}

	if m1.zp > m2.zp {
		m1.zv--
		m2.zv++
	} else if m1.zp < m2.zp {
		m1.zv++
		m2.zv--
	}

	return m1, m2
}

func updatePosition(m moon) moon {
	m.xp += m.xv
	m.yp += m.yv
	m.zp += m.zv

	return m
}

func getTotEnergy(moons []moon) (tot int) {
	for _, m := range moons {
		pot := int(math.Abs(float64(m.xp))) + int(math.Abs(float64(m.yp))) + int(math.Abs(float64(m.zp)))
		kin := int(math.Abs(float64(m.xv))) + int(math.Abs(float64(m.yv))) + int(math.Abs(float64(m.zv)))

		tot += pot * kin
	}

	return tot
}

// From https://play.golang.org/p/SmzvkDjYlb :

// greatest common divisor (GCD) via Euclidean algorithm
func GCD(a, b int) int {
	for b != 0 {
		t := b
		b = a % b
		a = t
	}
	return a
}

// find Least Common Multiple (LCM) via GCD
func LCM(a, b int, integers ...int) int {
	result := a * b / GCD(a, b)

	for i := 0; i < len(integers); i++ {
		result = LCM(result, integers[i])
	}

	return result
}
