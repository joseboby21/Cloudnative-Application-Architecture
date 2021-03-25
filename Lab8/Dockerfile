FROM golang:1.14-alpine AS build
WORKDIR /src/
RUN apk add git
RUN go mod init src
RUN go get go.mongodb.org/mongo-driver
COPY webserver.go .
RUN CGO_ENABLED=0 go build -o /bin/server
EXPOSE 8000
CMD ["server"]
