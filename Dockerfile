FROM alpine:3.12.0

RUN addgroup -S app && adduser -S -G app app
USER app

WORKDIR /app

ADD render/template.twig render/

# add binary
COPY build/linux/shopware-orders-scanner/ .
ENTRYPOINT ["/app/shopware-orders-scanner"]
