FROM golang:1.21.1 AS builder

WORKDIR /app

COPY .. .

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o doshbank ./main.go
	
ENV GOPATH=/app

FROM scratch

COPY --from=builder /app/director /director


CMD ["/doshbank"]
