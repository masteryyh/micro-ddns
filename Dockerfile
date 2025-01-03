FROM ubuntu:oracular

LABEL author=masteryyh
LABEL email="yyh991013@163.com"

RUN apt-get update && \
    apt-get install ca-certificates curl -y && \
    rm -rf /var/lib/apt/lists/* && \
    apt-get clean

COPY bin/micro-ddns-$TARGETARCH /usr/local/bin/micro-ddns

USER 1000:1000

HEALTHCHECK --interval=30s --timeout=5s --start-period=3s --retries=3 \
    CMD bash -c "curl -sf http://127.0.0.1:8080/ping > /dev/null || exit 1"

ENTRYPOINT ["micro-ddns"]
CMD ["run"]
