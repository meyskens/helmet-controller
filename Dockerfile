FROM scratch
COPY helmet-controller /
ENTRYPOINT ["/helmet-controller"]