FROM golang:1.19

WORKDIR /usr/src/app
COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build -v -o ./go_object_store

CMD ["./go_object_store"]