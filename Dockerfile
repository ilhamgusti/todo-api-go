FROM golang:1.17 as builder
WORKDIR /go/src/github.com/ilhamgusti/todo-go-api/
ADD . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .

FROM alpine:latest
WORKDIR /home/ilhamgusti/
ADD .env ./
COPY --from=builder /go/src/github.com/ilhamgusti/todo-go-api/app ./

## https://blog.phusion.nl/2015/01/20/docker-and-the-pid-1-zombie-reaping-problem/
# RUN apk add --no-cache dumb-init
# ENTRYPOINT ["/usr/bin/dumb-init", "--"]

CMD ./app