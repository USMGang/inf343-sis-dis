package main

import (
	"os"
	"strconv"

	f "l3/floors"
	g "l3/globals"
	u "l3/ui"
)

var (
    choice = 0
    s = f.Server{}
    fs = FloorsServers{}
    quit = make(chan bool)
)

func main(){
    N_MERCENARIES, _ := strconv.Atoi(os.Args[1])

    // ================== Inicializar el servidor ==================

    s.InitServer(N_MERCENARIES)
    defer s.Dosh.Conn.Close()
    defer s.Dosh.Ch.Close()

    fs.initDoshbank("localhost", "8081")

    s.Ui = u.NewUI(g.N_NOTIFICATIONS)
    s.Ui.ChangeOptions(g.DIRECTOR_PROMPT, g.DIRECTOR_OPTIONS)
    s.Ui.InitInterfaceChoice()
    showInterface()

    fs.setListener("localhost", "8080", &s)


    // Interfaz del director
    go func(){
        <-quit
        fs.lis.Close()
        fs.grpcServer.GracefulStop()
        os.Exit(0)
    }()
}
