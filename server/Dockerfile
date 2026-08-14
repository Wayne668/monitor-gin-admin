FROM golang:alpine AS builder

ARG APP=ginadmin
ARG VERSION=v10.1.0

# Install build dependencies (CGO for SQLite)
RUN apk add --no-cache gcc musl-dev sqlite-dev

ENV CGO_CFLAGS="-D_LARGEFILE64_SOURCE"
ENV CGO_ENABLED=1
ENV GOPROXY="https://goproxy.cn"

WORKDIR /go/src/${APP}

# Layer cache: download deps first
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -ldflags "-w -s -X main.VERSION=${VERSION}" -o ./${APP} .

# ---------------- runtime ----------------
FROM alpine

RUN apk add --no-cache ca-certificates tzdata sqlite-libs
ENV TZ=Asia/Shanghai

ARG APP=ginadmin
WORKDIR /go/src/${APP}
RUN mkdir -p data

COPY --from=builder /go/src/${APP}/${APP} /usr/bin/
COPY --from=builder /go/src/${APP}/configs /usr/bin/configs

ENTRYPOINT ["ginadmin", "start", "-d", "/usr/bin/configs", "-c", "prod"]
EXPOSE 8040
