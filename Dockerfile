FROM golang:1.20-alpine3.17 as base

WORKDIR /sala-2/
# COPY go.mod go.sum ./
COPY go.mod ./
COPY . .
RUN go build -o sala-2 ./cmd

FROM alpine:3.17 as binary
COPY --from=base /sala-2/ .
EXPOSE 3000
CMD ["./sala-2"]