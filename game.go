package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"golang.org/x/exp/rand"
	"log"
	"net/http"
	"time"
)

const FSIZE  = 8

var matrix [][]int

func init() {
	rand.Seed(uint64(time.Now().Unix()))
	matrix = initField()
	matrix = randomizeField(matrix)
}

func main() {

	r := mux.NewRouter()
	r.HandleFunc("/game", GetNext).Methods("GET")

	srv := &http.Server{
		Handler: r,
		Addr:    ":6000",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Print("Running on :6000")
	log.Print(srv.ListenAndServe())

}

func GetNext(writer http.ResponseWriter, request *http.Request) {
	nb := getNeighbours(matrix)
	matrix = nextField(matrix, nb)


	writer.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(writer).Encode(matrix)
	if err != nil {
		return
	}
	log.Print(err)
}



//init the game field
func initField() [][]int {
	matrix := make([][]int, FSIZE)
	for m := range matrix {
		matrix[m] = make([]int, FSIZE)
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

	neighbours := make([][]int, FSIZE)
	for m := range neighbours {
		neighbours[m] = make([]int, FSIZE)
	}

	for i := range matrix {
		for j := range matrix[i] {
			neighbours[i][j] += matrix[(FSIZE + i-1) % FSIZE][(FSIZE + j-1) % FSIZE]
			neighbours[i][j] += matrix[(FSIZE + i-1) % FSIZE][j]
			neighbours[i][j] += matrix[(FSIZE + i-1) % FSIZE][(FSIZE + j+1) % FSIZE]
			neighbours[i][j] += matrix[i][(FSIZE + j-1) % FSIZE]
			neighbours[i][j] += matrix[i][(FSIZE + j+1) % FSIZE]
			neighbours[i][j] += matrix[(FSIZE + i+1) % FSIZE][(FSIZE + j-1) % FSIZE]
			neighbours[i][j] += matrix[(FSIZE + i+1) % FSIZE][j]
			neighbours[i][j] += matrix[(FSIZE + i+1) % FSIZE][(FSIZE + j+1) % FSIZE]
		}
		log.Print(neighbours[i])
	}

	return neighbours
}

func nextField(matrix, neighbours [][]int) [][]int {

	mt := make([][]int, FSIZE)
	for m := range mt {
		mt[m] = make([]int, FSIZE)
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

	return mt
}