FROM alpine:3.17.0
RUN apk --no-cache add ca-certificates git
COPY aiac /usr/local/bin/aiac
ENTRYPOINT ["aiac"]
