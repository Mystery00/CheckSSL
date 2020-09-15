FROM golang:alpine3.12 as build
RUN mkdir /app
ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64
COPY . /go/src
WORKDIR /go/src
RUN go build -v -o /app/run CheckSSL

###
FROM scratch as final
COPY --from=build /app/run /
CMD ["/run"]