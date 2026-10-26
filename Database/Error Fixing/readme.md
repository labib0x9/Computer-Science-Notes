---
title: Error Fixing
date: 2026-10-26
tags:
  - postgres
  - database
---

# Overview

Every technical error i have encounted

# Issue - 01

```bash
$ brew services start postgres    
Bootstrap failed: 5: Input/output error
Try re-running the command as root for richer errors.
Error: Failure while executing; `/bin/launchctl bootstrap gui/501 /Users/labib0x9/Library/LaunchAgents/sh.brew.postgresql@18.plist` exited with 5.

$ brew services list
Name          Status User     File
cloudflared   none            
postgresql@18 other  labib0x9 ~/Library/LaunchAgents/sh.brew.postgresql@18.plist
rabbitmq      none            
redis         none            

$ launchctl bootout gui/$(id -u)/sh.brew.postgresql@18

$ brew services start postgresql@18
==> Successfully started `postgresql@18` (label: sh.brew.postgresql@18)
```