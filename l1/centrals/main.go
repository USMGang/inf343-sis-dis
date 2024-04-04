package main

import (
    "os"
	"fmt"
    "net"
    "bufio"
	"math/rand"
	"pirates/globals"
)

// ===== Funciones Central =====
// func handle_asignation(){
//     
//     // While infinito
// }

func get_planet_treasures() map[string]int {
    planets := make(map[string]int)

    for i := 0; i < globals.N_leters; i++ {
        planets[string(globals.Base_leter + byte(i))] = rand.Intn(globals.Max_tresures)
    }

    return planets
}

func handlePetition(conn net.Conn) {
    defer conn.Close() // Asegúrate de cerrar la conexión al final.

    // Leer el mensaje del cliente
    message, err := bufio.NewReader(conn).ReadString('\n')
    if err != nil {
        fmt.Println("Cliente desconectado")
        return // Salir del bucle si hay un error (p. ej., el cliente se desconecta).
    }
    fmt.Print("Mensaje recibido:", string(message))

    // Enviar una respuesta al cliente
    respuesta := "Mensaje recibido de la central: " + message
    _, err = conn.Write([]byte(respuesta))
    if err != nil {
        fmt.Println(err)
    }
}

func main() {

    // Crear los planetas y asignarls los tesoros
    planets := get_planet_treasures()
    fmt.Println(planets)

    // Inicializar el listener en el puerto 8080
    ln, err := net.Listen("tcp", ":8080")
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
    defer ln.Close()
    fmt.Println("Central escuchando en el puerto 8080")

    // Bucle principal de la central para manejar las peticiones
    // for {
    //     // Aceptar una nueva conexión
    //     conn, err := ln.Accept()
    //     if err != nil {
    //         fmt.Println(err)
    //         continue
    //     }
    //     fmt.Println("Conexión establecida con", conn.RemoteAddr().String())

            // go handlePetition(conn)
    // }

    // Aceptar una nueva conexión
    for i := 0; i < 2; i++ {
        conn, err := ln.Accept()
        if err != nil {
            fmt.Println(err)
        }
        // fmt.Println("Conexión establecida con", conn.RemoteAddr().String())

        go handlePetition(conn)
    }


    
}
