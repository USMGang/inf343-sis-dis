FROM golang:1.21.1 AS builder

WORKDIR /app

COPY .. .

RUN go mod download

RUN mkdir -p /app/txt /app/txt/datanode1 /app/txt/datanode2 /app/txt/datanode3
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o server datanode/main.go
	
ENV GOPATH=/app

FROM scratch

COPY --from=builder /app/server /server
COPY --from=builder /app/txt /app/txt

CMD ["/server"]
