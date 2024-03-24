# Build stage
FROM golang:1.22 AS builder
WORKDIR /app
COPY . .
RUN go build -o main main.go
RUN apt-get update && apt-get install -y wget && apt-get install -y netcat-openbsd
RUN wget -O migrate.tar.gz https://github.com/golang-migrate/migrate/releases/download/v4.14.1/migrate.linux-amd64.tar.gz \
    && tar -xzvf migrate.tar.gz \
    && rm migrate.tar.gz

# Run stage
FROM golang:1.22
WORKDIR /app
COPY --from=builder /app/main .
COPY --from=builder /app/migrate.linux-amd64 ./migrate
COPY app.env .
COPY start.sh .
COPY wait-for.sh .
COPY db/migration ./migration
RUN chmod +x /app/wait-for.sh
RUN chmod +x /app/start.sh

EXPOSE 8080
CMD [ "/app/main" ]
ENTRYPOINT [ "/app/start.sh" ]