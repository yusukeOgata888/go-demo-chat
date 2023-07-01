FROM golang:1.16-alpine as builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

#build
COPY ./src/chat .

RUN pwd
RUN ls -l .

RUN go build -o .
EXPOSE 8080

CMD [ "/go-demo-chat" ]


