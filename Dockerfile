FROM --platform=$BUILDPLATFORM golang:1.23 AS builder

COPY . /app
WORKDIR /app/
RUN echo $TARGETPLATFORM
RUN export CGO_ENABLED=0 && export GOARCH=$(echo "$TARGETPLATFORM" | cut -d "/" -f2) && go build -ldflags "-s -w"

FROM alpine:latest

COPY --from=builder /app/icon-metrics /app/
WORKDIR /app/
EXPOSE 8080

ENTRYPOINT [ "./icon-metrics" ]
