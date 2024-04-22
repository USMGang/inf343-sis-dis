client:
	go run helldiver/main.go

proto:
	protoc -I=. --go_out=. --go-grpc_out=. stratagems/stratagems.proto

docker: proto
	sudo docker build -t lab2_super_tierra .

serverd: docker
	sudo docker run -d -p 8080:8080 lab2_super_tierra

server: docker
	sudo docker run -p 8080:8080 lab2_super_tierra

s:
	sudo docker run -p 8080:8080 lab2_super_tierra

.PHONY: client proto docker server
