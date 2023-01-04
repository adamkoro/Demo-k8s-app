FROM registry.suse.com/bci/golang:1.19 as build
WORKDIR /go/src/build
COPY app/ ${WORKDIR}
RUN go mod download
RUN go build -ldflags="-s -w" -o /mq-producer

FROM registry.suse.com/bci/bci-micro:15.4
WORKDIR /home/user
COPY --from=build /mq-producer ${WORKDIR}
ENV GIN_MODE=release
RUN echo "user:x:10000:10000:user:/home/user:/bin/bash" >> /etc/passwd && chown user /home/user/mq-producer && chmod u+x /home/user/mq-producer
EXPOSE 8080
USER user
ENTRYPOINT ["./mq-producer"]