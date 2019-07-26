FROM golang:alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN apk add git
RUN go mod download
ADD . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o math main.go
RUN chmod +x math

FROM scratch
WORKDIR /app
COPY --from=builder /app/math /app/
EXPOSE 3000
EXPOSE 3001
ENTRYPOINT ["./math"]
CMD ["server"]
