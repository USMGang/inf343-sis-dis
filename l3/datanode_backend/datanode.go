package datanode_backend

import (
	"bufio"
	"fmt"
	g "l3/globals"
	"os"
	"sync"

	"golang.org/x/net/context"
)


type Server struct {
    UnimplementedDatanodeServiceServer
    Id int
    mutex sync.Mutex
}

func (s *Server) SaveStep(ctx context.Context, request *SaveStepRequest) (*SaveStepResponse, error) {
    s.mutex.Lock()
    name := fmt.Sprintf("datanode%d/%s_%d.txt", s.Id, g.GetName(int(request.Id)), request.Floor)

    file, err := os.OpenFile(name, os.O_APPEND|os.O_WRONLY, 0644)
    g.FailOnError(err, "Fallo al abrir el archivo")
    defer file.Close()


    for _, step := range request.Steps {
        _, err = file.WriteString(fmt.Sprintf("* %s\n", step))
        g.FailOnError(err, "Fallo al escribir en el archivo")
    }
    s.mutex.Unlock()

    return &SaveStepResponse{}, nil
}

func (s *Server) GetIdStepts(ctx context.Context, request *GetIdSteptsRequest) (*GetIdSteptsResponse, error) {
    s.mutex.Lock()
    name := fmt.Sprintf("datanode%d/%s_%d.txt", s.Id, g.GetName(int(request.Id)), request.Floor)

    file, err := os.OpenFile(name, os.O_APPEND|os.O_WRONLY, 0644)
    g.FailOnError(err, "Fallo al abrir el archivo")
    defer file.Close()

    var steps []string
    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        steps = append(steps, scanner.Text())
    }

    // TODO: Aqui deberiamos recupear el []string segun la id y el piso del matewea
    return &GetIdSteptsResponse{ Steps: steps }, nil
}


