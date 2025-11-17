FROM golang:1.21-alpine AS build
WORKDIR /app
COPY go.mod ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o gateway ./cmd/gateway

FROM alpine:3.20
WORKDIR /srv
COPY --from=build /app/gateway /srv/gateway
COPY api/swagger.json /srv/api/swagger.json
ENV PORT=8080
EXPOSE 8080
ENTRYPOINT ["/srv/gateway"]
