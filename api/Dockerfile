FROM golang:1.16.2-buster

WORKDIR /api

COPY . .

RUN go mod download

RUN go get -u github.com/cosmtrek/air
RUN go get github.com/joho/godotenv@master

ENTRYPOINT air