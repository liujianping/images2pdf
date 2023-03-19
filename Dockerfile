FROM golang:1.19-alpine AS builder
WORKDIR /app
COPY . .
RUN GOOS=linux go build -o main 

FROM  alpine
COPY --from=builder /app/main /usr/local/bin
ENTRYPOINT [ "main" ]
