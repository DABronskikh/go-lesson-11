FROM alpine:3.7

ADD bank /app/

ENTRYPOINT ["/app/main"]

EXPOSE 9999
