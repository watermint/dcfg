
FROM golang:1.6.2

RUN apt-get update -y
RUN apt-get upgrade -y
RUN apt-get install -y zip

RUN cd $GOPATH

RUN mkdir /dist
RUN go get -u github.com/cihub/seelog
RUN go get -u github.com/dropbox/dropbox-sdk-go-unofficial
RUN go get -u google.golang.org/api/admin/directory/v1
RUN go get -u google.golang.org/cloud/compute/metadata

ADD . $GOPATH
ENTRYPOINT $GOPATH/build/build_inside_docker.sh
