FROM hashicorp/terraform:0.12.2 as builder
RUN mkdir -p /usr/src/infrastructure
WORKDIR /usr/src/infrastructure
RUN apk add -U make curl git gcc musl-dev go
COPY . /usr/src/terraform-provider-liquidweb
WORKDIR /usr/src/terraform-provider-liquidweb
RUN make build

FROM hashicorp/terraform:0.12.2
COPY --from=builder /usr/src/terraform-provider-liquidweb/terraform-provider-liquidweb /bin
