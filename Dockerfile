# First stage: complete build environment
FROM golang:1.12 AS builder-go
# workdir
WORKDIR /app
# add source code
ADD . .
# build
ENV GOPROXY="https://goproxy.io"
RUN go build main.go


# Second stage: minimal runtime environment
From alpine:latest

# apk mirrors
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories
# install base lib
RUN apk update \
    && apk upgrade \
    && apk add bash \
    && apk add ca-certificates \
    && apk add wget \
    && apk add iputils \
    && apk add iproute2 \
    && apk add libc6-compat \
    && apk add -U tzdata \
    #&& apk add curl \
    #&& apk add tcpdump \
    #&& apk add nghttp2-dev \
    #&& apk add s6 \
    && rm -rf /var/cache/apk/*

# timezone
RUN rm -rf /etc/localtime \
 && ln -s /usr/share/zoneinfo/Asia/Shanghai /etc/localtime

# workdir
WORKDIR /app
# copy  from the first stage
COPY --from=builder-go /app/main phonelocation
COPY --from=builder-go /app/config/config.toml config/config.toml
COPY --from=builder-go /app/data/ data/

ENV PHONE_DATA_DIR="./data"


# port
EXPOSE 8199

# run jar
CMD ["./phonelocation"]