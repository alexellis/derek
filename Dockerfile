FROM alpine:3.5

RUN apk --no-cache add curl ca-certificates \ 
    && echo "Pulling watchdog binary from Github." \
    && curl -sSL https://github.com/alexellis/faas/releases/download/0.6.4/fwatchdog > /usr/bin/fwatchdog \
    && chmod +x /usr/bin/fwatchdog \
    && apk del curl --no-cache

WORKDIR /root/
COPY derek.pem  .
COPY derek      .
ENV cgi_headers="true"

ENV fprocess="./derek"
ENV secret_key="docker"
ENV installation=45362
ENV private_key="derek.pem"

ENV validate_hmac="true"

EXPOSE 8080
CMD ["fwatchdog"]
