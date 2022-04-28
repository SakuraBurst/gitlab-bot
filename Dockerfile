FROM golang:1.17

WORKDIR /myapp

COPY . .

RUN cd cmd/gitlab-bot && go get && go install && go build

ENV TZ=Europe/Moscow

RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

CMD cmd/gitlab-bot/gitlab-bot