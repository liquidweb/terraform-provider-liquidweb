FROM golang:1.13-alpine3.10 as dev
ENV GOCACHE /usr/src/terraform-provider-liquidweb/go/.cache
ENV GOPATH /usr/src/terraform-provider-liquidweb/go
ARG uid=1003
RUN apk add -U make curl git gcc musl-dev bind-tools bash terraform
RUN mkdir -p /usr/src/infrastructure
RUN mkdir -p /usr/src/terraform-provider-liquidweb
RUN adduser -h /usr/src/terraform-provider-liquidweb -g "" -D -u $uid builder
USER builder
WORKDIR /usr/src/terraform-provider-liquidweb
COPY --chown=builder:builder . .

FROM builder as builder
RUN make build

FROM hashicorp/terraform:0.12.2
COPY --from=builder /usr/src/terraform-provider-liquidweb/terraform-provider-liquidweb /bin
