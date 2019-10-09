FROM openfaas/classic-watchdog:0.18.1 as watchdog

FROM golang:1.11-alpine as build

ENV CGO_ENABLED=0

WORKDIR /go/src/github.com/alexellis/derek
COPY . .

RUN go test $(go list ./... | grep -v /vendor/) -cover

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o derek .

FROM alpine:3.10 as ship

COPY --from=watchdog /fwatchdog /usr/bin/fwatchdog
RUN chmod +x /usr/bin/fwatchdog

RUN apk --no-cache add ca-certificates

RUN addgroup -S app && adduser -S -g app app
RUN mkdir -p /home/app

WORKDIR /home/app

COPY --from=build /go/src/github.com/alexellis/derek/derek derek

RUN chown -R app /home/app

USER app

ENV cgi_headers="true"
ENV validate_hmac="true"
ENV validate_customers="true"

ENV combine_output="true"

ENV fprocess="./derek"

HEALTHCHECK --interval=5s CMD [ -e /tmp/.lock ] || exit 1

EXPOSE 8080
CMD ["fwatchdog"]
