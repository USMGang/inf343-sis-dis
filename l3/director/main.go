package main

import (
	"fmt"
	"os"
	"strconv"

	f "l3/floors"
	g "l3/globals"
	u "l3/ui"
)

var (
    s f.Server
    fs FloorsServers
    quit = make(chan bool)
)

func main(){
    N_MERCENARIES, _ := strconv.Atoi(os.Args[1])

    // ================== Inicializar el servidor ==================
    s := f.Server{}
    s.InitServer(N_MERCENARIES)

    s.Ui = u.NewUI(g.N_NOTIFICATIONS)
    fmt.Printf("equide: %d\n", s.Ui.Width)
    s.Ui.ChangeOptions(g.DIRECTOR_PROMPT, g.DIRECTOR_OPTIONS)
    showInterface()

    defer s.Dosh.Conn.Close()
    defer s.Dosh.Ch.Close()


    fs := FloorsServers{}
    fs.setListener("localhost", "8080", &s)
    fs.initDoshbank("localhost", "8081")

    // Interfaz del director
    go func(){
        <-quit
        fs.lis.Close()
        fs.grpcServer.GracefulStop()
        os.Exit(0)
    }()
}
