FROM golang:1.14 as builder
WORKDIR /workspace
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on go build ./cmd/nsinjector

FROM gcr.io/distroless/static:nonroot
WORKDIR /
COPY --from=builder /workspace .
USER nonroot:nonroot

ENTRYPOINT ["/nsinjector"]