FROM golang:1.18-alpine As builder

WORKDIR /katch/

RUN apk update && apk add git upx

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -ldflags "-s -w" -o /usr/bin/katch ./cmd/

RUN upx -9 /usr/bin/katch

FROM chromedp/headless-shell

WORKDIR /katch/

COPY --from=builder /usr/bin/katch /usr/bin/katch

ENTRYPOINT []
