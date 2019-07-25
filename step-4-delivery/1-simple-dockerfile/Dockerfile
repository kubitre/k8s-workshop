FROM golang:alpine
WORKDIR /app
COPY go.mod go.sum ./
RUN apk add git
RUN go mod download
COPY . .
RUN go build -o math
EXPOSE 3000
EXPOSE 3001
ENTRYPOINT ["./math"]
