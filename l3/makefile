N ?= 8

build_mercenary:
	go build -o mercenary/mercenary mercenary/main.go mercenary/floors.go mercenary/interface.go

player: build_mercenary
	./mercenary/mercenary 1

bot: build_mercenary
	./mercenary/mercenary 0

build_director:
	go build -o director/director director/main.go director/server.go director/interface.go

director: build_director
	docker start rabbitmq
	./director/director $(N)

doshbank:
	go run doshbank/main.go

rabbit:
	docker run -d --name rabbitmq -p 5672:5672 -p 15672:15672 rabbitmq:3-management

clear:
	# TODO: Agregar que borre todos los .txt
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

proto:
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

.PHONY: client proto docker server
