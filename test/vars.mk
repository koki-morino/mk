VAR=val1
VAR=$VAR val2
UNDEFINED=$UNDEFINED

test/vars.mk:V:
    var=val1
    echo $var
    echo $VAR
    echo $UNDEFINED
