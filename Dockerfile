FROM golang

EXPOSE 8889

ENV TZ Asia/Shanghai

RUN git clone https://github.com/cs-shuai/leafGame.git

WORKDIR $GOPATH/src/leafGame

CMD cd leafGame

CMD go mod vendor

CMD go run main.go


