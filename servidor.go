package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
)

const GRID = 4 // 4x4 = 16 quadrantes

var matrizFinal = make(map[int][]int)
var mu sync.Mutex

func main() {
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Erro ao iniciar servidor:", err)
		return
	}
	defer ln.Close()

	fmt.Println("Servidor aguardando conexões...")

	for {
		conn, err := ln.Accept()
		if err != nil {
			continue
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	buffer := make([]byte, 4096)
	n, err := conn.Read(buffer)
	if err != nil {
		return
	}

	msg := string(buffer[:n])
	fmt.Println("Recebido:", msg)

	partes := strings.Split(msg, "|")
	if len(partes) != 2 {
		fmt.Println("Formato inválido")
		return
	}

	id, err := strconv.Atoi(partes[0])
	if err != nil {
		return
	}

	valoresStr := strings.Split(partes[1], ",")

	valores := []int{}
	for _, v := range valoresStr {
		num, err := strconv.Atoi(v)
		if err == nil {
			valores = append(valores, num)
		}
	}

	mu.Lock()
	matrizFinal[id] = valores

	totalRecebido := len(matrizFinal)
	mu.Unlock()

	fmt.Printf("Progresso: %d/%d quadrantes\n", totalRecebido, GRID*GRID)

	// Quando todos chegarem
	if totalRecebido == GRID*GRID {
		fmt.Println("\nTodos os quadrantes recebidos!")

		matriz := montarMatriz(matrizFinal)

		mostrarMatriz2D(matriz)
		salvarPGM("imagem.pgm", matriz)
	}
}

// ================= MATRIZ =================

func montarMatriz(mapa map[int][]int) [][]int {

	partes := GRID * GRID
	tam := len(mapa[0])

	linhasQuad := 2
	colunasQuad := tam / linhasQuad

	alturaTotal := GRID * linhasQuad
	larguraTotal := GRID * colunasQuad

	matriz := make([][]int, alturaTotal)
	for i := range matriz {
		matriz[i] = make([]int, larguraTotal)
	}

	for id := 0; id < partes; id++ {

		q := mapa[id]

		rowBase := (id / GRID) * linhasQuad
		colBase := (id % GRID) * colunasQuad

		k := 0
		for i := 0; i < linhasQuad; i++ {
			for j := 0; j < colunasQuad; j++ {
				matriz[rowBase+i][colBase+j] = q[k]
				k++
			}
		}
	}

	return matriz
}

// ================= PRINT =================

func mostrarMatriz2D(matriz [][]int) {
	fmt.Println("\nMATRIZ FINAL:\n")

	for i := 0; i < len(matriz); i++ {
		for j := 0; j < len(matriz[i]); j++ {
			fmt.Printf("%3d ", matriz[i][j])
		}
		fmt.Println()
	}
}

// ================= IMAGEM =================

func salvarPGM(nome string, matriz [][]int) {
	file, err := os.Create(nome)
	if err != nil {
		fmt.Println("Erro ao criar arquivo")
		return
	}
	defer file.Close()

	altura := len(matriz)
	largura := len(matriz[0])

	fmt.Fprintf(file, "P2\n%d %d\n255\n", largura, altura)

	for i := 0; i < altura; i++ {
		for j := 0; j < largura; j++ {
			fmt.Fprintf(file, "%d ", matriz[i][j])
		}
		fmt.Fprintln(file)
	}

	fmt.Println("\nImagem salva como", nome)
}