FROM --platform=$BUILDPLATFORM golang:1.24 AS builder

ARG BUILDPLATFORM
ARG TARGETPLATFORM
ARG TARGETOS
ARG TARGETARCH

COPY . /app
WORKDIR /app/
RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} && go build -ldflags "-s -w"

FROM --platform=$TARGETPLATFORM alpine:latest

COPY --from=builder /app/icon-metrics/ /app/
WORKDIR /app/
EXPOSE 8080

ENTRYPOINT [ "./icon-metrics" ]
