package main

import (
	d "l3/doshbank_backend"
	g "l3/globals"
	u "l3/ui"
	"net"
	"os"

	"google.golang.org/grpc"
)

func main(){
    // ================== Archivo ==================
    file, err := os.Create("txt/doshbank.txt")
    g.FailOnError(err, "Fallo al crear el archivo")
    file.Close()

    // ================== RabbitMQ ==================
    dosh := d.DoshBank{}
    // TODO: Contolar esto por args
    dosh.InitDoshBank()
	defer dosh.Conn.Close()
	defer dosh.Ch.Close()

    dosh.Consume()

    go dosh.HandleDeadMercenary()

    dosh.Ui = u.NewUI(g.N_NOTIFICATIONS)
    dosh.Ui.ChangeOptions(g.VOID_PROMPT, g.VOID_OPTIONS)
    dosh.Ui.AddNotification("[DoshBank] Iniciando el doshbank...")

    // ================== gRPC ==================
    // TODO: Contolar esto por args
    lis, err := net.Listen("tcp", ":8081")
    g.FailOnError(err, "Fallo al escuchar el puerto 8081")

    grpcServer := grpc.NewServer()

    d.RegisterDoshBankServer(grpcServer, &dosh)
    err = grpcServer.Serve(lis)
    g.FailOnError(err, "Fallo al ejecutar grcp en el puerto 8081")

}
