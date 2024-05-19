#!/bin/bash

nPlayers=$1
nBots=$2
nMercenaries=$((nBots + nPlayers))
DIR="$HOME/code/usm/inf343-sis-dis/l3"
CURRENT_SESSION=$(tmux display-message -p '#S')


# Cerrar todas las ventanas
if [ "$nPlayers" -eq -1 ]; then
    tmux list-windows -t $CURRENT_SESSION -F "#{window_name}" | grep -i "bot" | while read -r WINDOW_NAME
    do
        tmux kill-window -t $CURRENT_SESSION:"$WINDOW_NAME"
    done

    tmux list-windows -t $CURRENT_SESSION -F "#{window_name}" | grep -i "player" | while read -r WINDOW_NAME
    do
        tmux kill-window -t $CURRENT_SESSION:"$WINDOW_NAME"
    done

    tmux kill-window -t $CURRENT_SESSION:"director"
    tmux kill-window -t $CURRENT_SESSION:"doshbank"
    exit 0
fi

WINDOWS_NAME="director"
tmux new-window -t $CURRENT_SESSION -n $WINDOWS_NAME -c $DIR -S
tmux send-keys -t $CURRENT_SESSION:$WINDOWS_NAME "make director N_MERCENARIES=$nMercenaries" C-m

WINDOWS_NAME="doshbank"
tmux new-window -t $CURRENT_SESSION -n $WINDOWS_NAME -c $DIR -S
tmux send-keys -t $CURRENT_SESSION:$WINDOWS_NAME "go run doshbank/main.go" C-m

sleep 1

# Crear los players
for (( i=1; i<=$nPlayers; i++ ))
do
    WINDOWS_NAME="player_$i"
    tmux new-window -t $CURRENT_SESSION -n $WINDOWS_NAME -c $DIR -S
    tmux send-keys -t $CURRENT_SESSION:$WINDOWS_NAME "make player" C-m
done

# Crear los bots
for (( i=1; i<=$nBots; i++ ))
do
    WINDOWS_NAME="bot_$i"
    tmux new-window -t $CURRENT_SESSION -n $WINDOWS_NAME -c $DIR -S
    tmux send-keys -t $CURRENT_SESSION:$WINDOWS_NAME "make bot" C-m
done

