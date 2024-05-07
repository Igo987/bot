# trunk-ignore-all(trivy/DS026)
# trunk-ignore-all(trivy/DS002)
FROM golang:1.22-alpine as builder
WORKDIR /build
# trunk-ignore(hadolint/DL3018)
COPY ./config/config.yaml ./config/config.yaml
RUN apk add --no-cache git bash git musl-dev
COPY . .
RUN go mod download
COPY . .
RUN go build -o ./bin/cmd/app cmd/app/main.go

# trunk-ignore(trivy/DS001)
FROM alpine:latest
WORKDIR /app

COPY --from=builder /build/config/config.go /config/config.go
COPY --from=builder /build/config/config.yaml /config/config.yaml
COPY --from=builder /build/bin/cmd/app /cmd/app
CMD [ "/cmd/app" ]


