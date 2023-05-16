FROM golang:1.19-bullseye

WORKDIR /app

COPY . .
RUN go mod download

COPY . .

RUN go build -o /go-docker

EXPOSE 8080

CMD [ "/go-docker" ]