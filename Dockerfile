FROM golang:1.23.4 as builder
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o api cmd/api/*.go

FROM scratch
WORKDIR /app
COPY --from=builder /app/api .
COPY --from=builder /app/.env .env
EXPOSE 8080
CMD ["./api"]
