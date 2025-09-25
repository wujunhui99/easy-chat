#!/usr/bin/env bash

# Detect the host IP used by local Docker services.
# Prefer IPv4 addresses that belong to active network interfaces.

if [[ -z "${HOST_IP:-}" ]]; then
    detect_host_ip() {
        local candidate=""

        if command -v ip >/dev/null 2>&1; then
            candidate=$(ip route get 1 2>/dev/null | awk '/src/ {for (i=1; i<=NF; i++) if ($i=="src") {print $(i+1); exit}}')
        fi

        if [[ -z "$candidate" ]] && command -v ipconfig >/dev/null 2>&1; then
            for iface in en0 en1 en2; do
                candidate=$(ipconfig getifaddr "$iface" 2>/dev/null || true)
                [[ -n "$candidate" ]] && break
            done
        fi

        if [[ -z "$candidate" ]]; then
            candidate=$(ifconfig | awk '/inet / {print $2}' | grep -Ev '^(127\.|169\.254\.|0\.)' | head -n 1)
        fi

        if [[ -z "$candidate" ]]; then
            candidate="127.0.0.1"
        fi

        printf '%s' "$candidate"
    }

    HOST_IP=$(detect_host_ip)
fi

export HOST_IP
echo "HOST_IP=${HOST_IP}"

# sudo chown -R 1001:1001 ./components/etcd/data ./components/etcd/logs
# sudo chmod -R 700 ./components/etcd/data ./components/etcd/logs
