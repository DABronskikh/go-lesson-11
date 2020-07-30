FROM alpine:3.7

ADD bank /

ENTRYPOINT ["/main"]

EXPOSE 9999
