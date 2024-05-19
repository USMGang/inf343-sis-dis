package floors

import (
	"fmt"
	u "l3/ui"
    dosh "l3/doshbank_backend"
	"math/rand"
	"sync"

	"golang.org/x/net/context"
)

var (
	ID = 1
)

type Server struct {
	UnimplementedFloorsServiceServer
	NMercenaries       int
	CurrentMercenaries int
	Wait               chan bool
    MercenariesMutex   sync.Mutex
	Mutex              sync.Mutex
	Cond               *sync.Cond
	Ui                 u.UI

    // --- Doshbank ---
    Dosh dosh.DoshBank
}

// ================== Auxiliares ==================
func (s *Server) InitServer(n_mercenaries int) {
    s.NMercenaries = n_mercenaries
    s.CurrentMercenaries = 0

    s.Cond = sync.NewCond(&s.Mutex)
    s.Wait = make(chan bool, n_mercenaries)

    s.Dosh = dosh.DoshBank{}
    s.Dosh.InitDoshBank()
}


func (s *Server) waitMercenaries(id int, floor int){
    s.MercenariesMutex.Lock()
    s.CurrentMercenaries++
	s.Ui.AddNotification(fmt.Sprintf("Estoy listo para el piso %d. [Mercenario %d]", floor, id))
    s.MercenariesMutex.Unlock()

    s.Mutex.Lock()
    if s.CurrentMercenaries == s.NMercenaries {
		s.Ui.AddNotification(fmt.Sprintf("[Director] Todos los mercenarios estan listos para el piso %d", floor))
    } 
    s.Cond.Wait()
    s.Mutex.Unlock()
}

func (s *Server) StartMission(ctx context.Context, request *Start) (*Start, error) {
	s.MercenariesMutex.Lock()
	selected_id := ID
	s.Ui.AddNotification(fmt.Sprintf("[Director] Bienvenido mercenario %d", ID))
	ID++
	defer s.MercenariesMutex.Unlock()

	return &Start{Id: int32(selected_id)}, nil
}

func (s *Server) MercenaryReady(ctx context.Context, request *ReadyRequest) (*ReadyAnswer, error) {
	return &ReadyAnswer{Continue: true}, nil
}

// ================== Piso 1 ==================
func (s *Server) Floor1(ctx context.Context, request *Floor1ResultsRequest) (*Floor1ResultsAnswer, error) {
    s.waitMercenaries(int(request.Id), 1)

	// Aqui agregar la wea de las probabilidades
	x := rand.Intn(101)
	y := rand.Intn(101-x) + x

	var prob int
	switch request.SelectedWeapon {
	case 1:
		prob = x
	case 2:
		prob = y - x
	case 3:
		prob = 100 - y
	default:
		prob = x
	}

	if prob < int(request.RandNumber) {
		// Aqui va la logica para avisar al banco
		s.Ui.AddNotification(fmt.Sprintf("[Director] El mercenario %d ha muerto.", request.Id))
		s.NMercenaries--
        s.Dosh.Publish(int(request.Id), 1)
		return &Floor1ResultsAnswer{Message: "Has muerto... [Director]", IsDead: true}, nil
	}

	s.Ui.AddNotification(fmt.Sprintf("[Director] El mercenario %d ha superado el piso 1", request.Id))
	return &Floor1ResultsAnswer{Message: "Piso 1 superado! [Director]", IsDead: false}, nil
}

// ================== Piso 2 ==================
func (s *Server) Floor2(ctx context.Context, request *Floor2PathRequest) (*Floor2PathAnswer, error) {
    s.waitMercenaries(int(request.Id), 2)

	directors_choice := rand.Intn(2) + 1
	if directors_choice == int(request.SelectedPath) {
		s.Ui.AddNotification(fmt.Sprintf("[Director] El mercenario %d ha superado el piso 2.", request.Id))
		return &Floor2PathAnswer{Message: "Piso 2 superado! [Director]", IsOut: false}, nil
	}

	s.NMercenaries--
    s.Dosh.Publish(int(request.Id), 2)
	s.Ui.AddNotification(fmt.Sprintf("[Director] El mercenario %d está fuera de la mision.", request.Id))
	return &Floor2PathAnswer{Message: "Has sido traicionado... [xxx]", IsOut: true}, nil
}

// ================== Piso 3 ==================
func (s *Server) Floor3(ctx context.Context, request *Floor3Try) (*Floor3Try, error) {

	rand_number := rand.Intn(16) + 1
	if rand_number <= int(request.RandNumber) {
		request.NGoodTries++
	}

	request.NTries++

	return &Floor3Try{Id: request.Id, NTries: request.NTries, NGoodTries: request.NGoodTries, RandNumber: int32(rand_number)}, nil

}

func (s *Server) Floor3Results(ctx context.Context, request *Floor3ResultsRequest) (*Floor3ResultsAnswer, error) {
    s.waitMercenaries(int(request.Id), 2)

	if request.NGoodTries >= 2 {
		s.Ui.AddNotification(fmt.Sprintf("[Director] El mercenario %d ha superado el piso 3.", request.Id))
		return &Floor3ResultsAnswer{Message: "El patriarca ha sido derrotado! [Director]", IsDead: false}, nil
	}

	s.NMercenaries--
    s.Dosh.Publish(int(request.Id), 3)
	s.Ui.AddNotification(fmt.Sprintf("[Director] El mercenario %d ha muerto.", request.Id))
	return &Floor3ResultsAnswer{Message: "Has muerto... [Director]", IsDead: true}, nil
}
