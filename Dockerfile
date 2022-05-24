FROM golang:1.17

WORKDIR /gitlab-bot

COPY . .

#RUN go get -d -t ./... && go install ./... && cd cmd/gitlab-bot &&  go build
RUN go mod download && cd cmd/gitlab-bot && go build

ENV TZ=Europe/Moscow

RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

CMD cmd/gitlab-bot/gitlab-bot

#sample run: docker run --rm -v main:/gitlab-bot --env PROJECT="gitlab/test" --env GITLAB_TOKEN="test" --env TELEGRAM_CHANEL="-1" --env TELEGRAM_BOT_TOKEN="test" --env FATAL_REMINDER="1" colapes/gitlab-bot:latest