# Igneous.io
Sample TFTP Code Project with limits. 
Files are "PUT" only the the stdout - they are not written to the fs. Not held in memory after being displayed
Files that are Read "Get" must exist on th  

# Installation:
Download to your go/src folder typically 

/Users/{user}/go/src
  or
c:\Go\src

# Build:

cd Igneous/tftp/cmd/tftpd

go build

Expected output tftpd

# Run
chmod +x tftpd
go tftpd

# Testing with a Mac built in tftp Client
/usr/bin/tftp 127.0.0.1
put testFiles/test_small.txt

# Unit Tests:
None

# Integration Tests
./test_put.sh . - sequentially uploads one small file and one large (multi-packet) file
./test_get.sh . - sequentially uploads one small file and one large (multi-packet) file and retrieves them 
./multi_file_upload.sh - Concurently runs the test_put.sh file three times 

Note: There is no validation. You must view the results in the output window and check file sizes


# ToDo:
Error Handling - none

Memory Leak Testing - none

Concurrent File handling - Likely works but with limited testing

Uploading multiple files with the same file name concurrently is not supported. It will produce incorrect results. There is no attempt to increment file names


