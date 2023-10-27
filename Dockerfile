FROM alpine:latest

RUN mkdir /app
WORKDIR /app

COPY bookstore /app
COPY .env /app

CMD [ "/app/bookstore"]


# FROM golang:1.20-alpine as builder

# RUN mkdir /app

# COPY . /app

# WORKDIR /app

# RUN CGO_ENABLED=0 go build -o bookstore  .

# # commented above to simplify docker image but i can leave it

# # build tiny docker image
# FROM alpine:latest

# RUN mkdir /app

# WORKDIR /app

# COPY --from=builder /app/bookstore /app
# COPY .env /app

# CMD [ "/app/bookstore" ]