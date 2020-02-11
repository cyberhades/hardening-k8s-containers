FROM golang:1.13 as build

WORKDIR /go/src/app
COPY main.go .
COPY go.mod .

RUN go get -d -v ./... && \
    CGO_ENABLED=0 go install -v ./...

FROM alpine

COPY --from=build /go/bin/notes /app
COPY assets /assets/
COPY pages /pages/
RUN mkdir /notes

EXPOSE 8080

CMD ["/app"]
