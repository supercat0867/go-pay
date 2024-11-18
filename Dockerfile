FROM swr.cn-north-4.myhuaweicloud.com/ddn-k8s/docker.io/golang:1.23-alpine

ENV GOPROXY=https://goproxy.cn,direct

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download -x

COPY . .

WORKDIR /app/cmd

RUN go build -o /app/main .

WORKDIR /app

EXPOSE 8080

CMD ["./main"]
