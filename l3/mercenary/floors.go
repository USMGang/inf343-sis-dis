package main

import (
	"context"
	"fmt"
	d "l3/doshbank_backend"
	f "l3/floors"
	g "l3/globals"
	"math/rand"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type FloorsServers struct {
	// ================== General ==================
	Id     int
	Player int

	// ================== Director ==================
	DirectorConn *grpc.ClientConn
	Director     f.FloorsServiceClient

	// ================== Doshbank ==================
	DoshbankConn *grpc.ClientConn
	Doshbank     d.DoshBankClient
}

// ================== Auxiliar ==================
func Death() {
	<-is_dead
	os.Exit(0)
}

func (s *FloorsServers) setPlayer(player int) {
	s.Player = player
}

// ================== Inicializar los servidores ==================
func (s *FloorsServers) initDirector(ip string, port string) {
	ip_conn := fmt.Sprintf("%s:%s", ip, port)
	conn, err := grpc.NewClient(ip_conn, grpc.WithTransportCredentials(insecure.NewCredentials()))

	g.FailOnError(err, "Error, no se pudo establecer comunicación con el director")

	s.DirectorConn = conn
	s.Director = f.NewFloorsServiceClient(conn)
}

func (s *FloorsServers) initDoshbank(ip string, port string) {
	ip_conn := fmt.Sprintf("%s:%s", ip, port)
	conn, err := grpc.NewClient(ip_conn, grpc.WithTransportCredentials(insecure.NewCredentials()))
	g.FailOnError(err, "Error, no se pudo establecer comunicación con el director")

	s.DoshbankConn = conn
	s.Doshbank = d.NewDoshBankClient(conn)
}

// ================== Auxiliares para Pisos  ==================
func (s *FloorsServers) isReady(id int32, is_ready bool) bool {
	ready := f.ReadyRequest{Id: id, IsReady: true}
	response, err := s.Director.MercenaryReady(context.Background(), &ready)
	g.FailOnError(err, "Error, no se pudo recibir el mensaje")

	return (response.Continue)
}

func (s *FloorsServers) getAnswer(prompt string, n_bot_choices int) {
	if s.Player == 1 {
		<-get_answer

		if prompt != "" {
			notifications.AddNotification(prompt)
		} else {
			notifications.ShowNotifications()
		}

		choice = notifications.GetUserChoice()
		get_choice <- true

	} else {
		choice = getBotChoice(n_bot_choices)
	}
}

func (s *FloorsServers) confirmReady() {
	if s.Player == 1 {
		notifications.AddNotification("Mande 1 si está listo para el siguiente... ")

		<-get_answer
		notifications.ShowNotifications()
		choice = notifications.GetUserChoice()
		get_choice <- true

		if choice == 1 {
			s.isReady(int32(s.Id), true)
		}

	} else {
		s.isReady(int32(s.Id), true)
	}
}

// ================== Pisos ==================
func (s *FloorsServers) startMission() {
	// === Obtener id ===
	request_id := f.Start{Id: -1}
	response_id, err := s.Director.StartMission(context.Background(), &request_id)
	g.FailOnError(err, "Error, no se pudo recibir la id")

	s.Id = int(response_id.Id)
	notifications.AddNotification(fmt.Sprintf("ID de la mision: %d", s.Id))
}

func (s *FloorsServers) floor1() {
	// === Obtener las armas disponibles ===
	weapons := []string{"Escopeta", "Rifle automático", "Puños eléctricos"}
	notifications.AddNotification("Entrando al piso 1")
	notifications.AddNotification("Armas disponibles:")
	for i := 1; i <= len(weapons); i++ {
		notifications.AddNotification(fmt.Sprintf("    - [%d] %s", i, weapons[i-1]))
	}

	// === Seleccionar el arma y obtener el resultado ===
	s.getAnswer("Seleccione el arma que desea usar...", 3)
	notifications.AddNotification(fmt.Sprintf("Arma seleccionada: %s", weapons[choice-1]))

	// Mandar el arma elegida y esperar el resultado
	request_results := f.Floor1ResultsRequest{Id: int32(s.Id), SelectedWeapon: int32(choice), RandNumber: int32(rand.Intn(101))}
	response_results, err := s.Director.Floor1(context.Background(), &request_results)
	g.FailOnError(err, "Error, no se pudo recibir el resultado")

	notifications.AddNotification(response_results.Message)

	// === Verificar si el mercenario murio ===
	if response_results.IsDead {
		is_dead <- true
		return
	}

	// === Confirmar que se esta listo para continuar ===
	s.confirmReady()

}

func (s *FloorsServers) floor2() {
	// === Elegir el camino ===
	notifications.AddNotification("Entrando al piso 2")
	notifications.AddNotification("Recorriendo el piso te encuentras con 2 caminos:")
	notifications.AddNotification("- [1] Camino a la izquierda")
	notifications.AddNotification("- [2] Camino a la derecha")

	s.getAnswer("Elige tu camino...", 2)

	// === Ver si el camino es correcto ===
	notifications.AddNotification(fmt.Sprintf("Camino seleccionado: %d", choice))
	path_request := f.Floor2PathRequest{Id: int32(s.Id), SelectedPath: int32(choice)}
	path_response, err := s.Director.Floor2(context.Background(), &path_request)
	g.FailOnError(err, "Error, no se pudo recibir el camino")

	notifications.AddNotification(path_response.Message)

	// === Verificar si el mercenario fue traicionado ===
	if path_response.IsOut {
		is_dead <- true
		return
	}

	s.confirmReady()
}

func (s *FloorsServers) floor3() {
	// === Realizar las 5 rondas ===
	n_good_tries := 0
	var rand_number int

	notifications.AddNotification("Entrando al piso 3")
	for i := 1; i <= 5; i++ {
		rand_number = rand.Intn(16) + 1
		notifications.AddNotification(fmt.Sprintf("Intento %d: %d", i, rand_number))
		tries_request := f.Floor3Try{Id: int32(s.Id), NTries: int32(i), NGoodTries: int32(n_good_tries), RandNumber: int32(rand_number)}
		tries_response, err := s.Director.Floor3(context.Background(), &tries_request)
		g.FailOnError(err, "Error, no se pudo realizar la ronda")

		n_good_tries = int(tries_response.NGoodTries)
	}

	// === Revisar el resultado de las rondas ===
	floor3_request := f.Floor3ResultsRequest{Id: int32(s.Id), NGoodTries: int32(n_good_tries)}
	floor3_response, err := s.Director.Floor3Results(context.Background(), &floor3_request)
	g.FailOnError(err, "Error, no se pudo recibir el resultado")

	notifications.AddNotification(floor3_response.Message)

	// === Verificar si el mercenario murio ===
	if floor3_response.IsDead {
		is_dead <- true
		return
	}
}
