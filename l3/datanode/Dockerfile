FROM golang:1.21.1 AS builder

WORKDIR /app

COPY .. .

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o datanode datanode/main.go datanode/server.go datanode/interface.go
	
ENV GOPATH=/app

FROM scratch

COPY --from=builder /app/datanode /datanode


CMD ["/namenode"]
