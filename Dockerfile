FROM golang:1.14.3-stretch AS build
WORKDIR /src
COPY . .
RUN ls
RUN apt-get update
RUN apt-get install git
RUN go get github.com/gorilla/mux github.com/sirupsen/logrus gopkg.in/yaml.v2
RUN go build -o /out/server .
FROM scratch as bin
COPY --from=build /out/server /