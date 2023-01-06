Execute given below comamnds:
===============================

pwd=$PWD
export GOPATH=$pwd
export PATH=$PATH:$GOPATH/bin
export GO111MODULE=off
cd src
go get github.com/gorilla/mux
cd main
go build main.go
./main

