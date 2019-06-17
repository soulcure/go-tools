PWD=`pwd`
echo $PWD
export GOPATH=$PWD
go build  -o ./bin/go-tools  dev
