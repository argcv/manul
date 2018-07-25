FROM golang:1.9.3 as builder

ADD . /go/src/github.com/argcv/manul

RUN cd /go/src/github.com/argcv/manul && bash ./build.sh

FROM scratch

# x509: failed to load system roots and no roots provided
COPY --from=builder  /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

#COPY --from=builder /go/bin/sunlab-api /sunlab-api
COPY --from=builder /go/src/github.com/argcv/manul/manul /manul

COPY --from=builder /go/src/github.com/argcv/manul/manul-entrypoint /manul-entrypoint

EXPOSE 35000
EXPOSE 35001

ENTRYPOINT ["/manul"]
