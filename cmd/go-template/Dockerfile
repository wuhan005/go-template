FROM golang:1.24 as go_builder

WORKDIR /app
ARG GITHUB_SHA

ENV CGO_ENABLED=1

RUN go mod tidy
RUN go build -v -trimpath -ldflags "-w -s -extldflags '-static' -X 'github.com/wuhan005/go-template/internal/appconst.BuildCommit=$GITHUB_SHA'" -o go-template ./cmd/go-template

FROM ubuntu:24.10

RUN apt update && apt install -y ca-certificates tzdata && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
    && rm -rf /var/lib/apt/lists/* && update-ca-certificates

WORKDIR /home/app

COPY --from=go_builder /app/go-template .

RUN chmod 655 /home/app/go-template

ENTRYPOINT ["./go-template"]
EXPOSE 8000
