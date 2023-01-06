FROM golang:1.19 as build

RUN mkdir -p /app/build/
WORKDIR /app/

COPY . .

RUN go mod download

RUN GO111MODULE=on CGO_ENABLED=0 go build -a -installsuffix cgo -o /app/build/app cmd/*.go

FROM alpine:3.17

RUN apk --no-cache add ca-certificates

COPY --from=build /app/build/app /root/

EXPOSE 80

WORKDIR /root/

CMD ["./app"]
