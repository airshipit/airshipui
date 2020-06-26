ARG GO_IMAGE=docker.io/golang:1.14.4-stretch
ARG RELEASE_IMAGE=scratch
FROM ${GO_IMAGE} as builder

# Inject custom root certificate authorities if needed
# Docker does not have a good conditional copy statement and requires that a source file exists
# to complete the copy function without error.  Therefore the README.md file will be copied to
# the image every time even if there are no .crt files.
COPY ./certs/* /usr/local/share/ca-certificates/
RUN update-ca-certificates

SHELL [ "/bin/bash", "-cex" ]
WORKDIR /usr/src/airshipui

# Take advantage of caching for dependency acquisition
COPY go.mod go.sum /usr/src/airshipui/
RUN go mod download

COPY . /usr/src/airshipui/
ARG MAKE_TARGET=build
RUN for target in $MAKE_TARGET; do make $target; done

FROM ${RELEASE_IMAGE} as release
COPY --from=builder /usr/src/airshipui/bin/airshipui /usr/local/bin/airshipui
USER 65534
ENTRYPOINT [ "/usr/local/bin/airshipui" ]