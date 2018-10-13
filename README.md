# IP Drive

Proof of concept storage file system built on top of [IPFS](https://ipfs.io/) + [Infura](https://infura.io)

## Ubuntu 16.04

- CLI

```bash
go run main.go -root=/home/{user}/ipfs-folder
```

- Register Service

```bash
# Move ip-drive service file to system
sudo mv ./ip-drive.service /lib/systemd/system/.

# Set permission
sudo chmod 755 /lib/systemd/system/ip-drive.service

# Load and start service
sudo systemctl daemon-reload
sudo systemctl enable ip-drive.service
sudo systemctl start ip-drive
```