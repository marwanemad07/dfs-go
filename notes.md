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

## Docker
- first you need to build the 2 images (master and data)
  - do that by running the command `docker build -t datanode -f DatanodeDockerfile .` to build the data node image
  - to build the master node image run `docker build -t masternode -f MasternodeDockerfile .`
- now you can start the containers
  - first starting the master container using ` docker run --name master1 -p 9090:9090 masternode 9090` 9090 is the port the master will be using. You can add `-d` flag to detach the container terminal from host terminal
  - now starting the data containers so you can use `docker run --name data1 -p 2025:2025 datanode 2025 9090` 2025 is the port for the data node and 9090 is the port for the master node. Same as before you can use the `-d` flag
  - Last you can now use `go run .\client\client.go upload example.mp4 9090` to upload a file. Note that 9090 is the master port
- You can check the files uploaded in the container using `docker exec -it <container_name_or_id> /bin/sh` to log into the container with the terminal then run linux commands like `ls storage` to check the uploaded files. Use `ctrl + d` to close the container.
- You can use `docker container stop $(docker ps -a -q)` to stop all running containers
- You can use `docker container rm $(docker ps -a -q)` to remove all stopped containers