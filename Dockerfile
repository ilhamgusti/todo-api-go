FROM golang:1.17 as builder
WORKDIR /go/src/github.com/ilhamgusti/todo-go-api/
ADD . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .

FROM alpine:latest
WORKDIR /home/ilhamgusti/
ADD .env ./
COPY --from=builder /go/src/github.com/ilhamgusti/todo-go-api/app ./
CMD ["./app"]  