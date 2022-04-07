FROM alpine:3.15.4

RUN apk update
RUN apk add curl bash tzdata

COPY bin/linux/private-ghp /app/private-ghp
RUN chmod u+x /app/private-ghp

COPY entrypoint.sh /app/entrypoint.sh
RUN chmod u+x /app/entrypoint.sh

ENTRYPOINT [ "/app/entrypoint.sh" ]