package main

import (
	"fmt"
	"math"
	"os"
	"sort"
	"time"
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

func main() {
	inicio = time.Now()

	// Lee el archivo tsp
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

	// Numero de filas de coordenadas
	numRows := len(slice) / 3

	// slice vacio
	grid := make([][]int, numRows)

	//  organiza la estructura del slice de 3 x n
	for i := 0; i < numRows; i++ {
		grid[i] = make([]int, 3)
		// fmt.Println(grid[i])

	}
	//se llena el slice con las cordenadas
	c := 0
	for i := 0; i < numRows; i++ {
		for j := 0; j < 3; j++ {
			grid[i][j] = slice[c]
			c++
		}
	}
	// / print out slices
	// fmt.Println(grid)

	// Tamaño de la lista de coordenadas
	numeroNodos := len(grid)
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
				a1 := grid[i][1]
				a2 := grid[i][2]
				b1 := grid[j][1]
				b2 := grid[j][2]

				distanciaEuclidiana := int(math.Sqrt(float64((a1-b1)*(a1-b1)+(a2-b2)*(a2-b2))) + 0.5)
				// fmt.Print("   A1 ", a1, " a2 ", a2, " b1 ", b1, " b2 ", b2, " distancia ", distanciaEuclidiana)
				matrizAdyacencia[i][j] = distanciaEuclidiana
				matrizAdyacenciaDIST[i][j] = distanciaEuclidiana
			}
		}
	}
	//se crea la varible tour y los mapas de nodos
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
	// posibles nodos a ingresar

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

	fmt.Println(" Tour inicial -> ", tour)
	// se llena el tour y el mapa cubierto
	// v := len(mapaSinCubrir)
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
		//oredena el listado, agregando todos los nodos
		sort.SliceStable(listadoInserciones, func(i, j int) bool {
			return listadoInserciones[i].diferencia < listadoInserciones[j].diferencia
		})

		// agrega a tour y mapas cubiertos y elimina nodo mapa sin cubrir
		que := listadoInserciones[0].a
		donde := listadoInserciones[0].b
		tour = append(tour[:donde+1], tour[donde:]...)
		tour[donde] = que
		mapaCubierto[que] = struct{}{}
		delete(mapaSinCubrir, que)
		// fmt.Println(" Tour ->", tour)
	}
	// cuneta la distancia recorrida al final
	cont := len(tour) - 1
	fo := 0
	for i := 0; i < cont; i++ {
		fo = fo + matrizAdyacenciaDIST[tour[i]][tour[i+1]]
		// fmt.Print("Función Objetivo->", fo)
	}
	// adiciona el ultimo retorno
	aux := tour[len(tour)-1:]
	// fmt.Print(aux[0], ' ', tour[0], ' ', matrizAdyacenciaDIST[aux[0]][tour[0]], len(tour))
	fo = fo + matrizAdyacencia[aux[0]][tour[0]]
	// imprime y cuenta el tiempo
	fmt.Print()
	fmt.Println("\nFunción Objetivo->", fo)
	fmt.Print("\nRuta recorrida ", tour)

	fmt.Println("\nTiempo de inicio -> ", inicio.Second())

	fin = time.Now()
	fmt.Println("\nTiempo de fin -> ", fin.Second())

	diff := fin.Sub(inicio)
	fmt.Println("\nTotal tiempo -> ", diff)

}
