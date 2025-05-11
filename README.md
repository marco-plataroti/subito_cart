# Purchase Cart Service

This is a simple RESTful microservice that calculates pricing and VAT for a purchase order.

The service exposes a `POST /order` endpoint which accepts a list of items (with `product_id` and `quantity`) and returns a detailed price breakdown for the order.

## 🧱 Tech Stack & info

- Language: Go (Golang)
- Runtime: Docker
- Port: 9090
- Structure: Simple clean modular approach (cmd/, internal/, router/, etc.)

---

## 🚀 Installation and Launch

## 🔨 Step 1: Build the Docker image
docker build -t subito_cart .

## 🧱 Step 2: Build the Go binary inside the container
docker run -v $(pwd):/mnt -w /mnt subito_cart ./scripts/build.sh

## 🚀 Step 3: Run the application (localhost:9090)
docker run -v $(pwd):/mnt -p 9090:9090 -w /mnt subito_cart ./scripts/run.sh

## 🧪 Step 4: Run tests inside the container
docker run -v $(pwd):/mnt -w /mnt subito_cart ./scripts/test.sh

---

## 🧪 Example Request

```json
POST /order
{
  "order": {
    "items": [
      { "product_id": 1, "quantity": 1 },
      { "product_id": 2, "quantity": 5 },
      { "product_id": 3, "quantity": 1 }
    ]
  }
}
