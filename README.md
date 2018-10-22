# Orbit Drive

Proof of concept storage file system built on top of [IPFS](https://ipfs.io/) + [Infura](https://infura.io)

## Requirements

- golang 1.9
- [ipfs](https://docs.ipfs.io/introduction/install/) (if running a gateway locally)

## Ubuntu 16.04

- CLI

Initialize user settings
```bash
go run orbit-drive.go init -r [Path of folder to sync] -p [Password] -n [Ipfs gateway]
```

Start synchronizing folder
```bash
go run orbit-drive sync
```

- Register Service

```bash
# Move orbit-drive service file to system
sudo mv ./orbit-drive.service /lib/systemd/system/.

# Set permission
sudo chmod 755 /lib/systemd/system/orbit-drive.service

# Load and start service
sudo systemctl daemon-reload
sudo systemctl enable orbit-drive.service
sudo systemctl start orbit-drive
```