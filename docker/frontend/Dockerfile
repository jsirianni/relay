FROM golang:1.14-alpine3.11 as build
COPY stage/ /stage
WORKDIR /stage/cmd/frontend
RUN apk add git
RUN go build

# use alpine instead of scratch / distroless in order to allow
# k8s to run "sleep" commands for lifecycle pre stop hooks
FROM alpine:3.11
WORKDIR /relay
COPY --from=build /stage/cmd/frontend/frontend /relay/frontend
ENTRYPOINT [/relay/frontend]
