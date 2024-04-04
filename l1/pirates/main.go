package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"net"
	"os"
	"pirates/globals"
)

func get_random_leter() string {
    random_leter := int(globals.Base_leter) + rand.Intn(globals.N_leters)
    return string(rune(random_leter))
}

// ===== Funciones Piratas =====
func find_tresure(name string){ 
    fmt.Printf("Capitan %s encontro botin en el planeta P%s, enviando solicitud de asignacion\n", name, get_random_leter())

    // Enviar solicitud a la central

    // Esperar la wea

    // Recibir cuanto tesoro le toco pal planeta
}

func sendMessage(conn net.Conn, mensaje string) {
    fmt.Fprintf(conn, mensaje+"\n")

}

// recibirMensaje se encarga de leer y mostrar un mensaje del servidor.
func getAnswer(conn net.Conn) {
    respuesta, err := bufio.NewReader(conn).ReadString('\n')
    if err != nil {
        fmt.Println("Error al recibir mensaje:", err)
        return
    }

    fmt.Print("Mensaje del servidor: ", respuesta)
}

func main(){
    // Conectarse al servidor TCP en el puerto 8080
    conn, err := net.Dial("tcp", "localhost:8080")
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
    defer conn.Close()


    for i:=0; i<2; i++ {
        // Enviar un mensaje al servidor
        sendMessage(conn, "Mensaje")

        // Leer la respuesta del servidor
        getAnswer(conn)
    }
    // Enviar un mensaje al servidor


    // Leer la respuesta del servidor

    find_tresure("Pedro")
}
