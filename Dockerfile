# Use Debian to obtain CA certificates.
FROM --platform=${BUILDPLATFORM} debian:trixie-slim AS certificates

# Install ca-certificates.
RUN apt-get update && apt-get install --yes ca-certificates

# Use Go for installing Ginkgo and building tests.
FROM --platform=${BUILDPLATFORM} golang:1.26 AS tests

ARG TARGETOS
ARG TARGETARCH

# Install Ginkgo for building tests on build platform.
RUN go install github.com/onsi/ginkgo/v2/ginkgo@latest

# Install Ginkgo for running tests on target platform.
RUN GOOS="${TARGETOS}" GOARCH="${TARGETARCH}" go install github.com/onsi/ginkgo/v2/ginkgo@latest

# Copy Ginkgo for running tests on target platform.
RUN [ "${TARGETOS}" = "$(go env GOOS)" ] && [ "${TARGETARCH}" = "$(go env GOARCH)" ] && \
    cp "/go/bin/ginkgo" /tmp/ginkgo || \
    cp "/go/bin/${TARGETOS}_${TARGETARCH}/ginkgo" /tmp/ginkgo

# Copy sources.
WORKDIR /app
COPY . .

# Build tests for target platform.
RUN GOOS="${TARGETOS}" GOARCH="${TARGETARCH}" ginkgo build --skip-package /X -r ./

# Use Debian for running tests.
FROM debian:trixie-slim

# Copy CA certificates.
COPY --from=certificates /etc/ssl/certs/ /etc/ssl/certs/
COPY --from=certificates /usr/share/ca-certificates/ /usr/share/ca-certificates/

# Copy Ginkgo.
COPY --from=tests /tmp/ginkgo /usr/local/bin/ginkgo

# Copy tests.
COPY --from=tests /app /app

# Define entrypoint.
WORKDIR /app
ENTRYPOINT [ "/app/entrypoint.sh" ]
