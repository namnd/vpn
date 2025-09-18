#!/bin/bash

# Install tailscale
curl -fsSL https://tailscale.com/install.sh | sh

sudo tailscale up \
  --advertise-exit-node \
  --hostname={{.HostName}} \
  --ssh \
  --advertise-tags=tag:vpn \
  --authkey={{.TailscaleAuthKey}}
