# syntax=docker/dockerfile:1
FROM --platform=$BUILDPLATFORM node:26-alpine AS web
WORKDIR /src/web
COPY web/package.json web/package-lock.json ./
RUN --mount=type=cache,target=/root/.npm npm ci
COPY web/ ./
# The help shows the changelog (web/src/lib/help/markdown.ts).
COPY CHANGELOG.md /src/CHANGELOG.md
RUN npm run build

FROM --platform=$BUILDPLATFORM golang:1.27 AS build
ARG TARGETOS TARGETARCH
WORKDIR /src
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod go mod download
COPY . .
COPY --from=web /src/internal/webui/dist ./internal/webui/dist
ARG VERSION
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH go build -trimpath -ldflags="-s -w -X main.version=$VERSION" -o /standblatt ./cmd/standblatt

FROM gcr.io/distroless/static:nonroot
COPY --from=build /standblatt /standblatt
ENV STANDBLATT_LISTEN=:8080
EXPOSE 8080 9090
USER nonroot
ENTRYPOINT ["/standblatt"]
