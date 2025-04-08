# 🛒 E-Commerce Microservices Platform

This is a basic e-commerce platform built as part of the **Advanced Programming II – Assignment 1**. The project is built using **Clean Architecture** principles and includes **three microservices** communicating via REST.

> **Student:** Beibars Yergali  
> **University:** Astana IT University  
> **Course:** Advanced Programming II  
> **Assignment:** Implementing a Clean Architecture-Based Microservices

---

## 🧠 Objective

Build a basic e-commerce platform composed of three independent services:

1. **API Gateway (Gin)** – Handles routing, logging, telemetry, and authentication.
2. **Inventory Service (Gin + DB)** – Manages product inventory, categories, and stock levels.
3. **Order Service (Gin + DB)** – Handles order creation, updates, payments, and tracking.

---

## 🏗️ Project Structure



Each service follows Clean Architecture, separating domain logic, interfaces, infrastructure, and delivery.

---

## 📦 Inventory Service

Manages **products**, **categories**, **stock**, and **prices**.

### ✨ Features
- CRUD for products and categories
- Filtering & pagination for product listing

### 🔌 Endpoints

| Method | Endpoint             | Description                |
|--------|----------------------|----------------------------|
| POST   | `/products`          | Create a new product       |
| GET    | `/products/:id`      | Retrieve a product by ID   |
| PATCH  | `/products/:id`      | Update a product           |
| DELETE | `/products/:id`      | Delete a product           |
| GET    | `/products`          | List all products          |

---

## 🧾 Order Service

Handles **order creation**, **status updates**, and **payment tracking**.

### ✨ Features
- Associate orders with products and quantities
- View order status

### 🔌 Endpoints

| Method | Endpoint             | Description                     |
|--------|----------------------|---------------------------------|
| POST   | `/orders`            | Create a new order              |
| GET    | `/orders/:id`        | Retrieve order by ID            |
| PATCH  | `/orders/:id`        | Update order status             |
| GET    | `/orders`            | List all orders for a user      |

---

## 🚪 API Gateway

Acts as the single access point to the services:
- Handles routing to Inventory and Order services
- Logging and telemetry integration
- Placeholder for future auth middleware

---

## 🛠️ Tech Stack

- **Go (Gin Web Framework)** – Backend for all services
- **Any DB (e.g., PostgreSQL, MongoDB)** – Persistent storage per service
- **Docker (optional)** – For containerized microservices
- **REST API** – Communication between services and client

---

## 📂 How to Run

Each service can be run individually:

```bash
# Example: Run Inventory Service
cd inventory-service
go run main.go


