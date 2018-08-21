FROM golang:alpine

RUN apk update && apk upgrade && \
    apk add --no-cache bash git openssh curl

ENV BASE_PATH=$GOPATH/src/github.com/dmitrymomot/go-test-task-management-tool
ENV PORT=8080

RUN mkdir -p $BASE_PATH
WORKDIR $BASE_PATH

COPY ./ .
COPY ./config /config
COPY ./tpl /tpl

RUN curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
RUN dep ensure -v
RUN go install -v ./cmd/server/...

EXPOSE $PORT

HEALTHCHECK CMD curl --fail http://localhost:$PORT/health || exit 1

CMD ["server"]
