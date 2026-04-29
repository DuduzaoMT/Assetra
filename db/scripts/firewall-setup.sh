#!/bin/bash
set -euo pipefail
# ==============================================================================
# Database Firewall Configuration (PostgreSQL + K8s Cluster)
# Interfaces: enp0s3 (Internal) | enp0s8 (Internet NAT)
# Usage: sudo ./firewall-setup.sh [INTERNAL_IF] [EXTERNAL_IF]
# ==============================================================================

# --- CONFIGURATION ---
# Internal interface (DB network - 192.168.0.x)
IF_INT="${1:-enp0s3}"

# Internet interface (Only for updates/packages)
IF_NET="${2:-enp0s8}"

# Bastion Host / Tailscale Gateway (the only one allowed to SSH)
# ADMIN_IP="192.168.0.X"

# IPs of Kubernetes nodes (Control Plane and Worker Nodes)
TRUSTED_IPS=(
    "192.168.0.10" # CP
    "192.168.0.9"  # WN-1
    "192.168.0.12" # WN-2
)

# Allowed ports
DB_PORT="5432" # PostgreSQL
SSH_PORT="22"

echo "[*] Starting firewall configuration..."

# ==============================================================================
# 1. Kernel hardening
# ==============================================================================
# Disable IPv6
sysctl -w net.ipv6.conf.all.disable_ipv6=1 >/dev/null
sysctl -w net.ipv6.conf.default.disable_ipv6=1 >/dev/null
sysctl -w net.ipv6.conf.lo.disable_ipv6=1 >/dev/null

# Anti-spoofing (Prevents forged packets)
sysctl -w net.ipv4.conf.all.rp_filter=1 >/dev/null
sysctl -w net.ipv4.conf.default.rp_filter=1 >/dev/null

# Allow pings (ICMP echo) - changed for connectivity testing
sysctl -w net.ipv4.icmp_echo_ignore_all=0 >/dev/null

# ==============================================================================
# 2. Clear existing rules
# ==============================================================================
iptables -F
iptables -X
iptables -t nat -F
iptables -t nat -X

# ==============================================================================
# 3. Default policies (Drop everything)
# ==============================================================================
iptables -P INPUT DROP
iptables -P FORWARD DROP
iptables -P OUTPUT DROP

# ==============================================================================
# 4. Base system rules
# ==============================================================================
# Loopback
iptables -A INPUT -i lo -j ACCEPT
iptables -A OUTPUT -o lo -j ACCEPT

# Allow established traffic (responses to permitted requests)
iptables -A INPUT -m conntrack --ctstate ESTABLISHED,RELATED -j ACCEPT
iptables -A OUTPUT -m conntrack --ctstate ESTABLISHED,RELATED -j ACCEPT

# Drop invalid packets
iptables -A INPUT -m conntrack --ctstate INVALID -j DROP

# ==============================================================================
# 5. Internet (External Interface) - Outbound only 
# TODO: All traffic should be routed through the gateway instead of allowing direct outbound here
# ==============================================================================
echo "[*] Allowing outbound traffic on $IF_NET (for apt/docker updates)..."
iptables -A OUTPUT -o "$IF_NET" -j ACCEPT

# ==============================================================================
# 6. Logging (For debug/audit)
# ==============================================================================
iptables -N LOGDROP || true
iptables -A LOGDROP -m limit --limit 5/min -j LOG --log-prefix "FW-BLOCKED: " --log-level 4
iptables -A LOGDROP -j DROP

# ==============================================================================
# 7. Internal Network (Internal Interface) - Strictly restricted
# ==============================================================================

# 7.1. Admin access (SSH through Bastion/Gateway)
#echo "[*] Allowing exclusive SSH for Admin via $ADMIN_IP..."
#iptables -A INPUT -i "$IF_INT" -s "$ADMIN_IP" -p tcp --dport "$SSH_PORT" -j ACCEPT
#iptables -A OUTPUT -o "$IF_INT" -d "$ADMIN_IP" -m conntrack --ctstate ESTABLISHED,RELATED -j ACCEPT

# Iterate over all cluster nodes to grant them access
for IP in "${TRUSTED_IPS[@]}"; do
    echo "[*] Allowing access from K8s node: $IP via $IF_INT..."
    
    # PostgreSQL
    iptables -A INPUT -i "$IF_INT" -s "$IP" -p tcp --dport "$DB_PORT" -j ACCEPT
    
    # Ping (ICMP Echo Request)
    iptables -A INPUT -i "$IF_INT" -s "$IP" -p icmp --icmp-type echo-request -j ACCEPT
    
done

# Block any attempt by the DB to initiate a new connection (NEW) to the internal network
# The DB only responds to requests, it never initiates them.
iptables -A OUTPUT -o "$IF_INT" -m conntrack --ctstate NEW -j LOGDROP

# Drop and log everything else
iptables -A INPUT -j LOGDROP
iptables -A OUTPUT -j LOGDROP

# ==============================================================================
# 8. Persisting Rules (Ubuntu/Debian)
# ==============================================================================
if ! command -v netfilter-persistent >/dev/null 2>&1; then
    echo "[*] Installing netfilter-persistent..."
    apt-get update -qq
    DEBIAN_FRONTEND=noninteractive apt-get install -y -qq netfilter-persistent iptables-persistent
fi

echo "[*] Saving iptables rules to persist after reboot..."
netfilter-persistent save

# Persist Kernel settings
SYSCTL_CONF="/etc/sysctl.d/99-db-firewall.conf"
echo "[*] Saving sysctl settings to ${SYSCTL_CONF}..."
cat > "${SYSCTL_CONF}" <<EOF
# DB firewall hardening settings
net.ipv6.conf.all.disable_ipv6=1
net.ipv6.conf.default.disable_ipv6=1
net.ipv6.conf.lo.disable_ipv6=1
net.ipv4.conf.all.rp_filter=1
net.ipv4.conf.default.rp_filter=1
net.ipv4.icmp_echo_ignore_all=0
EOF
sysctl --system >/dev/null || true

echo "======================================================="
echo "   PRODUCTION FIREWALL ACTIVE AND SAVED"
echo "======================================================="
echo " - Internet ($IF_NET): Outbound only."
echo " - Internal ($IF_INT): Access limited to K8s node IPs."
echo " - Pings: Allowed only from the cluster."
