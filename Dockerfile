FROM golang:1.26 AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY api ./api
COPY cmd ./cmd
COPY internal ./internal
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/server ./cmd/api

FROM debian:bookworm-slim
WORKDIR /app
COPY --from=build /out/server /app/server
COPY ui/html /app/ui/html
COPY ui/static/css-GIS /app/ui/static/css-GIS
COPY ui/static/js-GIS /app/ui/static/js-GIS
COPY ui/static/json-GIS /app/ui/static/json-GIS
COPY ui/static/assets /app/assets-base
ENV GIN_MODE=release
CMD ["sh", "-c", "mkdir -p /app/ui/static/assets; cp -rn /app/assets-base/. /app/ui/static/assets/; exec /app/server"]
