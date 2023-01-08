FROM harbor.adamkoro.com/bci/golang:1.19 as build
WORKDIR /go/src/build
COPY cmd/ ${WORKDIR}
RUN go mod download && go install github.com/swaggo/swag/cmd/swag@latest && mkdir /build
RUN swag init . && cp -r ./docs /build
RUN go build -ldflags="-s -w" -o /build/mq-producer

FROM harbor.adamkoro.com/bci/bci-micro:15.4
WORKDIR /home/user
COPY --from=build /build/mq-producer ${WORKDIR}
COPY --from=build /build/docs/swagger.json ${WORKDIR}/docs
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