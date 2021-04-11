FROM golang:1.16.2-buster

WORKDIR /api

COPY . .

RUN apt-get update
RUN apt-get install -y python3 python3-pip
RUN pip3 install PyNaCl

RUN go mod download

RUN go get -u github.com/cosmtrek/air
RUN go get github.com/joho/godotenv@master

ENTRYPOINT air