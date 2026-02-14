PROG=mk
GOFILES=`ls *.go`
GOOS=
GOARCH=

all:V: $PROG

$PROG: $GOFILES
    CGO_ENABLED=0 GOOS=$GOOS GOARCH=$GOARCH go build -o $target $prereq

clean:V:
    rm -rf $PROG
