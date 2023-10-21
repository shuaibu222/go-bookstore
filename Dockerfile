FROM golang:1.20-alpine as builder

RUN mkdir /app
WORKDIR /app

# Copy project dependencies and build the application
COPY go.mod go.sum ./
COPY . .
RUN go mod download

RUN CGO_ENABLED=0 go build -o bookstore .

# build tiny docker image
FROM alpine:latest

RUN mkdir /app
WORKDIR /app

COPY --from=builder /app/bookstore /app
COPY .env /app

# EXPOSE the port
EXPOSE 8000

CMD ["/app/bookstore"]


