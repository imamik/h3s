[![codecov](https://codecov.io/gh/imamik/h3s/branch/main/graph/badge.svg)](https://codecov.io/gh/imamik/h3s)

# H3S - Hetzner K3s Cluster Manager

A powerful CLI tool for creating, managing, and operating K3s clusters on
Hetzner Cloud. H3S simplifies the process of deploying and managing Kubernetes
clusters while providing advanced features for cluster lifecycle management.

## Features

- 🚀 Quick cluster creation and management
- 🔒 Built-in security features and firewall configuration
- 🔄 Easy cluster scaling and node management
- 🛠 Integrated kubectl command support
- 📦 Automated dependency management
- 🔑 SSH access to cluster nodes
- ⚡ Fast and efficient cluster operations
- 🌐 Automatic DNS management with Hetzner DNS
- 🔐 Built-in cert-manager integration
- 💻 Interactive configuration creation

## Prerequisites

- Hetzner Cloud API token
- Hetzner DNS API token (optional, for DNS management)
- kubectl installed locally
- SSH key pair for node access

## Installation

```bash
# Using go install
go install github.com/imamik/h3s@latest

# Or clone and build manually
git clone https://github.com/imamik/h3s.git
cd h3s
make build
```

## Quick Start

1. Create your configuration interactively:

```bash
# Create cluster configuration
h3s create config

# Create credentials configuration
h3s create credentials
```

Or manually create the `h3s.yaml` and `h3s-secrets.yaml` configuration files as
described in the [Configuration](#configuration) section below.

2. Create a cluster:

```bash
h3s create cluster
```

## Available Commands

### Create Commands

```bash
# Create a new configuration interactively
h3s create config

# Create credentials configuration interactively
h3s create credentials

# Create a new cluster
h3s create cluster
```

### Get Commands

```bash
# Get cluster kubeconfig
h3s get kubeconfig

# Get cluster access token
h3s get token
```

### Cluster Management

```bash
# Create a cluster
h3s create cluster

# SSH into a node
h3s ssh control-plane-1

# Use kubectl commands
h3s kubectl get nodes
h3s kubectl apply -f your-app.yaml

# Destroy cluster
h3s destroy
```

### Installation Commands

```bash
# Install dependencies
h3s install dependencies
```

## Configuration

### Basic Configuration (h3s.yaml)

The configuration is split into two files for better security practices:

- `h3s.yaml`: Contains non-sensitive cluster configuration
- `h3s-secrets.yaml`: Contains sensitive information like API tokens

Both files can be created interactively using the CLI or manually.

#### Main Configuration (h3s.yaml)

```yaml
name: my-cluster
k3s_version: v1.31.1+k3s1
network_zone: eu-central
domain: my-domain.com

# SSH Configuration
ssh_key_paths:
    private_key_path: $HOME/.ssh/id_ed25519
    public_key_path: $HOME/.ssh/id_ed25519.pub

# Certificate Management
cert_manager:
    email: your-email@domain.com
    production: false

# Node Configuration
control_plane:
    as_worker_pool: true
    pool:
        name: control-plane
        nodes: 3
        location: nbg1
        instance: cax11

worker_pools:
    - name: worker
      nodes: 2
      location: nbg1
      instance: cax11
```

#### Secrets Configuration (h3s-secrets.yaml)

```yaml
hcloud_token: your-hetzner-cloud-token
hetzner_dns_token: your-hetzner-dns-token
k3s_token: your-k3s-token
```

## Architecture

H3S follows a modular architecture with the following components:

- **cmd**: Command-line interface and command handlers
  - create: Cluster and configuration creation
  - get: Retrieve cluster information
  - destroy: Cluster cleanup
  - ssh: Node access
  - kubectl: Kubernetes operations
  - install: Dependency management
- **internal**: Core business logic and implementations
  - cluster: Cluster management logic
  - hetzner: Hetzner Cloud API integration
  - k3s: K3s specific operations
  - config: Configuration management

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the
[LICENSE.txt](LICENSE.txt) file for details.
