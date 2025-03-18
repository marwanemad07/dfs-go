### Requirements

1. Unique File names on the same datanode  (checksum)
2. Client download
3. run on diffrent machines 

### Enhancements

- For the datakeeper that should upload the file, we can take the one with larger available storage _ 0.7 + number of available ports _ 0.3. Now we are taking random one
- download in parall
### Bugs

### Atef tasks:

- images:
  - to build `docker build -t datanode`
  - to run `docker run -d -p <host_port:container_port datanode container_port>`
- checksum

### tests
- large files 
- one datanode down in the middle of upload or download
- multiple user upload on the same time
- same file name upload

### Questions 
- in replication in data keeper node i use the tcp like the client is that correct 