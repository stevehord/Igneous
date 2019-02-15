# Igneous.io
Sample TFTP Code Project with limits. 
Files are "PUT" only the the stdout - they are not written to the fs. Not held in memory after being displayed
Files that are Read "Get" must exist on th  

#Installation:
Download to your go/src folder typically 

/Users/{user}/go/src
  or
c:\Go\src

#Build:

cd Igneous/tftp/cmd/tftpd
go build

Expected output tftpd

#Run
chmod +x tftpd
go tftpd

#Unit Tests:
There are two new unit tests that cover readed and writting to the FS

#Integration Tests
./test_put.sh . - sequentially uploads one small file and one large (multi-packet) file
./multi_file_upload.sh - Concurently runs the test_put.sh file three times

Note: There is no validation. You must view the results in the output window



#ToDo:
Error Handling - none
Memory Leak Testing - none
Concurrent File handling - Likely works but with limited testing
