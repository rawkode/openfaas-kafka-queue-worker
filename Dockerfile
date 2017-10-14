FROM golang:1.8 AS builder

RUN wget -qO - http://packages.confluent.io/deb/3.3/archive.key | apt-key add -
RUN echo "deb [arch=amd64] http://packages.confluent.io/deb/3.3 stable main" >> /etc/apt/sources.list

RUN DEBIAN_FRONTEND=noninteractive apt update \
    && apt install -y librdkafka-dev

WORKDIR /go/src/github.com/openfaas/faas-kafka-queue-worker
COPY main.go .

RUN go get
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o openfaas_kafka_queue_worker .

# FROM alpine:3.5

# RUN apk add --no-cache ca-certificates

# EXPOSE 8080
# ENV http_proxy ""
# ENV https_proxy ""

# COPY --from=builder /go/src/github.com/openfaas/faas-kafka-queue-worker/openfaas_kafka_queue_worker /

CMD ["/go/src/github.com/openfaas/faas-kafka-queue-worker/openfaas_kafka_queue_worker"]
