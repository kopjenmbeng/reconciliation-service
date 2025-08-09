# reconciliation-service

## Overview

The `reconciliation-service` is a microservice designed to handle data reconciliation tasks efficiently. It compares datasets, identifies discrepancies, and provides actionable insights to ensure data consistency across systems. This project is for answering interview chalanges <b>Example 2: reconciliation service (Algorithmic and scaling)</b>

## Features

- **Data Comparison**: Supports reconciliation of structured data from multiple sources.
- **Scalability**: Optimized for handling large datasets. (using go routine when get the banks statement & system-transaction data)
- **Reporting**: Generates detailed reports of discrepancies.

## Installation

1. Clone the repository:
    ```bash
    git clone https://github.com/kopjenmbeng/reconciliation-service.git
    ```
2. Navigate to the project directory:
    ```bash
    cd reconciliation-service
    ```
3. Install dependencies:
    ```bash
    go mod tidy
    ```

## Usage

1. Start the service:
    ```bash
    go run main.go
    ``
## System Design

## System Design
![Optional Text](../master/files/system-design/reconcile-process.jpg)

# Technology Stack
-   GO