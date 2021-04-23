FROM golang:1.16.0-alpine3.9
WORKDIR /app
ADD . .
RUN go build -o main .
ENTRYPOINT ["/app/main"]
CMD ["serve"]
