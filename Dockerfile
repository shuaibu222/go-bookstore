FROM alpine:latest

RUN mkdir /app
WORKDIR /app

COPY bookstore /app
COPY .env /app

CMD [ "/app/bookstore"]