# Docker Cleanup

<p align="center">
  <img src="https://raw.githubusercontent.com/docker-library/docs/c350af05d3fac7b5c3f6327ac82fe4d990d8729c/docker/logo.png" alt="Docker Logo" width="200">
</p>

A powerful CLI tool designed to simplify Docker resource cleanup and recover disk space by automatically removing unused Docker resources.

## 🚀 Features

- **Comprehensive Cleanup**: Remove multiple types of unused Docker resources:
  - ✅ Stopped containers
  - ✅ Unused images
  - ✅ Dangling images (untagged and unreferenced)
  - ✅ Unused volumes
  - ✅ Unused networks
  - ✅ Build caches

- **Safe Operations**:
  - 🔍 Dry-run mode to preview what would be removed
  - 🕒 Age-based filtering (remove resources older than N days)
  - 📊 Display size information before cleanup

- **Easy to Use**:
  - 🧠 Intuitive command structure
  - 📝 Descriptive outputs
  - 🎨 Color-coded results for better visibility

## 📋 Requirements

- Go 1.23+
- Docker installed and running

## 🔧 Installation

### From Source

```bash
# Clone the repository
git clone https://github.com/your-username/docker-cleanup.git

# Navigate to the directory
cd docker-cleanup

# Build the binary
go build -o docker-cleanup

# Move to a directory in your PATH (optional)
sudo mv docker-cleanup /usr/local/bin/
```

## 📖 Usage

### Global Flags

All commands support these flags:

```
--dry-run       Preview what would be removed without actually deleting anything
--older-than N  Only remove resources older than N days
--show-size     Display size information for resources
```

### Commands

#### Cleanup Everything

```bash
docker-cleanup all
```

#### Cleanup Containers

```bash
docker-cleanup containers
```

#### Cleanup Images

```bash
docker-cleanup images
```

#### Cleanup Dangling Images

```bash
docker-cleanup dangling-images
```

#### Cleanup Volumes

```bash
docker-cleanup volumes
```

#### Cleanup Networks

```bash
docker-cleanup networks
```

#### Cleanup Build Caches

```bash
docker-cleanup builds
```

### Examples

Safely preview what would be cleaned up:
```bash
docker-cleanup all --dry-run
```

Remove images older than 30 days:
```bash
docker-cleanup images --older-than 30
```

Clean everything but show disk usage first:
```bash
docker-cleanup all --show-size
```

## 🏗️ Architecture

The application follows a clean architecture pattern:

```
app/
├── cmd/         # Command handlers for CLI interface
├── controllers/ # Business logic layer
├── models/      # Data access layer (Docker API interactions)
└── views/       # Presentation layer
```

## 🤝 Contributing

Contributions are welcome! Here's how you can contribute:

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📜 License

This project is licensed under the MIT License - see the LICENSE file for details.

## 🙏 Acknowledgements

- Docker for their excellent API
- The Go community for their helpful packages