FROM golang:1.15

COPY output/main-linux-amd64 /app/main
COPY config/config.json /app/config/config.json

EXPOSE 8080

CMD ["/app/main"]