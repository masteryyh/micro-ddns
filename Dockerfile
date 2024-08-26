FROM ubuntu:24.04

LABEL author=masteryyh
LABEL email="yyh991013@163.com"

ARG USE_CN_MIRROR
ENV USE_CN_MIRROR=${USE_CN_MIRROR}

RUN if [ -n "$USE_CN_MIRROR" ]; then \
        sed -i 's/archive.ubuntu.com/mirrors.ustc.edu.cn/g' /etc/apt/sources.list.d/ubuntu.sources && \
        sed -i 's/security.ubuntu.com/mirrors.ustc.edu.cn/g' /etc/apt/sources.list.d/ubuntu.sources; \
    fi; \
    apt-get update && \
    apt-get install -y ca-certificates

WORKDIR /app

COPY build/micro-ddns .

ENTRYPOINT ["./micro-ddns"]
CMD ["run", "-c", "/etc/micro-ddns/config.yaml"]
