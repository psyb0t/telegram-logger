FROM golang:1.17 as build

RUN mkdir -p /app/build/
WORKDIR /app/

ENV GO111MODULE=on

COPY . .

RUN go mod download

RUN CGO_ENABLED=0 go build -a -installsuffix cgo -o /app/build/app cmd/*.go

FROM alpine:3.16

RUN apk --no-cache add ca-certificates

COPY --from=build /app/build/app /root/

EXPOSE 80

WORKDIR /root/

CMD ["./app"]
