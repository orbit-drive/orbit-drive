#!/bin/bash

# Move ipfs-drive service file to system
sudo mv ./ipfs-drive.service /lib/systemd/system/.

# Permission
sudo chmod 755 /lib/systemd/system/ipfs-drive.service

# Enable service
sudo systemctl daemon-reload
sudo systemctl enable ipfs-drive.service