FROM golang:1.16 AS base
  WORKDIR /usr/src/app

  COPY . .

  RUN go mod vendor

FROM base AS development
  RUN go get -v github.com/cortesi/modd/cmd/modd
  RUN go mod tidy
