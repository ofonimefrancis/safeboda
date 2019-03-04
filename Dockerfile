FROM golang:1.8.5-jessie

# create a working directory
WORKDIR /go/src/github.com/ofonimefrancis/safeboda/

# install packages

# add source code
ADD ./ /go/src/github.com/ofonimefrancis/safeboda/


RUN echo $GOPATH

RUN go get -d -v

# build main.go
RUN go build main.go
# run the binary
EXPOSE 5000

CMD ["./main"]
