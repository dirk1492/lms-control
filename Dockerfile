FROM golang:alpine

RUN apk update && apk add gcc musl-dev upx ca-certificates dep git

COPY control/*.go /go/src/github.com/dirk1492/lms-control/control/
COPY *.go /go/src/github.com/dirk1492/lms-control/
COPY Gopkg.* /go/src/github.com/dirk1492/lms-control/

WORKDIR /go/src/github.com/dirk1492/lms-control

RUN dep ensure
RUN go build -ldflags "-linkmode external -extldflags -static -s -w" -o lms-control
RUN upx lms-control

ENV TIMETABLE=
ENV LMS_SERVER=
ENV LMS_PORT=
ENV INTERVAL= 

FROM scratch
COPY --from=0 /go/src/github.com/dirk1492/lms-control/lms-control /lms-control
COPY --from=0 /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=0 /usr/local/go/lib/time/zoneinfo.zip /
ENV TZ=Pacific/Auckland
ENV ZONEINFO=/zoneinfo.zip
COPY passwd /etc/passwd
USER nobody
CMD ["/lms-control"]