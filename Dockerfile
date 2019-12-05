ARG GO_IMAGE=docker.io/golang:1.13.1-stretch
ARG RELEASE_IMAGE=scratch
FROM ${GO_IMAGE} as builder

SHELL [ "/bin/bash", "-cex" ]
WORKDIR /usr/src/airshipui

# Take advantage of caching for dependency acquisition
COPY go.mod go.sum /usr/src/airshipui/
RUN go mod download

COPY . /usr/src/airshipui/
ARG MAKE_TARGET=build
RUN make ${MAKE_TARGET}

FROM ${RELEASE_IMAGE} as release
COPY --from=builder /usr/src/airshipui/bin/airshipui /usr/local/bin/airshipui
USER 65534
ENTRYPOINT [ "/usr/local/bin/airshipui" ]
