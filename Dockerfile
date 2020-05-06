FROM golang:alpine AS builder

RUN apk add --no-cache git

#ENV GO111MODULE="on go build"
RUN go get -u github.com/putdotio/go-putio
RUN go get github.com/azak-azkaran/putio-go-aria2
WORKDIR /go/src/github.com/azak-azkaran/putio-go-aria2/
RUN go install


FROM alpine:latest
WORKDIR /app/
COPY --from=builder /go/bin/ ./
CMD ["./putio-go-aria2"]
