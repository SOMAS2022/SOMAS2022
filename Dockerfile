# syntax=docker/dockerfile:1

## Build
FROM golang:1.19.3-alpine AS build

WORKDIR /app

COPY ./pkg/infra ./

RUN go build -o /main.out main.go

## Deploy
FROM alpine

WORKDIR /somas

COPY --from=build /main.out ./game

ENTRYPOINT ["/somas/game"]
