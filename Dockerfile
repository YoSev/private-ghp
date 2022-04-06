FROM alpine:3.15.4

RUN apk update
RUN apk add curl bash

COPY bin/linux/private-ghp /app/private-ghp
RUN chmod u+x /app/private-ghp

ENTRYPOINT [ "/app/private-ghp" ]