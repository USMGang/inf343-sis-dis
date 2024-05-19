package doshbank

import (
	"encoding/json"
	"fmt"
	g "l3/globals"

	amqp "github.com/streadway/amqp"
)

type DoshBank struct {
    // --- General ---
    Reward int

    // --- RabbitMQ ---
	Ch     *amqp.Channel
	Conn   *amqp.Connection
	q      amqp.Queue
	msgs   <-chan amqp.Delivery

    // --- gRPC ---
}

type Signal struct {
	Id     int `json:"id"`
	Floor  int `json:"floor"`
}

func (d *DoshBank) InitDoshBank() { 
	var err error
	d.Conn, err = amqp.Dial("amqp://guest:guest@localhost:5672/")
	g.FailOnError(err, "Error al conectar con RabbitMQ")

	d.Ch, err = d.Conn.Channel()
	g.FailOnError(err, "Error al abrir un canal")

	d.q, err = d.Ch.QueueDeclare(
		"doshBank", // name
		false,      // durable
		false,      // delete when unused
		false,      // exclusive
		false,      // no-wait
		nil,        // arguments
	)
	g.FailOnError(err, "Error al declarar una cola")
}

func (d *DoshBank) Publish(id int, floor int) {
	body := Signal{Id: id, Floor: floor}
	jsonBody, err := json.Marshal(body)
	g.FailOnError(err, "Error al transformar la señal a JSON")

	err = d.Ch.Publish(
		"",       // exchange
		d.q.Name, // routing key
		false,    // mandatory
		false,    // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        []byte(jsonBody),
		})
	g.FailOnError(err, "Error al publicar un mensaje")
}

func (d *DoshBank) Consume(){
	var err error
	d.msgs, err = d.Ch.Consume(
		d.q.Name, // queue
		"",       // consumer
		true,     // auto-ack
		false,    // exclusive
		false,    // no-local
		false,    // no-wait
		nil,      // args
	)
	g.FailOnError(err, "Error al registrar un consumidor")
}

func (d *DoshBank) HandleDeadMercenary() {
    for s := range d.msgs {
        var signal Signal
        err := json.Unmarshal(s.Body, &signal)
        g.FailOnError(err, "Error al transformar el mensaje a JSON")

        d.Reward += g.REWARD_BONUS
        fmt.Printf("Signal received: %v - %d\n", signal, d.Reward)

        // TODO: Crear el archivo con los datos de signal y la reward
    }
}
