FROM alpine:3.12.0

RUN addgroup -S app && adduser -S -G app app
USER app

WORKDIR /app

ADD consumers/html/template.twig consumers/html/

# add binary
COPY build/linux/shopware-orders-scanner/ .
ENTRYPOINT ["/app/shopware-orders-scanner"]
