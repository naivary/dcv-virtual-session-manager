# DCV Virtual Session Manager

This is a simple virtual session manager implemented for Amazon DCV. It allows
you to periodically create virtual sessions for newly added users based on GECOS
information.

# Usage

```bash
git clone
```

```bash
go build -o dcv-virtual-session-manager
```

```ini
# /etc/systemd/system/dcv-virtual-session-manager.service
[Unit]
Description=DCV Virtual Session Manager
Requires=dcvserver.service
After=dcvserver.service

[Service]
Type=simple
Restart=always
RestartSec=1
User=root
ExecStart=/bin/dcv-virtual-session-manager

[Install]
WantedBy=multi-user.target
```

```bash
sudo systemctl start dcv-virtual-session-manager
sudo systemctl enable dcv-virtual-session-manager
```

# Creating a Go DCV managed linux user

```bash
# Create a linux user with the GECOS information `go_dcv_managed` to inform the dcv-virtual-session-manager 
# process that it should manage that user. The first GECOS information should always be your display name and then
# you can add the go_dcv_managed annotation.
sudo useradd -d /home/u004 -s /bin/bash -c u004,go_dcv_managed u004
```
