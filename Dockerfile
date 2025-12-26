FROM scratch
COPY airgradient-exporter /
ENTRYPOINT ["/airgradient-exporter"]
