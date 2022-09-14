FROM golang:1.13.8

RUN mkdir /impulse-api
WORKDIR /impulse-api

COPY go.mod go.mod
COPY go.sum go.sum

RUN go mod download

RUN mkdir /app
COPY /app/*.go app/
RUN mkdir /config
COPY /config/.env config/
RUN mkdir /internal \
    mkdir /internal/entity \
COPY /internal/entity/*.go internal/entity/
RUN mkdir /pkg \
    mkdir /pkg/handler \
COPY /pkg/handler/*.go pkg/handler/
RUN mkdir /pkg/repository
COPY /pkg/repository/*.go pkg/repository/
RUN mkdir /pkg/service
COPY /pkg/service/*.go pkg/service/
COPY Aspects.txt .
COPY houses_upr.json .
COPY planets_power.json .
COPY server.go .

RUN GOOS=linux GOARCH=amd64 go build -o sendpulse app/main.go

CMD ["./sendpulse"]