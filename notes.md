### Requirements

1. Master Sends the available port to the client
2. Unique File names on the same datanode
3. Replicates
4. Client download
5. Create config file


### Enhancements

- For the datakeeper that should upload the file, we can take the one with larger available storage * 0.7 + number of available ports * 0.3. Now we are taking random one