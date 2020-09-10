package main

import (
	"github.com/k-t-l-h/GameOfLife/protomsg"
	"encoding/json"
	"github.com/golang/protobuf/proto"
	"github.com/gorilla/mux"
	"golang.org/x/exp/rand"
	"log"
	"net/http"
	"time"
)

//TODO: добавить конфигурацию
const FSIZE = 1000

var matrix [][]int32
var index [][][]int32

func init() {
	//TODO: задать зерно через crypto
	rand.Seed(uint64(time.Now().Unix()))

	//инициализация глобальных переменных
	matrix = initField()
	matrix = randomizeField(matrix)
	index = generateIndexes()

}

func generateIndexes() [][][]int32 {
	//матрица векторов индексов клеток-соседей
	//индекс на четном месте определяет строчку, на нечетном столбец

	//создание матрицы индексов
	matrix := make([][][]int32, FSIZE)
	for m := range matrix {
		matrix[m] = make([][]int32, FSIZE)
		for n := range matrix[m] {
			matrix[m][n] = []int32{}
		}
	}

	//обработка угловых клеток (с 3 соседями)
	//левая верхняя
	matrix[0][0] = []int32{0, 1,
		1, 1,
		1, 0}
	//правая верхняя
	matrix[0][FSIZE-1] = []int32{0, FSIZE - 2,
		1, FSIZE - 2,
		1, FSIZE - 1}
	//левая нижняя
	matrix[FSIZE-1][0] = []int32{FSIZE - 2, 0,
		FSIZE - 2, 1,
		FSIZE - 1, 1}
	//правая нижняя
	matrix[FSIZE-1][FSIZE-1] = []int32{FSIZE - 1, FSIZE - 2,
		FSIZE - 2, FSIZE - 2,
		FSIZE - 2, FSIZE - 1}

	//обработка краевых клеток (с 5 соседями, краевые клетки исключены)

	for i := int32(1); i < FSIZE-1; i++ {
		//левый ряд
		matrix[i][0] = []int32{i - 1, 0, i + 1, 0, i + 1, 1, i, 1, i - 1, 1}
		//правый ряд
		matrix[i][FSIZE-1] = []int32{i - 1, FSIZE - 1, i - 1, FSIZE - 2, i, FSIZE - 2, i + 1, FSIZE - 2, i + 1, FSIZE - 1}
		//верхний ряд
		matrix[0][i] = []int32{0, i - 1, 0, i + 1, 1, i - 1, 1, i, 1, i + 1}
		//нижний ряд
		matrix[FSIZE-1][i] = []int32{FSIZE - 1, i - 1, FSIZE - 1, i + 1, FSIZE - 2, i - 1, FSIZE - 2, i, FSIZE - 2, i + 1}
	}

	//обработка внутреннего квадрата (все 8 соседей)
	for i := int32(1); i < FSIZE-1; i++ {
		for j := int32(1); j < FSIZE-1; j++ {
			left := i - 1
			right := i + 1
			up := j - 1
			down := j + 1
			matrix[i][j] = []int32{left, up, i, up, right, up, left, j, right, j, left, down, i, down, right, down}
		}
	}

	return matrix
}

//инициализация глобальной структуры
func initField() [][]int32 {
	matrix := make([][]int32, FSIZE)
	for m := range matrix {
		matrix[m] = make([]int32, FSIZE)
	}
	return matrix
}

//заполнение поля случайными значениями
func randomizeField(matrix [][]int32) [][]int32 {
	for i := range matrix {
		for j := range matrix[i] {
			matrix[i][j] = rand.Int31n(2)
		}
	}
	return matrix
}

//получение нового поля на основе текущего состояния
func nextField(matrix [][]int32) [][]int32 {

	//создание новой матрицы
	mt := make([][]int32, FSIZE)
	for m := range mt {
		mt[m] = make([]int32, FSIZE)
	}

	for i := range matrix {
		for j := range matrix[i] {

			neighbours := int32(0)

			//подсчёт количества живых соседей
			for n := 0; n < (len(index[i][j]) >> 1); n += 2 {
				neighbours += matrix[(index[i][j][n])][(index[i][j][n+1])]
			}

			//определение состояния клетки
			if matrix[i][j] == 0 {
				if neighbours == 3 {
					mt[i][j] = 1
				}

			} else {
				if neighbours == 3 || neighbours == 2 {
					mt[i][j] = 1
				}
			}

		}
	}
	return mt
}

func GetNext(writer http.ResponseWriter, request *http.Request) {

	matrix = nextField(matrix)
	//application/json
	writer.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(writer).Encode(matrix)
	if err != nil {
		return
	}

}

func GetNextProto(writer http.ResponseWriter, request *http.Request) {

	matrix = nextField(matrix)

	array := []*protomsg.FieldLine{}
	for i := range matrix {
		arr := new(protomsg.FieldLine)
		arr.Field = append(matrix[i])
		array = append(array, arr)
	}

	f := &protomsg.Field{Array: array}
	data, err := proto.Marshal(f)
	if err != nil {

	}

	//application/protobuf
	writer.Header().Set("Content-Type", "application/protobuf")
	writer.Write(data)

}

func Resurrect(writer http.ResponseWriter, request *http.Request) {

	matrix = randomizeField(matrix)

	writer.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(writer).Encode(matrix)
	if err != nil {
		return
	}

}

func ResurrectProto(writer http.ResponseWriter, request *http.Request) {

	matrix = randomizeField(matrix)

	array := []*protomsg.FieldLine{}
	for i := range matrix {
		arr := new(protomsg.FieldLine)
		arr.Field = append(matrix[i])
		array = append(array, arr)
	}

	f := &protomsg.Field{Array: array}
	data, err := proto.Marshal(f)
	if err != nil {

	}

	//application/protobuf
	writer.Header().Set("Content-Type", "application/protobuf")
	writer.Write(data)

}

func main() {

	r := mux.NewRouter()
	r.HandleFunc("/game", GetNext).Methods("GET")
	r.HandleFunc("/gameproto", GetNextProto).Methods("GET")
	r.HandleFunc("/new", Resurrect).Methods("GET")
	r.HandleFunc("/newproto", ResurrectProto).Methods("GET")

	srv := &http.Server{
		Handler:      r,
		Addr:         "127.0.0.1:5000",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Print("Running on :6000")
	log.Print(srv.ListenAndServe())

}
