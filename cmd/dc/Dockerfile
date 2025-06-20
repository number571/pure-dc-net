FROM --platform=linux/amd64 ubuntu:20.04

RUN apt-get update && apt-get install -y wget gcc
RUN wget https://go.dev/dl/go1.23.0.linux-amd64.tar.gz && \ 
    tar -C /opt -xzf go1.23.0.linux-amd64.tar.gz

WORKDIR /pure-dc-net
ENV PATH="${PATH}:/opt/go/bin"
COPY ./ ./
RUN go build -o dc ./cmd/dc

ENV SERVICE_PATH="."
ENV SERVICE_NAME="_"
ENV CONSUME_ADDR=":8080"
ENV PRODUCE_ADDR=":8081"
CMD ./dc 
