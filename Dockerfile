FROM golang:1.16 as builder

WORKDIR /workspace
COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod download

# Copy the go source
COPY main.go main.go
COPY pkg/ pkg/
COPY cmd/ cmd/

# Build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o blade main.go

FROM hextechs/alpine:3.13.5
WORKDIR /
COPY --from=builder /workspace/blade .
USER 65532:65532

ENTRYPOINT ["tini", "--", "/blade"]
