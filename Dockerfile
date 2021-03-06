FROM golang:latest as build
WORKDIR /app
COPY src .
RUN go build -o bin/service -ldflags="-s -w" main.go

##

FROM ubuntu:focal
WORKDIR /app
ARG VERSION=0.12.6-1

RUN apt-get update && \
    apt-get install -y wget && \
    wget https://github.com/wkhtmltopdf/packaging/releases/download/${VERSION}/wkhtmltox_${VERSION}.focal_amd64.deb && \
    apt-get install --no-install-recommends --fix-broken -y ./wkhtmltox_${VERSION}.focal_amd64.deb && \
    rm ./wkhtmltox_${VERSION}.focal_amd64.deb && \
    rm -rf /var/lib/apt/lists/*

COPY --from=build /app/bin/service .
CMD ["./service"]
