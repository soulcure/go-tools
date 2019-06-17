export GOPATH=$GOPATH:${PWD}
echo $GOPATH
go build  -o ./go-tools  dev
