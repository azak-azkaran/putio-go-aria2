FROM golang:alpine AS builder

RUN apk add --no-cache git

#ENV GO111MODULE="on go build"
RUN go get -u github.com/putdotio/go-putio
#RUN go get github.com/azak-azkaran/putio-go-aria2
WORKDIR /go/src/github.com/azak-azkaran/putio-go-aria2
COPY main.go ./main.go
COPY go.mod ./
COPY go.sum ./
COPY aria2downloader ./aria2downloader/
COPY organize ./organize/
COPY utils ./utils/

RUN go install


FROM alpine:latest
ENV ARIA2_ADDRESS="localhost:6800"
WORKDIR /app/
COPY --from=builder /go/bin/ ./
CMD ["./putio-go-aria2"]
