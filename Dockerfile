FROM golang:1.12.11-alpine3.9 AS builder

RUN apk add --update --no-cache ca-certificates git

WORKDIR /service/build

ADD . /service/build

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o ctrlsrv /service/build/cmd/ctrl

FROM golang:1.12.11-alpine3.9

RUN apk --no-cache add ca-certificates git

WORKDIR /ctrl

RUN mkdir git
RUN mkdir blueprint
RUN mkdir actions
RUN mkdir bin

COPY --from=builder /service/build/ctrlsrv /ctrl/bin/

RUN chmod +x /ctrl/bin/ctrlsrv

EXPOSE 6060
EXPOSE 6063

ENTRYPOINT ["./bin/ctrlsrv", "-dir=/ctrl"]

