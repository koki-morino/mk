PROG=mk

# Go 1.12 introduced *os.ProcessState.ExitCode()
GC=go1.12.17
GCFLAGS=-s -w
GOOS=
GOARCH=

SRC=\
    expand.go\
    graph.go\
    lex.go\
    mk.go\
    parse.go\
    recipe.go\
    rules.go\

TESTS=\
    mk_test.go

all:V: $PROG

$PROG: $SRC
    CGO_ENABLED=0 GOOS=$GOOS GOARCH=$GOARCH $GC build -o $target $prereq

clean:V:
    # Don't use the Go module system for backwards compatibility
    rm -f $PROG go.mod go.sum

fmt:V:
    $GC fmt $SRC $TESTS

install:V:
    $GC install -ldflags="$GCFLAGS" .

test:V: $PROG
    $GC test -v $TESTS
