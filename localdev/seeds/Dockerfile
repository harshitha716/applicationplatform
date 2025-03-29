FROM golang:1.22.5-alpine

WORKDIR /go/src/app

COPY go.* .

RUN go mod download

COPY . .

CMD ["go", "run", "main.go"]
