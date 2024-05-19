FROM golang:1.21.1 AS builder

WORKDIR /app

COPY .. .

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o namenode namenode/main.go namenode/server.go namenode/interface.go
	
ENV GOPATH=/app

FROM scratch

COPY --from=builder /app/namenode /namenode


CMD ["/namenode", "8"]
