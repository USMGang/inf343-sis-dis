package main

import (
	"os"
	"strconv"

	f "l3/floors"
	g "l3/globals"
	u "l3/ui"

	"github.com/joho/godotenv"
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
    err := godotenv.Load()
    g.FailOnError(err, "Error al cargar el archivo .env")

    doshHost := os.Getenv("DOSHBANK_HOST")
    doshPort := os.Getenv("DOSHBANK_PORT")

    directorHost := os.Getenv("DIRECTOR_HOST")
    directorPort := os.Getenv("DIRECTOR_PORT")

    s.InitServer(N_MERCENARIES)
    defer s.Dosh.Conn.Close()
    defer s.Dosh.Ch.Close()

    fs.initDoshbank(doshHost, doshPort)

    s.Ui = u.NewUI(g.N_NOTIFICATIONS)
    s.Ui.ChangeOptions(g.DIRECTOR_PROMPT, g.DIRECTOR_OPTIONS)
    s.Ui.InitInterfaceChoice()
    showInterface()

    fs.setListener(directorHost, directorPort, &s)
    // TODO: LE falta el listener pal wea del NameNode

    // Interfaz del director
    go func(){
        <-quit
        fs.lis.Close()
        fs.grpcServer.GracefulStop()
        os.Exit(0)
    }()
}
