FROM golang:alpine as builder
RUN apk add -U make curl git gcc musl-dev
RUN mkdir -p /usr/src/terraform-provider-liquidweb
COPY . /usr/src/terraform-provider-liquidweb
WORKDIR /usr/src/terraform-provider-liquidweb
RUN make build

FROM hashicorp/terraform:0.11.8
RUN mkdir -p /usr/src/infrastructure
WORKDIR /usr/src/infrastructure
COPY --from=builder /usr/src/terraform-provider-liquidweb/terraform-provider-liquidweb /bin/
RUN apk add -U jq && \
  curl -SsL https://github.com/kvz/json2hcl/releases/download/v0.0.6/json2hcl_v0.0.6_linux_amd64 -o /usr/local/bin/json2hcl && \
  chmod 755 /usr/local/bin/json2hcl && \
  json2hcl -version
