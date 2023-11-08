FROM --platform=${BUILDPLATFORM} golang:1.21 as builder

WORKDIR /build

COPY deployer/* .

ENV CGO_ENABLED=0
ARG TARGETOS
ARG TARGETARCH

RUN  --mount=type=cache,target=/go/pkg/mod \
     --mount=type=cache,target=/root/.cache/go-build \
     GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -o deployer .

FROM alpine:3.18

RUN apk add --no-cache curl \
    && rm -rf /var/cache/apk/*

# Install Auth0 CLI
RUN curl -sSfL https://raw.githubusercontent.com/auth0/auth0-cli/main/install.sh | sh -s -- -b . \
    && mv ./auth0 /usr/local/bin

RUN apk del curl

COPY --from=builder /build/deployer /usr/bin/deployer

COPY entrypoint.sh /etc/entrypoint.sh
RUN chmod +x /etc/entrypoint.sh

ENTRYPOINT ["/etc/entrypoint.sh"]

