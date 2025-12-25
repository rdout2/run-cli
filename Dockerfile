FROM alpine:3

COPY rabdis /bin/run

RUN addgroup -g 1000 -S run && \
    adduser -u 1000 -S run -G run && \
    chown run:run /bin/run

USER run:run

ENTRYPOINT ["/bin/run"]
