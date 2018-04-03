# FROM alpine
# LABEL maintainer="nicholas.morrow@ctgcompanies.com"
# RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*
FROM golang:1.8.3 as goget
ADD ./main.go  /go/src/github.com/vworri/CurrencryExchangeAPI/main.go
RUN go get .\..


FROM golang:1.8.3-alpine as gobuild
COPY --from=goget /go/src /go/src
RUN go install github.com/vworri/CurrencryExchangeAPI

FROM alpine
RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*
COPY --from=gobuild /go/bin /go/bin
RUN ls /go/bin
EXPOSE 80
CMD /go/bin/CurrencryExchangeAPI -port=80