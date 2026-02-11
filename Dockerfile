FROM scratch
ARG EXT=
COPY go-filterweb${EXT} /
ENTRYPOINT ["/go-filterweb${EXT}"]
