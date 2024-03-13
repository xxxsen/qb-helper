FROM golang:1.15.2-buster

COPY . /build 
WORKDIR /build 

RUN CGO_ENABLED=0 go build -a -tags netgo -ldflags '-w' -o qb-helper 

FROM alpine:3.12

COPY --from=0 /build/qb-helper /bin/

ENTRYPOINT ["/bin/qb-helper"]