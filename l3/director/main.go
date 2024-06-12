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
	s      = f.Server{}
	fs     = FloorsServers{}
	quit   = make(chan bool)
)

// comment
func main() {
	N_MERCENARIES, _ := strconv.Atoi(os.Args[1])

	// ================== Inicializar el servidor ==================
	doshHost := "10.35.169.80"
	doshPort := "8081"

	directorHost := "10.35.169.79"
	directorPort := "8080"

	s.InitServer(N_MERCENARIES)
	defer s.Dosh.Conn.Close()
	defer s.Dosh.Ch.Close()

	fs.initDoshbank(doshHost, doshPort)

	s.Ui = u.NewUI(g.N_NOTIFICATIONS)
	s.Ui.ChangeOptions(g.DIRECTOR_PROMPT, g.DIRECTOR_OPTIONS)
	s.Ui.InitInterfaceChoice()
	showInterface()

	fs.setListener(directorHost, directorPort, &s)

	// Interfaz del director
	go func() {
		<-quit
		fs.lis.Close()
		fs.grpcServer.GracefulStop()
		os.Exit(0)
	}()
}
