FROM golang

EXPOSE 8889
EXPOSE 8888

ENV TZ Asia/Shanghai
ENV GO111MODULE on
ENV GOPROXY https://goproxy.io

WORKDIR /code

RUN git clone https://github.com/cs-shuai/leafGame.git

WORKDIR /code/leafGame

RUN go mod vendor

RUN go build -o leafGame main.go

CMD ./leafGame


