FROM golang:alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN apk add git
RUN go mod download
COPY . .
RUN go build -o math

FROM scratch
WORKDIR /app
COPY --from=builder /app/math /app/
EXPOSE 3000
EXPOSE 3001
ENTRYPOINT ["./math"]
