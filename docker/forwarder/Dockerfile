FROM golang:1.14-alpine3.11 as build
COPY stage/ /stage
WORKDIR /stage/forwarder
RUN apk add git
RUN go build

FROM alpine:3.11
WORKDIR /relay
COPY --from=build /stage/forwarder/forwarder /relay/forwarder
ENTRYPOINT [/relay/forwarder]
