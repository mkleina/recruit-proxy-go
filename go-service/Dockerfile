FROM golang:1.8

ARG BUILD_PACKAGES="git-core openssh-client g++ make libssl-dev curl"
ARG DEBIAN_FRONTEND=noninteractive

WORKDIR /go/src/github.com/go-service

#RUN adduser --disabled-password --system --uid 1000 deploy

# install build deps
RUN apt-get update && \
    apt-get dist-upgrade -q -y && \
    apt-get install -q -y --no-install-recommends $BUILD_PACKAGES

COPY . .

RUN make deps
RUN make

CMD ["./server", "-amqp_host", "rabbitmq"]

