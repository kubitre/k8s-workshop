FROM golang:alpine
WORKDIR /app
COPY go.mod go.sum ./
RUN apk add git
RUN go mod download
ADD . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o math main.go
RUN chmod +x math
EXPOSE 3000
EXPOSE 3001
RUN ls -lah
ENTRYPOINT ["./math"]
CMD ["server"]
