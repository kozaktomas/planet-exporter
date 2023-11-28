FROM golang:1.21 AS builder
ENV CGO_ENABLED 0
ADD . /app
WORKDIR /app
RUN go build -ldflags "-s -w" -v -o planet-exporter .

FROM alpine:3
RUN apk update && \
    apk add openssl tzdata && \
    rm -rf /var/cache/apk/* \
    && mkdir /app

WORKDIR /app
ENV VSOP87 /app/data
ADD https://raw.githubusercontent.com/ctdk/vsop87/master/VSOP87B.ear /app/data/VSOP87B.ear
ADD https://raw.githubusercontent.com/ctdk/vsop87/master/VSOP87B.jup /app/data/VSOP87B.jup
ADD https://raw.githubusercontent.com/ctdk/vsop87/master/VSOP87B.mar /app/data/VSOP87B.mar
ADD https://raw.githubusercontent.com/ctdk/vsop87/master/VSOP87B.mer /app/data/VSOP87B.mer
ADD https://raw.githubusercontent.com/ctdk/vsop87/master/VSOP87B.nep /app/data/VSOP87B.nep
ADD https://raw.githubusercontent.com/ctdk/vsop87/master/VSOP87B.sat /app/data/VSOP87B.sat
ADD https://raw.githubusercontent.com/ctdk/vsop87/master/VSOP87B.ura /app/data/VSOP87B.ura
ADD https://raw.githubusercontent.com/ctdk/vsop87/master/VSOP87B.ven /app/data/VSOP87B.ven
ADD Dockerfile /Dockerfile
COPY --from=builder /app/planet-exporter /app/planet-exporter

RUN chown -R nobody /app \
    && chmod 500 /app/planet-exporter \
    && chmod -R 700 /app/data

USER nobody
ENTRYPOINT ["/app/planet-exporter"]
