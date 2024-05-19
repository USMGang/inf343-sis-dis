package main

import (
	"fmt"
	d "l3/doshbank_backend"
	g "l3/globals"
	u "l3/ui"
	"net"
	"os"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"
)

func main(){
    // ================== .env ==================
    err := godotenv.Load()
    g.FailOnError(err, "Error al cargar el archivo .env")

    rabbitHost := os.Getenv("RABBITMQ_HOST")
    rabbitPort := os.Getenv("RABBITMQ_PORT")

    doshHost := os.Getenv("DOSHBANK_HOST")
    doshPort := os.Getenv("DOSHBANK_PORT")

    // ================== Archivo ==================
    file, err := os.Create("txt/doshbank.txt")
    g.FailOnError(err, "Fallo al crear el archivo")
    file.Close()


    // ================== RabbitMQ ==================
    dosh := d.DoshBank{}

    dosh.InitDoshBank(rabbitHost, rabbitPort)
	defer dosh.Conn.Close()
	defer dosh.Ch.Close()
    
    dosh.Ui = u.NewUI(g.N_NOTIFICATIONS)
    dosh.Ui.ChangeOptions(g.VOID_PROMPT, g.VOID_OPTIONS)
    dosh.Ui.AddNotification("[DoshBank] Iniciando el doshbank...")

    dosh.Consume()

    go dosh.HandleDeadMercenary()

    // ================== gRPC ==================
    // TODO: Contolar esto por args
    lis, err := net.Listen("tcp", fmt.Sprintf("%s:%s", doshHost, doshPort))
    g.FailOnError(err, "Fallo al escuchar el puerto 8081")

    grpcServer := grpc.NewServer()

    d.RegisterDoshBankServer(grpcServer, &dosh)
    err = grpcServer.Serve(lis)
    g.FailOnError(err, "Fallo al ejecutar grcp en el puerto 8081")

}
