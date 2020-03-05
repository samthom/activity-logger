FROM golang:1.13.8-alpine3.11
RUN apk update && apk add git && apk add curl
RUN curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
# RUN go get github.com/gorilla/mux
# RUN go get go.mongodb.org/mongo-driver/mongo
ADD . /go/src/github.com/samthom/activity-logger
WORKDIR /go/src/github.com/samthom/activity-logger
RUN dep ensure

RUN go install github.com/samthom/activity-logger
# RUN go build -o main .

ENTRYPOINT /go/bin/activity-logger

EXPOSE 8000

# CMD [ "/go/src/github.com/samthom/activity-logger/main" ]