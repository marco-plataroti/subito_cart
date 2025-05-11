# 🛒 Purchase Cart Service

This microservice provides price calculation and VAT breakdown for a purchase order. It is intended to be part of a larger e-commerce platform and demonstrates clean architectural separation, containerized deployment, and simple API design.

---

## ⚙️ Tech Stack

- **Language**: Go (Golang)
- **Runtime**: Docker
- **Port**: 9090 (default)
- **Architecture**: Clean modular structure (`cmd/`, `internal/`, `router/`, etc.)

---

## 🚀 Getting Started

### 🧰 Prerequisites

Ensure you have the following tools installed:

- [Docker](https://docs.docker.com/get-docker/)
- Bash-compatible shell (Linux/macOS or Git Bash for Windows)
- (Optional) [Go](https://go.dev/) — only required for local development outside Docker

### 📦 Clone the Repository

```bash
git clone https://github.com/marco-plataroti/subito_cart.git
# OR SSH (if you have access configured):
# git clone git@github.com:marco-plataroti/subito_cart.git
cd subito_cart
```

## 🏗️ Build & Run

### 🔨 Step 1: Build the Docker image
```bash
docker build -t subito_cart .
```

### 🧱 Step 2: Build the Go binary inside the container

```bash
docker run -v $(pwd):/mnt -w /mnt subito_cart ./scripts/build.sh
```

### ▶️ Step 3: Run the application (localhost:9090)
```bash
docker run -v $(pwd):/mnt -p 9090:9090 -w /mnt subito_cart ./scripts/run.sh
```

### ✅ Step 4: Run tests inside the container
```bash
docker run -v $(pwd):/mnt -w /mnt subito_cart ./scripts/test.sh
```


## 📫 API Reference
POST /order
Calculates the total price and VAT breakdown for an order.

🔸 Request Body
```json
{
  "order": {
    "items": [
      { "product_id": 1, "quantity": 1 },
      { "product_id": 2, "quantity": 5 },
      { "product_id": 3, "quantity": 1 }
    ]
  }
}
```

🔹 Response

```json
{
  "subtotal": 115.0,
  "vat": 23.0,
  "total": 138.0,
  "items": [
    { "product_id": 1, "quantity": 1, "price": 20.0, "vat": 4.0 },
    { "product_id": 2, "quantity": 5, "price": 75.0, "vat": 15.0 },
    { "product_id": 3, "quantity": 1, "price": 20.0, "vat": 4.0 }
  ]
}
```