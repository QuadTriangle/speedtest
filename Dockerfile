FROM scratch

ARG TARGETARCH

COPY ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

COPY speedtest-linux-${TARGETARCH} /speedtest

ENTRYPOINT ["/speedtest"]
