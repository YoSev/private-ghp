FROM docker.io/library/golang:alpine3.17 AS build

WORKDIR /workspace
ADD . /workspace

RUN apk update && \
    apk add git make

RUN make build

FROM docker.io/library/alpine:3.17

RUN apk update && \
    apk add curl bash tzdata

COPY --from=build /workspace/bin/generic/private-ghp /app/private-ghp
RUN chmod u+x /app/private-ghp

ENV TZ=UTC
COPY entrypoint.sh /app/entrypoint.sh
RUN chmod u+x /app/entrypoint.sh

ENTRYPOINT [ "/app/entrypoint.sh" ]