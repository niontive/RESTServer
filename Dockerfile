FROM golang:1.13 AS build
WORKDIR /src
RUN apt-get update
RUN apt-get install git
COPY go.* .
RUN go mod download
COPY . .
RUN go build -o /out/server .
FROM scratch as bin
COPY --from=build /out/server /