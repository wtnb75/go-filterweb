FROM scratch
ARG EXT=
COPY filterweb${EXT} /
ENTRYPOINT ["/filterweb${EXT}"]
