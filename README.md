# Igneous.io
Sample TFTP Code Project with limits. 
Files are "PUT" only into local memory and written to stdout for dubgging. They are not written to the filesystem. "GETS" must process against a file that has been successfully "PUT" using the same filename.  "GET"ing a non existant file will return zero bytes

# Installation:
Download to your go/src folder typically 

/Users/{user}/go/src
  or
c:\Go\src

# Build:

cd Igneous/tftp/cmd/tftpd

go build


# Run
Depending on your system you may need to make the output executable
chmod +x tftpd

note: server opens port 69 and may require priveldges 

sudo ./tftpd

# Testing with a Mac built-in tftp Client
/usr/bin/tftp 127.0.0.1

put testFiles/test_small.txt

# Unit Tests:
None

# Integration Tests
./test_put.sh . - sequentially uploads one small file and one large (multi-packet) file

./test_get.sh . - sequentially uploads one small file and one large (multi-packet) file and retrieves them two new files are created in the testFiles directory

./multi_file_upload.sh - Concurently runs the test_put.sh file three times 

Note: There is no validation. You must view the results in the output window and check file sizes


# ToDo:
Windows Build - none

Error Handling - limited 

Memory Leak Testing - none

Concurrent File handling - Works but with lots of limits. Needs improved testing

Uploading multiple files with the same file name concurrently is not supported. It will produce incorrect results. There is no attempt to guarentee unique file names


