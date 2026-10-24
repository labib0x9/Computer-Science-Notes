---
title: "Container Root Isolation with Chroot and Namespaces (Go)"
date: 2026-10-25
tags:
  - golang
  - linux
  - containers
  - code-explained
summary: "Implements a minimal container with Linux namespaces, UID/GID mapping, chroot, and private mounts."
---

## Overview

This program demonstrates how container engines isolate processes on Linux using kernel namespaces and filesystem root changes (`chroot`). By re-executing itself with specific clone flags, it creates an unprivileged user namespace mapped to root, isolates process IDs and hostnames, remounts the filesystem privately, changes the root directory to a container rootfs, and launches an interactive shell inside the isolated environment.

## Source

https://github.com/labib0x9/Programming-Notes/blob/main/Programming/Golang%20System/02-Container/chroot.go

## Concepts used

- Linux namespaces (`CLONE_NEWUSER`, `CLONE_NEWPID`, `CLONE_NEWUTS`, `CLONE_NEWNS`)
- Rootless user namespace mapping (`UidMappings`, `GidMappings`)
- Process self-re-execution pattern
- Filesystem isolation with `chroot` and private mounts (`MS_PRIVATE | MS_REC`)
- Pseudo-filesystems (`tmpfs`, `proc`)

## How it works

The code specifies the `//go:build linux` build constraint, restricting compilation to Linux systems where namespace and mount system calls are supported.

The function `logger` provides standardized error logging, panicking if an error is encountered.

The `main` function evaluates command-line arguments passed via `os.Args[1:]`. When invoked with `"run"`, it invokes `parent()`. When invoked with `"child"`, it invokes `child()`. If arguments are missing or invalid, it prints usage instructions.

The `parent` function implements the self-re-execution pattern:
- Retrieves the current binary path via `os.Executable()`.
- Prepares a command `exec.Command(exePath, "child")`, binding standard input, standard output, and standard error.
- Configures `cmd.SysProcAttr` with `Cloneflags` specifying `CLONE_NEWUSER` (enabling rootless containerization without host root), `CLONE_NEWPID` (isolating process tree so child becomes PID 1), `CLONE_NEWUTS` (isolating hostname), and `CLONE_NEWNS` (isolating mount points).
- Configures `UidMappings` and `GidMappings` to map the current host user ID (`os.Getuid()`) and group ID (`os.Getgid()`) to container ID 0 (`root`).
- Sets `Pdeathsig: syscall.SIGKILL` so the child is killed if the parent process dies.
- Calls `cmd.Run()` to start the child and waits for its termination.

The `child` function executes inside the newly isolated namespaces:
- Calls `syscall.Mount("", "/", "", syscall.MS_PRIVATE|syscall.MS_REC, "")` to prevent mount events inside the container from propagating back to the host filesystem.
- Sets the container hostname to `"toor"` using `syscall.Sethostname([]byte("toor"))`.
- Calls `syscall.Chroot(rootfsPath)` pointing to `"/home/labib0x9/rootfs/alpine_tmp"`, locking the process filesystem root to that folder.
- Changes working directory to `"/"` via `syscall.Chdir("/")`.
- Mounts a `tmpfs` instance at `"/tmp"` and mounts the process filesystem at `"/proc"`.
- Sets deferred calls to unmount `/proc` and `/temp`.
- Spawns an interactive `sh` command with standard I/O wired to the terminal.

## Key excerpt

```go
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUSER |
			syscall.CLONE_NEWPID |
			syscall.CLONE_NEWUTS |
			syscall.CLONE_NEWNS,
		Pdeathsig: syscall.SIGKILL,
		UidMappings: []syscall.SysProcIDMap{
			{
				ContainerID: 0,
				HostID:      os.Getuid(),
				Size:        1,
			},
		},
		GidMappings: []syscall.SysProcIDMap{
			{
				ContainerID: 0,
				HostID:      os.Getgid(),
				Size:        1,
			},
		},
		GidMappingsEnableSetgroups: false,
	}
```

## Expected behavior (from reading the code, not run)

To run the program on a compatible Linux environment:
```bash
go run chroot.go run
```

Expected behavior:
The parent spawns a child process within isolated user, PID, UTS, and mount namespaces. Inside the child, the hostname changes to `toor`, the root switches to `/home/labib0x9/rootfs/alpine_tmp`, pseudo-filesystems are mounted, and an interactive `/bin/sh` shell prompt opens as `root` (UID 0 within the container). Typing `exit` or `Ctrl+D` closes the shell, triggers deferred cleanup, and returns to the host shell.

## Things to watch

- **Linux platform dependency:** The file contains `//go:build linux` and relies on Linux-specific syscalls (`CLONE_NEW*`, `syscall.Mount`, `syscall.Sethostname`). It cannot compile or run on macOS or Windows.
- **Hardcoded rootfs path:** The path `/home/labib0x9/rootfs/alpine_tmp` is hardcoded to a specific local user directory. If this directory does not exist on the target machine, `syscall.Chroot` fails with `no such file or directory`.
- **Typo in deferred unmount:** The deferred cleanup on line 77 specifies `syscall.Unmount("/temp", 0)`, but line 73 mounted to `"/tmp"`. Unmounting `"/temp"` will fail with an error because that path was never mounted.
- **chroot isolation limits:** Standard `chroot` does not change the root mount in the same manner as `pivot_root`. If a process inside the container has permission to call `chroot` or retains an open file descriptor outside the jail, escape is possible.

## Terms

- **Linux Namespace:** A kernel mechanism that partitions system resources (such as process trees, user IDs, network devices, and mount points) so processes see isolated virtual instances.
- **Rootless Container:** A container that executes without requiring real root privileges on the host, relying on user namespace mapping.
- **chroot:** A Unix system call that changes the apparent root directory for the current running process and its children.
- **tmpfs:** A temporary filesystem residing in volatile memory rather than on persistent disk storage.

## Related

- [docker api communication](../03-Docker%20Api/docker%20api%20comminucation.md)

*This explanation was generated by an AI from reading the source, not from running it. Verify before relying on it.*
