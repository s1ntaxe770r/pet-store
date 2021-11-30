FROM golang:latest as builder

WORKDIR /app
COPY . .

RUN cd category && go mod download

RUN cd category && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app

EXPOSE 3000

ENTRYPOINT ["/app/category/app"]