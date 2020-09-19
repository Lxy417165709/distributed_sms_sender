FROM golang

WORKDIR /root
ADD . /root
RUN /bin/cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && echo 'Asia/Shanghai' >/etc/timezone
RUN go env -w GOPROXY=https://goproxy.io,direct
RUN go build -o main main.go

EXPOSE 8080

ENTRYPOINT ["./main"]