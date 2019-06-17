PWD=`pwd`
echo $PWD
export GOPATH=$GOPATH:$PWD
go install dev
