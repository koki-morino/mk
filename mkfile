PROG=mk
GOFILES=`ls *.go`

all:V: $PROG

$PROG: $GOFILES
    CGO_ENABLED=0 go build -o $target $prereq

clean:V:
    rm -rf $PROG
