# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Phoenix is a file recovery tool written in Go that recovers files deleted by `rm` command. It works by analyzing ext4 filesystem structures at a low level to find deleted inodes and recover their associated data blocks.

## Architecture

The project is currently a single-file Go application (`main.go`) that:

1. **Requires root privileges** - Must run with sudo to access raw block devices
2. **Takes device path as argument** - e.g., `/dev/sda1`
3. **Uses low-level filesystem analysis** - References `FilesystemAnalyzer` type (not yet implemented)
4. **Reads ext4 superblocks** - To understand filesystem layout and parameters

Key missing components that need implementation:
- `FilesystemAnalyzer` struct and `NewFilesystemAnalyzer()` function
- `ReadSuperblock()` and `PrintFilesystemInfo()` methods
- ext4-specific data structures for superblocks, inodes, and block groups

## Technical Foundation

The project relies on detailed ext4 filesystem knowledge documented in:
- `ext4-official-docs.md` - Official Linux kernel ext4 documentation summary
- `inode-knowledge.md` - inode structure and deletion mechanics
- `file-recovery-analysis.md` - Go implementation strategies for filesystem recovery

Critical ext4 concepts:
- Superblock at offset 1024 bytes with magic number 0xEF53
- Block groups containing inode tables and data blocks
- Deleted files have `dtime` (deletion time) set but inode/data may still exist
- Little-endian byte order (except journal)

## Development Commands

```bash
# Build the application
go build -o phoenix-recovery main.go

# Run with device (requires root)
sudo ./phoenix-recovery /dev/sda1

# Check dependencies
go mod tidy

# Format code
go fmt ./...
```

## Dependencies

- `golang.org/x/sys v0.15.0` - Required for low-level system calls and Unix interfaces
- Go 1.21 minimum version

## Current Status

The project is in early development stage with basic CLI structure in place but core filesystem analysis functionality not yet implemented.