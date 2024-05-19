client:
	go run mercenary/main.go 1

bot:
	go run mercenary/main.go 0

server:
	go run director/main.go 1

doshbank:
	go run doshbank/main.go

rabbit:
	docker run -d --name rabbitmq -p 5672:5672 -p 15672:15672 rabbitmq:3-management

clear:
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
