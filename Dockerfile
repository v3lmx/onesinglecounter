FROM alpine:latest
COPY build/counter /counter
ENTRYPOINT ["/counter"]
