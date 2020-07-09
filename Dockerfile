FROM golang:latest as build
WORKDIR /app
COPY service .
RUN go build -o bin/service -ldflags="-s -w" main.go

##

FROM ubuntu:latest
WORKDIR /app
ARG VERSION=0.12.6-1

RUN apt-get update && \
    apt-get install -y wget && \
    wget https://github.com/wkhtmltopdf/packaging/releases/download/${VERSION}/wkhtmltox_${VERSION}.focal_amd64.deb && \
    apt-get install -y --fix-broken ./wkhtmltox_${VERSION}.focal_amd64.deb && \
    rm ./wkhtmltox_${VERSION}.focal_amd64.deb

COPY --from=build /app/bin/service .
CMD ["./service"]

