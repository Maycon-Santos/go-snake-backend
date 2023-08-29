FROM golang:1.21 AS base
WORKDIR /usr/src/app

COPY . .

RUN go mod vendor

FROM base AS development
RUN go install github.com/cortesi/modd/cmd/modd@latest
RUN go install github.com/golang/mock/mockgen@v1.6.0
RUN go mod tidy

FROM base AS production
