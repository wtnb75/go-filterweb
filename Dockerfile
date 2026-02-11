FROM scratch
ARG EXT=
COPY go-filterweb${EXT} /go-filterweb
ENTRYPOINT ["/go-filterweb"]
