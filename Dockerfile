FROM golang:1.17

WORKDIR /myapp

COPY . .

RUN go get

RUN go install

RUN go build

CMD ["/myapp/gitlab-bot"]