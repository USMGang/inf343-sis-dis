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

    notifications = ui.NewUI(g.N_NOTIFICATIONS)
)

func main() {
	// ================== Inicialización ==================

	// Verificar si el mercenario sera controlado o un bot
	player, err := strconv.Atoi(os.Args[1])
	g.FailOnError(err, "Error, no se pudo recibir el tipo de jugador")

    namenodeHost := "10.35.169.81"
    namenodePort := "8070"

    doshHost := "10.35.169.80"
    doshPort := "8081"

    directorHost := "10.35.169.79"
    directorPort := "8080"


	// Administrar gestion con el director 
    server := FloorsServers{}
    server.initDirector(directorHost, directorPort)
    server.initNameNode(namenodeHost, namenodePort)
    server.initDoshbank(doshHost, doshPort)
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
