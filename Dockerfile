FROM alpine:latest

RUN apk update
RUN apk add netcat-openbsd
COPY ./nc.sh /
RUN chmod +x /nc.sh
ENTRYPOINT ["/bin/sh", "/nc.sh"]
