FROM golang:1.18 AS build-stage

ARG BUILD_REF
ARG BUILD_DATE

WORKDIR /code

COPY . .

RUN go build -ldflags "-X main.buildVersion=${BUILD_VERSION} -X main.buildCommit=${BUILD_COMMIT} -X main.buildCommit=${BUILD_DATE}"

FROM archlinux:latest AS lux-stage

COPY --from=build-stage /code/barber /linux/amd64
COPY --from=build-stage /code/config/config.yaml /linux/config/config.yaml

WORKDIR /linux

RUN ./amd64

LABEL org.opencontainers.image.created="${BUILD_DATE}" \
      org.opencontainers.image.title="shabo" \
      org.opencontainers.image.authors="Filip Boye-Kofi <filip.boyekofi@gmail.com>" \
      org.opencontainers.image.source="https://github.com/spuxy/shabo/.docker/Dockerfile" \
      org.opencontainers.image.revision="${BUILD_VERSION}" \
      org.opencontainers.image.vendor="Spuxy"