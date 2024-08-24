FROM ubuntu:24.04

LABEL author=masteryyh
LABEL email="yyh991013@163.com"

WORKDIR /app

COPY build/micro-ddns .

ENTRYPOINT ["./micro-ddns"]
CMD ["run", "-c", "/etc/micro-ddns/config.yaml"]
