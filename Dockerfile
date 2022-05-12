FROM golang:1.17

WORKDIR /myapp

COPY . .

RUN go get -d -t ./... && go install ./... && cd cmd/gitlab-bot &&  go build

ENV TZ=Europe/Moscow

RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

CMD cmd/gitlab-bot/gitlab-bot