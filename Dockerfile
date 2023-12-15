FROM golang:1.21-alpine as builder
RUN mkdir /build
WORKDIR /build

# manage app deps
COPY go.mod .
COPY go.sum .
RUN go mod download

# prepare base deps
COPY . .
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -a -installsuffix cgo -o restaurans-crawl-data .


FROM alpine:3.18.0
RUN apk add ca-certificates
WORKDIR /
COPY --from=builder /build/restaurans-crawl-data .

CMD [ "./restaurans-crawl-data" ]