FROM golang

#ADD . /go/src/bookAPI/apiservice
#WORKDIR /go/src/bookAPI/apiservice
#RUN go get github.com/lib/pq
#RUN go install /go/src/bookAPI/apiservice

#ENTRYPOINT /go/bin/apiservice

ADD ./apiservice /go/bin
WORKDIR /go/bin
EXPOSE 8080
ENTRYPOINT /go/bin/apiservice

CMD [ "/bin/bash" ]