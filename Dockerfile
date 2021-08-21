FROM golang:buster

WORKDIR /go/src/ewallet
COPY . .
RUN go mod download 
ENV MYSQL_HOST=host.docker.internal

CMD ["go", "run", "controller.go"]
