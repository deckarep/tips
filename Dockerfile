FROM golang:1.21 AS build
WORKDIR /go/src/github.com/deckarep/tips
COPY . .
RUN go version
RUN VERSION=$(git describe --tags --abbrev=0) && \
CGO_ENABLED=0 go build -o bin/tips #-ldflags "-X="github.com/deckarep/tips/cmd.Version=${VERSION}

FROM alpine:3.16
RUN adduser -D tipsuser && apk add --no-cache bash git
COPY --from=build /go/src/github.com/deckarep/tips/bin/* /usr/bin/
USER tipsuser

RUN git config --global --add safe.directory '*'

ENTRYPOINT ["tips"]