FROM golang:1.26-alpine AS build

WORKDIR /app

COPY services/users/go.mod services/users/go.sum ./
RUN go mod download

COPY services/users/ ./
RUN CGO_ENABLED=0 go build -o /users .

FROM alpine:3.19
RUN apk --no-cache add ca-certificates wget
WORKDIR /
COPY --from=build /users .
EXPOSE 8080
CMD ["./users"]
