package main

import (
	g "l3/globals"
    ui "l3/ui"
	"os"
	"strconv"
)

var (
    // name = ""
    choice = 0
    is_dead = make(chan bool)
    get_answer = make(chan bool)
    get_choice = make(chan bool)

    notifications = ui.NewUI(g.WIDTH, g.N_NOTIFICATIONS)
)

func main() {
	// ================== Inicialización ==================

	// Verificar si el mercenario sera controlado o un bot
	player, err := strconv.Atoi(os.Args[1])
	g.FailOnError(err, "Error, no se pudo recibir el tipo de jugador")

	// Administrar gestion con el director 
    server := FloorsServers{}
    server.initDirector("localhost", "8080")
    server.initDoshbank("localhost", "8081")
    server.setPlayer(player)

    // Dejar la señal para cuando el mercenario muera
    go Death()

	// ================== Interfaz ==================
	ShowInterface(player)

	// ================== Preparacion ==================
    server.startMission()

	// ================== Piso 1 ==================
    server.floor1()

	// ================== Piso 2 ==================
    server.floor2()

	// ================== Piso 3 ==================
    server.floor3()

}
