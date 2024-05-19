package main

import (
	"fmt"
	g "l3/globals"
	namenode "l3/namenode_backend"
	"net"
	"os"
	"strconv"

	"google.golang.org/grpc"
)

func main(){
	n_mercenaries, _:= strconv.Atoi(os.Args[1])
    ip_conn := fmt.Sprintf("%s:%s", os.Args[2], "8070")

    ips := make([]string, 3)
    for i := 0; i < 3; i++ {
        ips[i] = os.Args[3+i]
    }

    // ================== Inicializar el servidor ==================
    lis, err := net.Listen("tcp", ip_conn)
    g.FailOnError(err, fmt.Sprintf("Error, no se pudo establece el listener en: %s", ip_conn))

    s := namenode.Server{}
    s.InitServer(ips, n_mercenaries)
    grpcServer := grpc.NewServer()

    fmt.Println("Servidor gRPC iniciado en: ", ip_conn)

    namenode.RegisterNamenodeServiceServer(grpcServer, &s)

    err = grpcServer.Serve(lis)
    g.FailOnError(err, fmt.Sprintf("Error, no se pudo establecer el servidor gRPC en: %s", ip_conn))


}
