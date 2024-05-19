N ?= 8
ID ?= 1

BIN_FOLDER=bin

all: proto build

proto:
	protoc -I=. --go_out=. --go-grpc_out=. floors/floors.proto
	protoc -I=. --go_out=. --go-grpc_out=. doshbank_backend/doshbank.proto
	protoc -I=. --go_out=. --go-grpc_out=. namenode_backend/namenode.proto
	protoc -I=. --go_out=. --go-grpc_out=. datanode_backend/datanode.proto

build: proto
	go build -o $(BIN_FOLDER)/mercenary mercenary/main.go mercenary/floors.go mercenary/interface.go
	go build -o $(BIN_FOLDER)/director director/main.go director/server.go director/interface.go
	go build -o $(BIN_FOLDER)/doshbank doshbank/main.go
	go build -o $(BIN_FOLDER)/namenode namenode/main.go
	go build -o $(BIN_FOLDER)/datanode datanode/main.go

player: build
	./$(BIN_FOLDER)/mercenary 1

bot: build
	./$(BIN_FOLDER)/mercenary 0

# director: build
# 	docker start rabbitmq
# 	./$(BIN_FOLDER)/director $(N)

build_director:
	docker compose build director

director: build_director
	docker compose up director

doshbank: build
	docker start rabbitmq
	./$(BIN_FOLDER)/doshbank

namenode: build
	./$(BIN_FOLDER)/namenode $(N)

datanode: build
	./$(BIN_FOLDER)/datanode $(ID)

install_rabbit:
	docker run -d --name rabbitmq -p 5672:5672 -p 15672:15672 rabbitmq:3-management

clear:
	rm -rf bin/*
	rm -rf txt/*.txt
	rm -rf txt/**/*.txt
	./bots.sh -1

test:
	make clear
	./bots.sh 2 3

test2:
	make clear
	./bots.sh 2 0

test3:
	make clear
	./bots.sh 0 4

proto1:
	protoc -I=. --go_out=. --go-grpc_out=. floors/floors.proto

proto2:
	protoc -I=. --go_out=. --go-grpc_out=. doshbank_backend/doshbank.proto

proto3:
	protoc -I=. --go_out=. --go-grpc_out=. namenode_backend/namenode.proto

proto4:
	protoc -I=. --go_out=. --go-grpc_out=. datanode_backend/datanode.proto

# docker: proto
# 	sudo docker build -t lab2_super_tierra .
#
# serverd: docker
# 	sudo docker run -d -p 8080:8080 lab2_super_tierra
#
# server: docker
# 	sudo docker run -p 8080:8080 lab2_super_tierra
#
# s:
# 	sudo docker run -p 8080:8080 lab2_super_tierra

.PHONY: all proto build player bot director doshbank namenode datanode install_rabbit clear test test2 test3 proto1 proto2 proto3 proto4
