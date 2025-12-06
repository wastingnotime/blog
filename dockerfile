# syntax=docker/dockerfile:1
FROM golang:1.25-alpine AS build-stage

WORKDIR /app

COPY go.mod /app

RUN go mod download

COPY . /app

RUN go run ./cmd/migrate check

RUN go build ./cmd/build

ENV IS_PROD_ENV=true
ENV SITE_BASE_URL=https://wastingnotime.org/

RUN ./build

FROM nginx:1.29.3-alpine AS deploy-stage

# act as doc only
EXPOSE 80
LABEL vendor=wastingnotime.org

COPY --from=build-stage /app/public /usr/share/nginx/html
