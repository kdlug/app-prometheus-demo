FROM alpine:3.4

RUN apk -U add ca-certificates

EXPOSE 8080

ADD app-prometheus-demo /bin/app-prometheus-demo

CMD ["app-prometheus-demo"]
