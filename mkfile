PROG=mk
GO=go
GOFILES=`ls *.go`
GOOS=
GOARCH=

all:V: $PROG

$PROG:V: $GOFILES
    CGO_ENABLED=0 GOOS=$GOOS GOARCH=$GOARCH $GO build -o $target $prereq

clean:V:
    rm -rf $PROG

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
