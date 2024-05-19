package doshbank_backend

import (
	"fmt"
	g "l3/globals"
	u "l3/ui"
	"os"
	sync "sync"

	"encoding/json"

	amqp "github.com/streadway/amqp"
	"golang.org/x/net/context"
)

type DoshBank struct {
	// --- General ---
	Reward int
	Ui     u.UI

	// --- RabbitMQ ---
	Ch   *amqp.Channel
	Conn *amqp.Connection
	q    amqp.Queue
	msgs <-chan amqp.Delivery
	mu   sync.Mutex

	// --- gRPC ---
    UnimplementedDoshBankServer
}

type Signal struct {
	Id    int `json:"id"`
	Floor int `json:"floor"`
}


// ================== gRPC ==================
func (s *DoshBank) GetCurrentReward(ctx context.Context, request *GetCurrentRewardRequest) (*GetCurrentRewardResponse, error) {
    return &GetCurrentRewardResponse{Reward: int32(s.Reward)}, nil
}

// ================== RabbitMQ ==================
func (d *DoshBank) InitDoshBank(ip string, port string) {
	var err error
    ip_conn := fmt.Sprintf("amqp://guest:guest@%s:%s/", ip, port)
    d.Conn, err = amqp.Dial(ip_conn)
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

func (d *DoshBank) Consume() {
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
    file, err := os.OpenFile("txt/doshbank.txt", os.O_APPEND|os.O_WRONLY, 0644)
    g.FailOnError(err, "Fallo al abrir el archivo")
    defer file.Close()

	for s := range d.msgs {
		var signal Signal
		err = json.Unmarshal(s.Body, &signal)
		g.FailOnError(err, "Error al transformar el mensaje a JSON")

        d.mu.Lock()
		d.Reward += g.REWARD_BONUS
        d.Ui.AddNotification(fmt.Sprintf("%s ha muerto en el piso %d - Botin actual: %d", g.GetName(signal.Id), signal.Floor, d.Reward))
        _, err = file.WriteString(fmt.Sprintf("%s %d %d\n", g.GetName(signal.Id), signal.Floor, d.Reward))
        g.FailOnError(err, "Fallo al escribir en el archivo")
        d.mu.Unlock()
	}
}
