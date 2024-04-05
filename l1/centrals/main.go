package main

import (
	"fmt"
	"math"
	"math/rand"
	"net"
    "encoding/json"
	"pirates/globals"
	"sync"
)


func get_planet_treasures() map[string]int {
    planets := make(map[string]int)

    for i := 0; i < globals.N_leters; i++ {
        planets[string(globals.Base_leter + byte(i))] = rand.Intn(globals.Max_tresures)
    }

    return planets
}

func central(planetas map[string]int, wg *sync.WaitGroup) {
	defer wg.Done()
	s, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println(err)
		return
	}

	for {
		conn, err := s.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}

        handleSolicitude(conn, planetas)
	}
}


func handleSolicitude(conn net.Conn, planetas map[string]int) {
	defer conn.Close()

    var msg globals.Message
    err := json.NewDecoder(conn).Decode(&msg)
    if err != nil {
        fmt.Println("[Central]: Error al decodificar el mensaje", err)
    }

    fmt.Printf("[Central]: Recepcion de solicitud desde el Planeta P%s, del capitan %s \n", msg.Planet, msg.Captain)

    aux := math.Inf(1)
    min_planeta := ""
    for planeta, botin := range planetas {
        if float64(botin) < aux {
            aux = float64(botin)
            min_planeta = planeta
        }
    }
    planetas[min_planeta] += msg.Tresure

    fmt.Printf("[Central]: Botin asignado al planeta P%s, cantidad actual: %d\n", min_planeta, planetas[min_planeta])
    fmt.Println("[Central]: Botin Planetas: ", planetas)
}


func main() {

    // Crear los planetas y asignarls los tesoros
    planetas := get_planet_treasures()
    fmt.Println("[Central]: Botin planetas: ", planetas)

    var wg sync.WaitGroup
    wg.Add(1)
    go central(planetas, &wg)
    wg.Wait()
}
