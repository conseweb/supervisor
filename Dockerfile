FROM ckeyer/obc:base

MAINTAINER Chuanjian Wang <me@ckeyer.com>

COPY . /go/src/github.com/conseweb/supervisor
WORKDIR /go/src/github.com/conseweb/supervisor
ENV PKG=github.com/conseweb/supervisor

RUN make build-local

EXPOSE 9376

CMD ["./bundles/supervisor", "node"]
