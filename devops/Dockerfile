FROM golang:1.14.2 as compiler
WORKDIR /app
COPY ./src/go.mod ./src/go.sum ./
RUN go mod download
COPY ./src .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build hive

FROM alpine
WORKDIR /app
COPY --from=compiler /app/hive /app/hive
COPY --from=compiler /app/migrations /app/migrations
ENV PATH="/app:${PATH}"