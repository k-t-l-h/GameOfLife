package main

import (
	"golang.org/x/exp/rand"
	"log"
	"time"
)

func main() {

	rand.Seed(uint64(time.Now().Unix()))
	m := initField()
	n := randomizeField(m)
	log.Print("\n")
	nb := getNeighbours(n)
	log.Print("\n")
	nextField(n, nb)

}

//init the game field
func initField() [][]int {
	matrix := make([][]int, 5)
	for m := range matrix {
		matrix[m] = make([]int, 5)
	}
	return matrix
}

func randomizeField(matrix [][]int) [][]int{

	for i := range matrix {
		for j := range matrix[i] {
			matrix[i][j] = rand.Intn(2)
		}
		log.Print(matrix[i])
	}
	return matrix
}

func getNeighbours(matrix [][]int) [][]int {

	neighbours := make([][]int, 5)
	for m := range neighbours {
		neighbours[m] = make([]int, 5)
	}

	for i := range matrix {
		for j := range matrix[i] {
			neighbours[i][j] += matrix[(5 + i-1) % 5][(5 + j-1) % 5]
			neighbours[i][j] += matrix[(5 + i-1) % 5][j]
			neighbours[i][j] += matrix[(5 + i-1) % 5][(5 + j+1) % 5]
			neighbours[i][j] += matrix[i][(5 + j-1) % 5]
			neighbours[i][j] += matrix[i][(5 + j+1) % 5]
			neighbours[i][j] += matrix[(5 + i+1) % 5][(5 + j-1) % 5]
			neighbours[i][j] += matrix[(5 + i+1) % 5][j]
			neighbours[i][j] += matrix[(5 + i+1) % 5][(5 + j+1) % 5]
		}
		log.Print(neighbours[i])
	}

	return neighbours
}

func nextField(matrix, neighbours[][]int) {

	mt := make([][]int, 5)
	for m := range mt {
		mt[m] = make([]int, 5)
	}


	for i := range matrix {
		for j := range matrix[i] {

			if matrix[i][j] == 0 {
				if neighbours[i][j] == 3 {
					mt[i][j] = 1
				}

			} else {
				if neighbours[i][j] == 3 || neighbours[i][j] == 2 {
					mt[i][j] = 1
				}
			}
		}

		log.Print(mt[i])
	}


}