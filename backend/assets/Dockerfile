FROM golang:1.23.4-alpine

WORKDIR /usr/src/app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

ARG MAXPROCS=1
COPY . .
RUN go build -v -o /usr/local/bin/gh-asset -gcflags=GOMAXPROCS=$MAXPROCS .

CMD ["/usr/local/bin/gh-asset"]
