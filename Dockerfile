FROM golang:1.17 as build

# These are used by the Makefile
# The `make build` directive uses these variables
# They will be baked into the build binary 
ARG BUILD_HASH=nil
ARG BUILD_VERSION=docker-unknown

RUN mkdir /go/src/app-build

COPY cmd /go/src/app-build/cmd
COPY internal /go/src/app-build/internal
COPY pkg /go/src/app-build/pkg
COPY test /go/src/app-build/test
COPY go.mod /go/src/app-build/go.mod
COPY go.sum /go/src/app-build/go.sum
COPY .golangci.yaml /go/src/app-build/.golangci.yaml
COPY Makefile /go/src/app-build/Makefile

WORKDIR /go/src/app-build/

RUN make build

## --- CI --- #
FROM build as ci

RUN go get github.com/golangci/golangci-lint/cmd/golangci-lint@v1.38.0

# --- FINAL --- #
FROM golang:1.17 as final

COPY --from=build /go/bin/ymir /bin/ymir
RUN chmod +x /bin/ymir

ENTRYPOINT ["/bin/ymir"]

# --- HOT-RELOAD --- #
FROM ci as hot-reload

RUN go get github.com/cespare/reflex

RUN mkdir /go/src/app-dev
VOLUME /go/src/app-dev

WORKDIR /go/src/app-dev

COPY .docker/ymir/hot-reload.sh /hot-reload.sh
RUN chmod +x /hot-reload.sh

ENTRYPOINT /hot-reload.sh