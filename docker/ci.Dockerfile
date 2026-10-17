# CI only: packs bin/standblatt built and tested by the jobs binary and e2e (build context bin/).
# Local builds use docker/Dockerfile.
FROM gcr.io/distroless/static:nonroot
LABEL org.opencontainers.image.source=https://github.com/standblatt/server
COPY standblatt /standblatt
ENV STANDBLATT_LISTEN=:8080
EXPOSE 8080 9090
USER nonroot
ENTRYPOINT ["/standblatt"]
