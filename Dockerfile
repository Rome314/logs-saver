FROM golang:1.17 AS builder
WORKDIR /app
COPY go.mod ./
RUN go mod download
#COPY build ./
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix nocgo -o /app cmd/main.go

FROM alpine
RUN apk add --no-cache tzdata
COPY --from=builder /app ./
ENTRYPOINT ["./app"]
