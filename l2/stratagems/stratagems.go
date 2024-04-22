package stratagems

import (
	"fmt"
	"sync"
	"golang.org/x/net/context"
)

type SuperEarthSupply struct {
    AT int
    MP int
}

type Server struct {
    UnimplementedSuperEarthServiceServer
    Quit chan struct{}
    Data SuperEarthSupply
    Mutex sync.Mutex
}

func (s *Server) UpdateAmmo(at int, mp int) {
    s.Mutex.Lock()
    defer s.Mutex.Unlock()
    at_max := 50
    mp_max := 20

    if s.Data.AT + at >= at_max { s.Data.AT = at_max } else { s.Data.AT += at }
    if s.Data.MP + mp >= mp_max { s.Data.MP = mp_max } else { s.Data.MP += mp }
}

func (s *Server) Resupply(ctx context.Context, request *HelldiverRequest) (*SuperEarthAnswer, error) {
    required_at := request.At
    required_mp := request.Mp

    if (s.Data.AT >= int(required_at)) && (s.Data.MP >= int(required_mp)) {
        s.Mutex.Lock()
        s.Data.AT -= int(required_at)
        s.Data.MP -= int(required_mp)
        defer s.Mutex.Unlock()
        fmt.Printf("[Super Tierra] - [Reserva (AT: %d, MP: %d)]: Peticion de [%s] => AT: %d, MP: %d [ --- APROBADA --- ]\n", required_at, required_mp, request.Id, s.Data.AT, s.Data.MP)
	    return &SuperEarthAnswer{ Status: true }, nil
    } 

    fmt.Printf("[Super Tierra] - [Reserva (AT: %d, MP: %d)]: Peticion de [%s] => AT: %d, MP: %d [ --- RECHAZADA --- ]\n", required_at, required_mp, request.Id, s.Data.AT, s.Data.MP)
    return &SuperEarthAnswer{ Status: false }, nil
}

func (s*Server) EndServer(ctx context.Context, request *QuitMessage) (*QuitMessage, error) {
    fmt.Printf("%s\n", request.Message)
    close(s.Quit)
    return &QuitMessage{ Message: "Conexion con la Super Tierra cerrada, Mision Cumplida!" }, nil
}
