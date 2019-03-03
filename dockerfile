FROM alpine

RUN apk add ca-certificates

COPY vms /

EXPOSE 3500

ENTRYPOINT ["/vms"]