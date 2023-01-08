FROM harbor.adamkoro.com/bci/golang:1.19 as build
WORKDIR /go/src/build
COPY cmd/ ${WORKDIR}
RUN go mod download
RUN go build -ldflags="-s -w" -o /mq-producer

FROM harbor.adamkoro.com/bci/bci-micro:15.4
WORKDIR /home/user
COPY --from=build /mq-producer ${WORKDIR}
ENV GIN_MODE=release \
    USERNAME=guest \
    PASSOWRD=guest \
    MQHOST=127.0.0.1 \
    PORT=5672 \
    VHOST=/ \
    HTTP_PORT=8080
RUN echo "user:x:10000:10000:user:/home/user:/bin/bash" >> /etc/passwd && chown user /home/user/mq-producer && chmod u+x /home/user/mq-producer
EXPOSE ${HTTP_PORT}
USER user
ENTRYPOINT ["./mq-producer"]