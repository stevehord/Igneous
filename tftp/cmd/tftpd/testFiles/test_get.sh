
# remove any existing test results
rm small_out.txt
rm large_out.txt

# open local tftp clinet - put 2 files and get said files
/usr/bin/tftp 127.0.0.1 << !
put test_small.txt
put test_large.txt
get test_small.txt small_out.txt
get test_large.txt large_out.txt