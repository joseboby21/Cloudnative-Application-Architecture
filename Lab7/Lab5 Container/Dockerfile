FROM golang:1.14-alpine AS build
WORKDIR /src/Lab5
COPY Lab5/movieserver/server.go /src/Lab5
RUN mkdir movieapi
COPY Lab5/movieapi /src/Lab5/movieapi
COPY go.mod /src/Lab5
COPY go.sum /src/Lab5
RUN go mod download

RUN CGO_ENABLED=0 go build -o /bin/server

FROM scratch
COPY --from=build /bin/server /bin/server
ENTRYPOINT ["/bin/server"]


