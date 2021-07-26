FROM golang:1.13 AS base
WORKDIR /src
RUN apt-get update
RUN apt-get install git
COPY go.* .
RUN go mod download
COPY . .

FROM base as build
RUN go build -o /out/server .

FROM base as unit-test
RUN go test -v .

FROM scratch as bin
COPY --from=build /out/server /