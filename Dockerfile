FROM golang:1.11-alpine as builder

WORKDIR /go/src/github.com/reo7sp/two-step-migrate

RUN apk add --no-cache git
RUN go get github.com/golang/dep/cmd/dep

COPY Gopkg.toml .
COPY Gopkg.lock .
RUN dep ensure -vendor-only

COPY . .
RUN go build


FROM alpine

WORKDIR /usr/bin
COPY --from=builder /go/src/github.com/reo7sp/two-step-migrate/two-step-migrate two-step-migrate

CMD two-step-migrate --help
