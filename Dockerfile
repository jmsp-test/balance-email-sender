FROM golang:1.18.1

LABEL maintainer = "jmsp-tests <jmsosa.tests@gmail.com>"

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN go build -o ./build/ main.go

CMD ["./build/main"]