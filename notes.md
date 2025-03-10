### Requirements

1. Master Sends the available port to the client
2. Unique File names on the same datanode
3. Replicates
4. Client download
5. Use addresses of devices instead of using localhost:port, See `getAddress()` for more info

### Enhancements

- For the datakeeper that should upload the file, we can take the one with larger available storage _ 0.7 + number of available ports _ 0.3. Now we are taking random one

### Bugs

1. `if time.Duration(now.Sub(lastHeartbeat).Seconds()) >= heartbeatTimeout` this line in master_tracker make the master know that the datanode is dawn after two seconds not one

### Atef tasks:

- filename & filepath
- images:
  - to build `docker build -t datanode`
  - to run `docker run -d -p <host_port:container_port datanode container_port>`
- checksum
- address & port
