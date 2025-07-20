# gocanarygo

A lightweight CLI tool for managing and testing canary deployments in Kubernetes clusters.

## 🔧 Features

- ✅ Deploy applications with autoscaling (HPA)
- 🧹 Cleanup deployments and autoscalers
- 📦 View current status of deployments and HPAs
- 📌 Check CLI version

---

## 🚀 Quickstart

### 1. Prerequisites

- Go 1.21+
- Docker & Kind
- kubectl
- metrics-server running in your cluster

### 2. Clone & Build

```bash
git clone https://github.com/jrrbb/gocanarygo.git
cd gocanarygo
go mod tidy
