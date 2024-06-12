package main

import (
	"fmt"
	d "l3/doshbank_backend"
	g "l3/globals"
	u "l3/ui"
	"net"
	"os"
    "log"

	"google.golang.org/grpc"
)

func main(){
    dirPath := "txt"
    if _, err := os.Stat(dirPath); os.IsNotExist(err) {
        if err := os.MkdirAll(dirPath, os.ModePerm); err != nil {
            log.Fatalf("Error al crear la carpeta: %v", err)
        }
        log.Println("Carpeta creada:", dirPath)
    }

    rabbitHost := "10.35.169.80"
    rabbitPort := "5672"

    doshHost := "10.35.169.80"
    doshPort := "8081"

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
