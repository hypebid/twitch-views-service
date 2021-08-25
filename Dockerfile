FROM golang:1.16.5-alpine3.13
RUN mkdir /app
COPY . /app
WORKDIR /app
RUN go build -o bin/twitch-views -v .
EXPOSE 8880
CMD [ "/app/bin/twitch-views" ]