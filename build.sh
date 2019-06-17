export GOPATH=$GOPATH:${PWD}
echo $GOPATH
go build  -o ./go-server  dev
