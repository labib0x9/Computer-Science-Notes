---
title: "Capstone Project: PIKA, a Toy Key-Value Protocol and Server in C"
date: 2026-10-28
tags:
  - c
  - networking
  - sockets
  - pthreads
  - capstone
summary: A small text protocol over TCP in C, built as protocol, then sockets, then three concurrency models, then manual memory vs a garbage collector
---

# Overview

PIKA is a small text protocol for a key-value store, served over TCP. The project has three stages: **protocol, sockets, multithreading**.

- PIKA Protocol
- Build the Protocol Parser
- Build the key-value store using hashmap and linked list
- Client-server communication using socket (over TCP) for single client
- Integrate Thread Pool for handling multiple clients
- Improve multiple clients handling by event loop
- Intigrate LRU cache