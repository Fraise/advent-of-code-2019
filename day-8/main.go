package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
)

func main() {
	f, err := os.Open(os.Args[1])
	defer f.Close()

	if err != nil {
		log.Fatalln(err)
	}

	s := bufio.NewScanner(f)

	s.Scan()

	rawImage := s.Text()
	height := 6
	length := 25
	layerNum := len(rawImage) / (height * length)

	fmt.Printf("len : %v\n", len(rawImage))

	image := make([][][]int, layerNum)

	for i := range image {
		image[i] = make([][]int, height)
		for j := range image[i] {
			image[i][j] = make([]int, length)
		}
	}

	for i := 0; i < layerNum; i++ {
		for j := 0; j < height; j++ {
			for k := 0; k < length; k++ {
				image[i][j][k], err = strconv.Atoi(string(rawImage[i*height*length+j*length+k]))

				if err != nil {
					fmt.Printf("Error reading raw image")
				}
			}
		}
	}

	num := 0
	min0 := 1000000

	for i := 0; i < layerNum; i++ {
		num0 := 0
		num1 := 0
		num2 := 0
		for j := 0; j < height; j++ {
			for k := 0; k < length; k++ {
				if image[i][j][k] == 0 {
					num0++
				} else if image[i][j][k] == 1 {
					num1++
				} else if image[i][j][k] == 2 {
					num2++
				}
			}
		}

		if num0 < min0 {
			min0 = num0
			num = num1 * num2
		}
	}

	fmt.Printf("The number of 1 * 2 on the layer containing the least 0 is %v\n", num)

	render := make([][]int, height)

	for i := range render {
		render[i] = make([]int, length)
	}

	//Init the rendered image
	for j := 0; j < height; j++ {
		for k := 0; k < length; k++ {
			render[j][k] = -1
		}
	}

	//Render the image
	for i := layerNum - 1; i >= 0; i-- {
		for j := 0; j < height; j++ {
			for k := 0; k < length; k++ {
				if render[j][k] < 0 {
					render[j][k] = image[i][j][k]
				} else if image[i][j][k] < 2 {
					render[j][k] = image[i][j][k]
				}
			}
		}
	}

	fmt.Printf("The decoded image is the following : \n")

	//Display the rendered image
	for j := 0; j < height; j++ {
		for k := 0; k < length; k++ {
			if render[j][k] == 1 {
				fmt.Printf("%v", "\u2588")
			} else {
				fmt.Printf("%v", " ")
			}
		}
		fmt.Printf("\n")
	}
}
