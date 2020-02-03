FROM alpine:3.10

RUN apk add --no-cache \
  bash curl nano

WORKDIR /app
COPY app ./

ENTRYPOINT [ "/app/entrypoint.sh" ]
