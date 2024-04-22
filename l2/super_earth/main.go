package main

import (
	"log"
	"net"
	"time"
	"l2/stratagems"
	"google.golang.org/grpc"
)

func makeAmmo(s *stratagems.Server, quit chan struct{}){
    for {
        select {
        case <- quit:
            return
        default: 
            time.Sleep(10 * time.Second)
            s.UpdateAmmo(10, 5)
        }
    }
}

func closeServer(server *grpc.Server, listener net.Listener, quit chan struct{}){
    <-quit
    server.GracefulStop()
    listener.Close()
}

func main(){
    lis, err := net.Listen("tcp", ":8080")
    if err != nil {
        log.Fatalf("Fallo al escuchar el puerto 8080: %v", err)
    }

    s := stratagems.Server{}
    s.Quit = make(chan struct{})
    s.Data = stratagems.SuperEarthSupply{AT: 0, MP: 0}

    quit := make(chan struct{})
    defer close(quit)
    go makeAmmo(&s, quit)

    grpcServer := grpc.NewServer()

    stratagems.RegisterSuperEarthServiceServer(grpcServer, &s)

    go closeServer(grpcServer, lis, s.Quit)

    if err := grpcServer.Serve(lis); err != nil {
        log.Fatalf("Fallo al ejecutar grcp en el puerto 8080: %v", err)
    }
} 
