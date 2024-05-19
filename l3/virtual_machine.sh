#!/bin/bash

base=69
n=$1

read -p "Ingresa el usuario: " user
read -s -p "Ingresa la contraseña: " password

bash -c "ssh $user@capela.inf.santiago.usm.cl" <<< password
