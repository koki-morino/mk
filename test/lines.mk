LINES=\
    echo "line 1"\
    # mk allows continuation of lines with backslashes even in comments\
    echo "line 2"\

test/lines.mk:V:
	echo $LINES
