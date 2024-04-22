package main

import (
	"fmt"
	"l2/stratagems"
	"log"
	"math/rand"
	"sync"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Helldivers struct {
    id string
    planet string
    at int32
    mp int32
}

func set_helldivers() []Helldivers {
    helldivers := make([]Helldivers, 4)

    planets := []string{"Malevelon Creek", "Tibit", "Mantes", "Ingmar"}
    rand.Shuffle(len(planets), func(i, j int) { planets[i], planets[j] = planets[j], planets[i] })
    
    for i:=0; i<4; i++ {
        at := int32(20 + rand.Intn(11))
        mp := int32(10 + rand.Intn(6))
        random_leter := string(rune(int(byte('A')) + rand.Intn(20)))
        helldivers[i] = Helldivers {
            id: fmt.Sprintf("%s%d", random_leter, i+1),
            planet: planets[i],
            at: at,
            mp: mp,
        }
    }

    return helldivers
}

func spreed_democracy(h Helldivers, c stratagems.SuperEarthServiceClient, wg *sync.WaitGroup) {
    defer wg.Done()

    initial_timer := time.NewTimer(time.Duration(10) * time.Second)
    <-initial_timer.C

    request := stratagems.HelldiverRequest{ Id: h.id, At: h.at, Mp: h.mp, }

    for {
        response, err := c.Resupply(context.Background(), &request)
        if err != nil {
            log.Fatalf("Error, estratagema no disponible: %s", err)
        }

        status_aux := response.Status
        if status_aux {
            fmt.Printf("[%s] Solicitando => [%d AT y %d MP]; [ --- APROBADA --- ]; %s conquistado!. Cerrando comunicación.\n", h.id, h.at, h.mp, h.planet)
            return;
        }

        fmt.Printf("[%s] Solicitando => [%d AT y %d MP]; [ --- DENEGADA --- ]; Reintentando en 3 segundos...\n", h.id, h.at, h.mp)
        time.Sleep(3 * time.Second)
    }
}

func main() {
    helldivers := set_helldivers()

    var conn *grpc.ClientConn
    conn, err := grpc.Dial(":8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
    if err != nil {
        log.Fatalf("No se pudo establecer comunicación con la Super Tierra: %s", err)
    }

    defer conn.Close()

    var wg sync.WaitGroup
    c := stratagems.NewSuperEarthServiceClient(conn)

    n_helldivers := 4
    wg.Add(n_helldivers)
    for i:=0; i<n_helldivers; i++ {
        go spreed_democracy(helldivers[i], c, &wg)
    }

    wg.Wait()

    // Mandar la señal de cierre para la Super Tierra
    request := stratagems.QuitMessage{ Message: "Conexion cerrada" }
    response, err := c.EndServer(context.Background(), &request)
    if err != nil {
        log.Fatalf("Error, estratagema no disponible: %s", err)
    }
    fmt.Println(response.Message)

}
