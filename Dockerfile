FROM golang:1.18.3 AS base

ADD . /build

WORKDIR /build
RUN go build -o /scrubeer ./cmd/scrubeer

FROM scratch

COPY --from=base /scrubeer /scrubeer

WORKDIR /work

ENTRYPOINT [ "/scrubeer" ]
