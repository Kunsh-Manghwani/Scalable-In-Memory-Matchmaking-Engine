FROM golang:1.21

WORKDIR /app

COPY go.mod . 
COPY go.sum .
RUN go mod download

COPY . .

RUN go build -o matchmaking cmd/main.go

EXPOSE 8080

CMD ["./matchmaking"]
