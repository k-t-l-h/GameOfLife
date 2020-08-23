package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"golang.org/x/exp/rand"
	"log"
	"net/http"
	"sync"
	"time"
)

//TODO: добавить конфигурацию
const FSIZE = 400

var matrix [][]int
var index [][][]int

func init() {
	rand.Seed(uint64(time.Now().Unix()))

	matrix = initField()
	matrix = randomizeField(matrix)
	index = generateIndexes()

}

func generateIndexes() [][][]int {

	//создание матрицы индексов
	matrix := make([][][]int, FSIZE)
	for m := range matrix {
		matrix[m] = make([][]int, FSIZE)
		for n := range matrix[m] {
			matrix[m][n] = []int{}
		}
	}

	//обработка краевых точек с 3 соседями
	//левая верхняя
	matrix[0][0] = []int{0, 1,
		1, 1,
		1, 0}
	//правая верхняя
	matrix[0][FSIZE-1] = []int{0, FSIZE - 2,
		1, FSIZE - 2,
		1, FSIZE - 1}
	//левая нижняя
	matrix[FSIZE-1][0] = []int{FSIZE - 2, 0,
		FSIZE - 2, 1,
		FSIZE - 1, 1}
	//правая нижняя
	matrix[FSIZE-1][FSIZE-1] = []int{FSIZE - 1, FSIZE - 2,
		FSIZE - 2, FSIZE - 2,
		FSIZE - 2, FSIZE - 1}

	//обработка краевых рядов с 5 соседями (краевые точки исключены)
	//левый ряд
	for i := 1; i < FSIZE-1; i++ {
		matrix[i][0] = []int{i - 1, 0, i + 1, 0, i + 1, 1, i, 1, i - 1, 1}
	}

	//правый ряд
	for i := 1; i < FSIZE-1; i++ {
		matrix[i][FSIZE-1] = []int{i - 1, FSIZE - 1, i - 1, FSIZE - 2, i, FSIZE - 2, i + 1, FSIZE - 2, i + 1, FSIZE - 1}
	}

	//верхний ряд
	for i := 1; i < FSIZE-1; i++ {
		matrix[0][i] = []int{0, i - 1, 0, i + 1, 1, i - 1, 1, i, 1, i + 1}
	}

	//нижний ряд
	for i := 1; i < FSIZE-1; i++ {
		matrix[FSIZE-1][i] = []int{FSIZE - 1, i - 1, FSIZE - 1, i + 1, FSIZE - 2, i - 1, FSIZE - 2, i, FSIZE - 2, i + 1}
	}

	//обработка внутреннего квадрата, все 9 соседей
	for i := 1; i < FSIZE-1; i++ {
		for j := 1; j < FSIZE-1; j++ {
			left := i - 1
			right := i + 1
			up := j - 1
			down := j + 1
			matrix[i][j] = []int{left, up, i, up, right, up, left, j, right, j, left, down, i, down, right, down}
		}
	}

	return matrix
}

func main() {

	r := mux.NewRouter()
	r.HandleFunc("/game", GetNext).Methods("GET")

	srv := &http.Server{
		Handler:      r,
		Addr:         "127.0.0.1:4000",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Print("Running on :6000")
	log.Print(srv.ListenAndServe())

}

func GetNext(writer http.ResponseWriter, request *http.Request) {

	t := time.Now()
	matrix = nextField(matrix)
	log.Print(time.Since(t))

	writer.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(writer).Encode(matrix[0][0])
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

func randomizeField(matrix [][]int) [][]int {

	for i := range matrix {
		for j := range matrix[i] {
			matrix[i][j] = rand.Intn(2)
		}
	}
	return matrix
}

func nextField(matrix [][]int) [][]int {

	mt := make([][]int, FSIZE)
	for m := range mt {
		mt[m] = make([]int, FSIZE)
	}

	var wg sync.WaitGroup

	for i := range matrix {
		for j := range matrix[i] {

			wg.Add(1)
			go func(i, j int) {
				neighbours := 0

				for n := 0; n < (len(index[i][j]) >> 1); n += 2 {
					neighbours += matrix[(index[i][j][n])][(index[i][j][n+1])]
				}

				if matrix[i][j] == 0 {
					if neighbours == 3 {
						mt[i][j] = 1
					}

				} else {
					if neighbours == 3 || neighbours == 2 {
						mt[i][j] = 1
					}
				}

				wg.Done()
			}(i, j)

		}
	}

	wg.Wait()

	return mt
}
