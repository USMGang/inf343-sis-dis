package main

import (
    // ui "l3/ui"
    d "l3/doshbank_backend"
)

func main(){
    dosh := d.DoshBank{}
    dosh.InitDoshBank()
	defer dosh.Conn.Close()
	defer dosh.Ch.Close()

    dosh.Consume()

    var forever chan struct{}
    go dosh.HandleDeadMercenary()

    // TODO: Copiar el codigo para TCP pero con otro puerto

    <-forever
}
