FROM golang:1.17

WORKDIR /myapp

COPY . .

RUN go get

RUN go install

RUN go build

ENV TZ=Europe/Moscow

RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

CMD ["/myapp/gitlab-bot"]