package main

import (
	"fmt"
	"math"
	"os"
	"sort"
	"time"
	// "math"
)

var inicio time.Time
var fin time.Time

type posibleNodo struct {
	id int
	x  int
}

type posibleInsercion struct {
	a          int
	b          int
	diferencia int
}

func leerArchivo(c1 chan [][]int) {
	var coordenadas [][]int
	// Se lee el archivo y se guarda en un slice
	nombre_archivo := "eil51.tsp"
	f, err := os.Open(nombre_archivo)
	if err != nil {
		fmt.Println(err)
		return
	}
	//Slice que recibe los datos
	var slice []int
	for {
		var a int
		_, err := fmt.Fscan(f, &a)

		if err != nil {
			break
		}
		slice = append(slice, a)
		defer f.Close()

	}
	// fmt.Println("slice", slice)
	canal := make(chan [][]int)

	// go organizarArchivo(slice, canal)
	if len(slice)%2 == 0 {
		go organizarArchivo(slice[len(slice)/2:], canal)
		go organizarArchivo(slice[:len(slice)/2], canal)
	} else {
		go organizarArchivo(slice[int(len(slice)/2)-1:], canal)
		go organizarArchivo(slice[:int(len(slice)/2)], canal)
	}
	coordenadasA := <-canal
	coordenadasB := <-canal

	for i := 0; i < len(coordenadasA); i++ {
		coordenadas = append(coordenadas, coordenadasA[i])
	}
	for i := 0; i < len(coordenadasB); i++ {
		coordenadas = append(coordenadas, coordenadasB[i])
	}
	// fmt.Println(" ")
	// fmt.Println(" Coordenadas ", coordenadas)
	c1 <- coordenadas
}

// organiza el archivo
func organizarArchivo(slice []int, canal chan [][]int) {
	// fmt.Println(" organizar archivo ", slice)

	numRows := len(slice) / 3

	grid := make([][]int, numRows)

	for i := 0; i < numRows; i++ {
		grid[i] = make([]int, 3)
	}

	c := 0
	for i := 0; i < numRows; i++ {
		for j := 0; j < 3; j++ {
			grid[i][j] = slice[c]
			c++
		}
	}

	// fmt.Println(" resultado ", grid)
	canal <- grid

}

func matrices(coordenadas [][]int, canalMatriz chan [][]int) {

	numeroNodos := len(coordenadas)
	// fmt.Print(numeroNodos)

	// slice vacio
	matrizAdyacencia := make([][]int, numeroNodos)
	matrizAdyacenciaDIST := make([][]int, numeroNodos)
	//  organiza la estructura del slice de 3 x n
	for i := 0; i < numeroNodos; i++ {
		matrizAdyacencia[i] = make([]int, numeroNodos)
		matrizAdyacenciaDIST[i] = make([]int, numeroNodos)
	}

	for i := 0; i < numeroNodos; i++ {
		for j := 0; j < numeroNodos; j++ {
			if i != j {
				a1 := coordenadas[i][1]
				a2 := coordenadas[i][2]
				b1 := coordenadas[j][1]
				b2 := coordenadas[j][2]

				distanciaEuclidiana := int(math.Sqrt(float64((a1-b1)*(a1-b1)+(a2-b2)*(a2-b2))) + 0.5)
				// fmt.Print("   A1 ", a1, " a2 ", a2, " b1 ", b1, " b2 ", b2, " distancia ", distanciaEuclidiana)
				matrizAdyacencia[i][j] = distanciaEuclidiana
				matrizAdyacenciaDIST[i][j] = distanciaEuclidiana
			}
		}
	}
	canalMatriz <- matrizAdyacencia
	canalMatriz <- matrizAdyacenciaDIST
}

// proceso de inserccion
func proceso(coordenadas [][]int, matrizAdyacencia [][]int, canalTourInicial chan []int, canalTourFinal chan []int) {

	numeroNodos := len(coordenadas)
	var tour []int

	mapaCubierto := make(map[int]struct{})
	mapaSinCubrir := make(map[int]struct{})
	// fmt.Print(numeroNodos)

	for i := 0; i < numeroNodos; i++ {
		mapaSinCubrir[i] = struct{}{}
	}

	nodoInicial := 0
	mapaCubierto[nodoInicial] = struct{}{}
	delete(mapaSinCubrir, nodoInicial)
	tour = append(tour, nodoInicial)

	var posiblesVerticesEnvolvente []posibleNodo

	for i := 0; i < numeroNodos; i++ {
		if matrizAdyacencia[nodoInicial][i] != 0 && i != nodoInicial {

			posiblesVerticesEnvolvente = append(posiblesVerticesEnvolvente, posibleNodo{id: i, x: matrizAdyacencia[nodoInicial][i]})
		}
	}

	sort.SliceStable(posiblesVerticesEnvolvente, func(i, j int) bool {
		return posiblesVerticesEnvolvente[i].x > posiblesVerticesEnvolvente[j].x
	})

	tour = append(tour, posiblesVerticesEnvolvente[0].id)
	mapaCubierto[posiblesVerticesEnvolvente[0].id] = struct{}{}
	delete(mapaSinCubrir, posiblesVerticesEnvolvente[0].id)
	tour = append(tour, posiblesVerticesEnvolvente[1].id)
	mapaCubierto[posiblesVerticesEnvolvente[1].id] = struct{}{}
	delete(mapaSinCubrir, posiblesVerticesEnvolvente[1].id)

	canalTourInicial <- tour

	for len(mapaSinCubrir) > 0 {

		var listadoInserciones []posibleInsercion

		for i := 0; i < len(tour)-1; i++ {
			costoOriginalArista := matrizAdyacencia[tour[i]][tour[i+1]]

			for k := range mapaSinCubrir {
				costoInsercionK := 0
				costoInsercionK = costoInsercionK + matrizAdyacencia[tour[i]][k]
				costoInsercionK = costoInsercionK + matrizAdyacencia[k][tour[i+1]]
				diferenciaInsercion := costoInsercionK - costoOriginalArista

				listadoInserciones = append(listadoInserciones, posibleInsercion{a: k, b: i + 1, diferencia: diferenciaInsercion})
			}
		}

		aux := tour[len(tour)-1:]
		costoOriginalArista := matrizAdyacencia[aux[0]][tour[0]]

		for k := range mapaSinCubrir {
			costoInsercionK := 0
			costoInsercionK = costoInsercionK + matrizAdyacencia[aux[0]][k]
			costoInsercionK = costoInsercionK + matrizAdyacencia[k][tour[0]]
			diferenciaInsercion := costoInsercionK - costoOriginalArista

			listadoInserciones = append(listadoInserciones, posibleInsercion{a: k, b: 0, diferencia: diferenciaInsercion})
		}

		sort.SliceStable(listadoInserciones, func(i, j int) bool {
			return listadoInserciones[i].diferencia < listadoInserciones[j].diferencia
		})

		que := listadoInserciones[0].a
		donde := listadoInserciones[0].b
		tour = append(tour[:donde+1], tour[donde:]...)
		tour[donde] = que
		mapaCubierto[que] = struct{}{}
		delete(mapaSinCubrir, que)
	}
	canalTourFinal <- tour
}

// cuenta la distancia recorrida final
func distanciaFinal(tour []int, matrizAdyacenciaDIST [][]int, matrizAdyacencia [][]int, d chan int) {
	cont := len(tour) - 1
	fo := 0
	for i := 0; i < cont; i++ {
		fo = fo + matrizAdyacenciaDIST[tour[i]][tour[i+1]]
		// fmt.Print("Función Objetivo->", fo)
	}

	aux := tour[len(tour)-1:]
	// fmt.Print(aux[0], ' ', tour[0], ' ', matrizAdyacenciaDIST[aux[0]][tour[0]], len(tour))
	fo = fo + matrizAdyacencia[aux[0]][tour[0]]

	d <- fo
}

func main() {
	inicio = time.Now()

	// Numero de filas de coordenadas
	// Canal uno lectura del archivo funciones leerArchivo, organizarArchivo
	c1 := make(chan [][]int)
	go leerArchivo(c1)
	coordenadasArchivo := <-c1

	// Creacion y llenado de las matrices de adyacencia
	canalMatriz := make(chan [][]int)
	go matrices(coordenadasArchivo, canalMatriz)
	matrizAdyacencia, matrizAdyacenciaDIST := <-canalMatriz, <-canalMatriz

	// Inicializacion de los mapas y el tour
	canalTourInicial := make(chan []int)
	canalTourFinal := make(chan []int)

	// proceso general
	go proceso(coordenadasArchivo, matrizAdyacencia, canalTourInicial, canalTourFinal)
	tourInicial := <-canalTourInicial

	fmt.Println("Tour inicial ", tourInicial)

	tourFinal := <-canalTourFinal

	d := make(chan int)
	go distanciaFinal(tourFinal, matrizAdyacenciaDIST, matrizAdyacencia, d)
	distanciaRecorrida := <-d
	fmt.Println("\nDistancia recorrida -> ", distanciaRecorrida)
	fmt.Println("\nTour final ", tourFinal)

	fmt.Println("\nTiempo de inicio -> ", inicio.Second())
	fin = time.Now()
	fmt.Println("\nTiempo de fin ->", fin.Second())
	diff := fin.Sub(inicio)

	var fin string
	fmt.Scan(&fin)

	fmt.Println("\nTotal tiempo ->", diff)
}
