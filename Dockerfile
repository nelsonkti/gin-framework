ARG DOCKER_USERNAME
FROM $DOCKER_USERNAME:alpine
#noroot
RUN addgroup -S nonroot && adduser -u 65530 -S nonroot -G nonroot
USER 65530
ARG service
ARG typ

WORKDIR /app

COPY main/main /app/main