FROM golang:1.21.1 AS builder

WORKDIR /app

COPY ../go.mod ../go.sum ./
RUN go mod download

COPY .. .
COPY .env .env

# RUN go build -o director director/main.go director/server.go director/interface.go
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o server director/main.go director/server.go director/interface.go
	
# ENV GOPATH=/app

FROM scratch

COPY --from=builder /app/server /server 
COPY .env /app/.env

CMD ["/server", "8"]
