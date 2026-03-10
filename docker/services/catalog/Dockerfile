FROM golang:1.26-alpine AS build

WORKDIR /app

COPY services/catalog/go.mod services/catalog/go.sum ./
RUN go mod download

COPY services/catalog/ ./
ARG VERSION=dev
RUN CGO_ENABLED=0 go build -ldflags "-X github.com/ecommerce/services/catalog/internal/handlers.Version=${VERSION}" -o /catalog .

FROM alpine:3.19
RUN apk --no-cache add ca-certificates wget
WORKDIR /
COPY --from=build /catalog .
EXPOSE 8081
CMD ["./catalog"]
