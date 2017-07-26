FROM golang:1.8

VOLUME /home/colin/gocode/src/github.com/colindev/ssh-agent:/go/src/app
WORKDIR /go/src/app

RUN ls /go/src/app
RUN go-wrapper download

CMD ["go", "build", "-o", "xx"]

