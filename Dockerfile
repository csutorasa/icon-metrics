FROM --platform=$BUILDPLATFORM golang:1.22 as builder

COPY . /app
WORKDIR /app/
RUN export CGO_ENABLED=0 && go build -ldflags "-s -w"

FROM alpine:latest

COPY --from=builder /app/icon-metrics /app/
WORKDIR /app/
EXPOSE 8080

ENTRYPOINT [ "./icon-metrics" ]
