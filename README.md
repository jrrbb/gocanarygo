# GoCanaryGo

A lightweight Go-based CLI tool for managing Kubernetes deployments with canary-style rollout and autoscaling support. Ideal for developers and SREs who want fast, scriptable deployment operations.

---

## 🚀 Features

- Deploy Kubernetes apps with HPA (Horizontal Pod Autoscaler)
- Canary-style rollouts and rollback support
- View deployment status, logs, and rollout history
- Simple CLI interface using `cobra`
- Built with client-go and autoscaling v2 API

---

## 📦 Installation

Clone the repo and build it locally:

```bash
git clone https://github.com/your-username/gocanarygo.git
cd gocanarygo
go build -o gocanarygo
```

Or run it without installing:

```bash
go run main.go [command]
```

---

## 🛠️ Usage

```bash
gocanarygo [command] [flags]
```

### Example

```bash
gocanarygo deploy --name nginx-autoscale --image nginx --replicas 2
```

---

## 🧰 Common Commands

### Deploy a new app with autoscaling

```bash
gocanarygo deploy --name nginx-autoscale --image nginx --replicas 2
```

### Check deployment status

```bash
gocanarygo status nginx-autoscale
```

### Stream logs from the first pod

```bash
gocanarygo logs nginx-autoscale
```

### View deployment details

```bash
gocanarygo describe nginx-autoscale
```

### Scale a deployment manually

```bash
gocanarygo scale --name nginx-autoscale --replicas 3
```

### Set up or update autoscaler

```bash
gocanarygo autoscale --name nginx-autoscale --min 2 --max 5 --cpu 60
```

### Roll back to the previous revision

```bash
gocanarygo rollback --name nginx-autoscale
```

### Show rollout history

```bash
gocanarygo history nginx-autoscale
```

### List all deployments

```bash
gocanarygo list
```

### Cleanup a deployment and its autoscaler

```bash
gocanarygo cleanup --name nginx-autoscale
```

---

## 🔧 Development

This tool uses Go modules. If you're hacking on it:

```bash
go mod tidy
go build ./...
```

---

## 📝 License

MIT License. Contributions welcome.

