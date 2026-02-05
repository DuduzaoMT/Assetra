# Assetra

Assetra is a testbed application created to experiment with security protocols and actively search for vulnerabilities. It is not a production asset management platform, but rather a playground for testing, learning, and evaluating security mechanisms in a controlled environment. The ultimate goal is to run Assetra as an application inside a Kubernetes (k8s) cluster, serving as a secure gateway to access self-hosted resources, with VPN integration for protected connectivity.

## Dependencies

- [Docker](https://docs.docker.com/get-docker/) (with `docker compose` support)
- [Go](https://go.dev/doc/install) (version 1.20+ recommended)
- [Node.js & npm](https://nodejs.org/) (for frontend, optional)

## Getting Started

1. **Clone the repository:**

2. **Run the setup script:**

   ```sh
   cd assetra
   ./setup.sh
   ```

3. **Start the application:**

   ```sh
   make build-start
   ```

4. **Access the app:**
   - Backend/API: [http://localhost:9000](http://localhost:9000)
   - Frontend: [http://localhost:3000](http://localhost:3000)

## Useful Commands

- Install Go dependencies:  
  `go mod tidy`
- Install frontend dependencies:  
  `cd frontend && npm install`
- Start services:  
  `make build-start`
- Stop services:  
  `make stop`
