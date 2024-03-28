FROM bitnami/golang:1.20 as builder
COPY / /app
WORKDIR /app
ENV GO111MODULE=on
ENV CGO_ENABLED=0
ENV GOSUMDB=off
ENV GOOS=linux
RUN go build -o /go/bin/main .

FROM alpine:latest
COPY --from=builder /go/bin/main .
EXPOSE 80
ENTRYPOINT ["/main"]
