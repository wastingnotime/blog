# syntax=docker/dockerfile:1.6

FROM golang:1.25-alpine AS build-stage

WORKDIR /app

ARG TARGETOS
ARG TARGETARCH

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH \
    go build -o analytics-api ./api/analytics

FROM alpine:3.22 AS deploy-stage

ENV ENVIRONMENT=production

EXPOSE 8080
LABEL vendor=wastingnotime.org

RUN addgroup -S nonroot && adduser -S appuser -G nonroot

WORKDIR /app
RUN chown appuser .

COPY --from=build-stage /app/analytics-api .

USER appuser

ENTRYPOINT ["./analytics-api"]
