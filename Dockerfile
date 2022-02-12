FROM golang:1.17-alpine as builder

ENV GOPROXY=https://goproxy.cn

WORKDIR /app
COPY . .

RUN go mod download && go build -v -ldflags "-s -w" -o go2rss
EXPOSE 8081

CMD [ "/app/go2rss" ]
