# syntax=docker/dockerfile:1

FROM golang:1.26.5-trixie AS build-stage

WORKDIR /app

COPY go.mod go.sum ./
COPY cmd/ ./cmd/
COPY internal/ ./internal/

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o /moviedb-api ./cmd/


FROM gcr.io/distroless/static-debian13 AS run-stage

COPY --from=build-stage /moviedb-api /moviedb-api

EXPOSE 8080/tcp

USER nonroot:nonroot

CMD ["/moviedb-api"]
