
# open tftp client and put two files with unique file names
/usr/bin/tftp 127.0.0.1 << !
put test_small.txt small$RANDOM.txt
put test_large.txt large$RANDOM.txt