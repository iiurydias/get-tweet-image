FROM ubuntu:latest

RUN apt-get update && apt-get install -y ca-certificates
WORKDIR /app
COPY ./bin/get-tweet-image /app/
COPY ./background.jpg /app/
COPY ./LilitaOne-Regular.ttf /app/
COPY ./Lora-Bold.ttf /app/
ENTRYPOINT ["/app/get-tweet-image"]
