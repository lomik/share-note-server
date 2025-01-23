FROM --platform=${BUILDPLATFORM} golang:alpine as compiler
ARG TARGETOS
ARG TARGETARCH
ENV CGO_ENABLED=0

WORKDIR /go/src/share-note-server
ADD go.mod .
ADD go.sum .
RUN go mod download
COPY . .

RUN GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -ldflags="-s -w" -o share-note-server cmd/share-note-server/main.go

FROM --platform=${TARGETPLATFORM} alpine
WORKDIR /
COPY --from=compiler /go/src/share-note-server/share-note-server /share-note-server
COPY config.yaml /config.yaml


ENTRYPOINT ["/share-note-server"]