FROM golang:alpine AS builder
RUN cd /src && go build -o math

FROM scratch
WORKDIR /app
COPY --from=builder /src/math /app/
ENTRYPOINT ./math