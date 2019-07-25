# Step 4: Delivery

1. Lets version our application: `git init && git add * && git commit -m "initial commit"`
    - Add `.gitignore` if needed;
2. Lets fix our dependencies:
    - For native go `GO111MODULE=on go mod init && go mod tidy`
    - For containerized go: `gontainer /bin/sh -c "GO111MODULE=on go mod init && go mod tidy"` 
3. We can build our project using i.e. `go build -o math`
    - Then delivery can look like moving compiled file to server. 
4. But we can use docker to unify application management.
    - Lets use dockerfile:
     ```dockerfile
     FROM golang:alpine
     WORKDIR /app
     COPY go.mod go.sum ./
     RUN apk add git
     RUN go mod download
     COPY .. .
     RUN go build -o main
     EXPOSE 3000
     EXPOSE 3001
     ENTRYPOINT ./math

     ```
   - We can build it: `docker build -t app_test .`
   - Then we can check it: `docker run -p 3000:3000 -p3001:3001 app_test server`
 5. Then we can separate building and running containers to performance purposes:
   - Lets use dockerfile:
   ```dockerfile
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
   ```
   - Use same `docker build` and `docker run` command for testing purposes.