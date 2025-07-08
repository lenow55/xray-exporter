FROM scratch

ARG ARCH

EXPOSE 9550

COPY dist/xray-exporter-linux-${ARCH} /usr/bin/xray-exporter

ENTRYPOINT [ "/usr/bin/xray-exporter" ]
