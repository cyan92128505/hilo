FROM golang:1.25.0-bullseye AS base
ENV GO111MODULE=on
ADD ./go.mod ./go.sum /workspace/

WORKDIR /workspace
RUN go mod download

FROM base AS build
ARG CMD_FOLDER
ADD . /workspace
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /go/bin/main ./cmd/$CMD_FOLDER

FROM gcr.io/distroless/static-debian11 AS runtime
COPY --from=build /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /etc/passwd /etc/passwd
COPY --from=build /etc/group /etc/group
COPY --from=build /go/bin/main /

CMD ["/main"]
