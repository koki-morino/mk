PROG=mk
GO=go
GOCOMPAT=1.12
GOPATH=`go env GOPATH`
GOFILES=`ls *.go`
GOOS=
GOARCH=

all:V: $PROG

$PROG:V: $GOFILES
    CGO_ENABLED=0 GOOS=$GOOS GOARCH=$GOARCH $GO build -o $target $prereq

clean:V:
    rm -rf $PROG

compat-install:V: $HOME/sdk/go$GOCOMPAT

$HOME/sdk/go$GOCOMPAT: $GOPATH/bin/go$GOCOMPAT
    go$GOCOMPAT download

$GOPATH/bin/go$GOCOMPAT:
    go install -v golang.org/dl/go$GOCOMPAT@latest

test:V: testpid testshell testfaile

testpid:V:
    echo pid=$pid

testshell:V:
    echo shell=$MKSHELL

testfail:V:
    false
    true

testfaild:VD:
    touch testfaild
    false

testfaile:VE:
    false
    true
