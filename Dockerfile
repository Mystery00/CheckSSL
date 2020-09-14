FROM golang:alpine3.12 as build
RUN mkdir /app
COPY . /go/src
WORKDIR /go/src
RUN go build -v -o /app/run CheckSSL

###
FROM scratch as final
COPY --from=build /app/run /
CMD ["/run"]