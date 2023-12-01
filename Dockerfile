# syntax=docker/dockerfile:1

FROM golang:1.21 AS build-stage

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . ./

RUN CGO_ENABLED=0 GOOS=linux go build -o bin/

FROM build-stage AS run-test-stage
RUN go test -v ./...

FROM gcr.io/distroless/base-debian11 AS build-release-stage

WORKDIR /

COPY --from=build-stage /app/bin/gotab /gotab

EXPOSE 8080

USER nonroot:nonroot
ENV GIN_MODE=release
ENTRYPOINT ["/gotab"]