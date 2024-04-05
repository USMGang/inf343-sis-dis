package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net"
	"os"
	"pirates/globals"
	"strconv"
	"sync"
)

// ===== Funciones Piratas =====
func get_random_planet() string {
    random_leter := int(globals.Base_leter) + rand.Intn(globals.N_leters)
    return string(rune(random_leter))
}

func find_tresure(wg *sync.WaitGroup) {
	defer wg.Done()
	c, err := net.Dial("tcp", ":8080")
	if err != nil {
		fmt.Println(err)
		return
	}

    msg := globals.Message {
        Tresure: rand.Intn(globals.Max_tresures),
        Captain: "C" + strconv.Itoa(rand.Intn(3) + 1),
        Planet: get_random_planet(),
    }

    jsonData, err := json.Marshal(msg)
    if err != nil {
        fmt.Println(err)
        return
    }
    fmt.Printf("[Capitan]: Capitan %s encontro un botin de %d en el Planeta P%s, enviando solicitud de asignacion\n", msg.Captain, msg.Tresure, msg.Planet)

	_, err = c.Write(jsonData)
    if err != nil {
        fmt.Println(err)
    }
	c.Close()
}

func main(){
    // Conectarse al servidor TCP en el puerto 8080
    conn, err := net.Dial("tcp", "localhost:8080")
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
    defer conn.Close()

    var wg sync.WaitGroup
    wg.Add(globals.Task_counter)

    for i:=0; i<globals.Task_counter; i++ {
        go find_tresure(&wg)
    }
    wg.Wait()
    fmt.Println("[Capitan]: Los capitanes dejaron de buscar botines")
}
