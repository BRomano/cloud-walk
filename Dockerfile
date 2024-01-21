ARG GO_VERSION=1.21.6

FROM golang:${GO_VERSION} as builder
ARG MAIN_PATH="./cmd/cloud-walk"

WORKDIR /build
ADD . .
RUN GIT_COMMIT=$(git show --format="%h" --no-patch) && \
    GIT_TAG=$(echo 'v0.0') && \
    BUILD_TIMESTAMP=$(date +%Y-%m-%dT%H:%M:%S%z) && \
    GO_LDFLAGS="-w -X main.Version=$GIT_TAG -X main.Commit=$GIT_COMMIT  -X main.Timestamp=$BUILD_TIMESTAMP" && \
    set -x && \
    CGO_ENABLED=0 go build --mod=vendor -ldflags="$GO_LDFLAGS"-a -installsuffix cgo -o app ${MAIN_PATH}

FROM golang:${GO_VERSION}-alpine as app
RUN apk update

WORKDIR /app
COPY --from=builder /build/app .
ENTRYPOINT ./app