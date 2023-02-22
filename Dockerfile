FROM --platform=${TARGETPLATFORM:-linux/amd64} ghcr.io/openfaas/classic-watchdog:0.2.2 as watchdog
FROM --platform=${BUILDPLATFORM:-linux/amd64} golang:1.18-alpine3.17 as build

ENV CGO_ENABLED=0
ENV GO111MODULE=on

WORKDIR /go/src/github.com/alexellis/derek
COPY . .

RUN GOOS=${TARGETOS} GOARCH=${TARGETARCH} CGO_ENABLED=${CGO_ENABLED} go test $(go list ./... | grep -v /vendor/) -cover
RUN GOOS=${TARGETOS} GOARCH=${TARGETARCH} CGO_ENABLED=${CGO_ENABLED} go build -mod=vendor -o derek .

FROM --platform=${TARGETPLATFORM:-linux/amd64} alpine:3.17.2 as ship

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
