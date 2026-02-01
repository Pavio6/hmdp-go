FROM golang:1.24-alpine AS build

ARG GOPROXY=https://goproxy.cn,direct
ENV GOPROXY=${GOPROXY}

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /bin/hmdp-server ./cmd/server

FROM alpine:3.20
RUN apk add --no-cache ca-certificates

WORKDIR /app
COPY --from=build /bin/hmdp-server /app/hmdp-server

ENV HMDP_CONFIG=/etc/hmdp/app.yaml
EXPOSE 8081
USER 65532:65532

ENTRYPOINT ["/app/hmdp-server"]
