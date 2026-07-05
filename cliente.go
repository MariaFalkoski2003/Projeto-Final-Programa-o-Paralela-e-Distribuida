package main

import (
	"fmt"
	"math/rand"
	"net"
	"strconv"
	"time"
)

const GRID = 4 // 4x4 = 16 quadrantes

func main() {
	rand.Seed(time.Now().UnixNano())

	total := GRID * GRID

	for id := 0; id < total; id++ {

		conn, err := net.Dial("tcp", "127.0.0.1:8080")
		if err != nil {
			fmt.Println("Erro na conexão")
			return
		}

		// gera valores aleatórios
		valores := ""
		for i := 0; i < 8; i++ { // quantidade de elementos por quadrante
			valores += strconv.Itoa(rand.Intn(256))
			if i < 7 {
				valores += ","
			}
		}

		msg := fmt.Sprintf("%d|%s", id, valores)

		fmt.Println("Enviando:", msg)

		conn.Write([]byte(msg))
		conn.Close()

		time.Sleep(100 * time.Millisecond) // simula envio distribuído
	}

	fmt.Println("\nTodos os quadrantes enviados!")
}