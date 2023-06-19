# ==========================================
# 1st Stage
# ==========================================
FROM golang:1.18 AS builder

## Set the working directory
WORKDIR /app

## Copy source
COPY . .

## Compile
RUN make build

# ==========================================
# 2nd Stage
# ==========================================
FROM alpine:latest

ENV APP_NAME=tinder-like-app

WORKDIR /app

## Add ssl cert
RUN apk add --update --no-cache ca-certificates

## Add timezone data
RUN apk --no-cache add tzdata

## Copy binary file from 1st stage
COPY --from=builder /app/bin/* ./

# Set environment variables for the PostgreSQL database
ENV POSTGRES_HOST=db
ENV POSTGRES_PORT=5432
ENV POSTGRES_USER=root
ENV POSTGRES_PASSWORD=password123
ENV POSTGRES_DB=postgres

# Install PostgreSQL client
RUN apk update && apk add postgresql-client

## Copy migration files
COPY ./database ./database

CMD ["./tinder-like-app", "server"]
