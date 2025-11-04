FROM golang:1.25.3-alpine3.22 AS base
WORKDIR /app
COPY go.* ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /restful ./cmd/restful

FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=base /restful /restful
EXPOSE 8080
CMD ["/restful"]