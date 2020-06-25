# Build the manager binary
FROM ubuntu:20.04 as base

RUN apt-get update && apt-get install tgt curl -y

ENV GOVERSION 1.13.12
ENV PATH $PATH:/usr/local/go/bin:/usr/local/kubebuilder/bin

RUN cd /tmp && curl -O https://dl.google.com/go/go${GOVERSION}.linux-amd64.tar.gz && \
    tar -C /usr/local -xzf go${GOVERSION}.linux-amd64.tar.gz
RUN os=$(go env GOOS) && \
    arch=$(go env GOARCH) && \
    curl -L https://go.kubebuilder.io/dl/2.3.1/${os}/${arch} | tar -xz -C /tmp/ && \
    mv /tmp/kubebuilder_2.3.1_${os}_${arch} /usr/local/kubebuilder


WORKDIR /workspace
# Copy the Go Modules manifests
COPY go.mod go.mod
COPY go.sum go.sum
# cache deps before building and copying source so that we don't need to re-download as much
# and so that source changes don't invalidate our downloaded layer
RUN go mod download

# Copy the go source
COPY cmd/ cmd/
COPY config/ config/
COPY api/ api/
COPY controllers/ controllers/
COPY utils/ utils/

FROM base as unit-test

ENV CGO_ENABLED=0
# RUN --mount=target=. \
#     --mount=type=cache,target=/root/.cache/go-build \
RUN go test -v ./...

FROM base as builder

# Build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on go build -a -o tgtd-operator cmd/tgtd-operator/main.go

FROM ubuntu:20.04 as bin
USER root

RUN apt-get update && apt-get install tgt curl -y

WORKDIR /
COPY --from=builder /workspace/tgtd-operator .

EXPOSE 3260

COPY ./entrypoint.sh /
ENTRYPOINT ["/entrypoint.sh"]
