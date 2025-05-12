# Non-FaasFlows-Demo-Functions

## Overview

This repository contains example serverless functions built using the standard OpenFaaS templates (`golang-http`) without the FaasFlows enhancements. It serves as a comparison baseline to evaluate the performance benefits and developer experience improvements offered by FaasFlows.

The functions implement the same ride-sharing application workflow as in the FaasFlows-Demo-Functions repository, but using traditional HTTP requests between functions rather than the optimized workflow approach.

---

## Prerequisites

### Required Infrastructure

1. **Kubernetes cluster** - A running Kubernetes cluster (local or remote)
    - Minikube, kind, or a managed Kubernetes service (EKS, GKE, AKS)
    - Properly configured `kubectl` with access to your cluster

2. **Docker** - For building and pushing function images
    - Docker Desktop or Docker Engine with Docker CLI
    - Access to Docker Hub or a private container registry

3. **Standard OpenFaaS installation** - No modifications required
    - Deployed using Helm or the official installation methods
    - Properly configured OpenFaaS gateway

### Required Tools
1. **faas-cli** - OpenFaaS CLI tool
   ```bash
   curl -sL https://cli.openfaas.com | sudo sh
   ```
2. **Node.js** - For running test scripts
    - Required for executing the migration and redirection tests
3. **Custom golang-flow Template** - For building functions
    ```bash
   faas-cli template pull https://github.com/danenherdi/golang-http-template
   
   # OR

    git clone https://github.com/danenherdi/golang-http-template.git
    
   # Copy the template to your local OpenFaaS templates directory
   mkdir -p ~/.openfaas/templates
   cp -r ~/golang-http-template/template/golang-flow ~/.openfaas/templates/
   
   # Verify the template is available
   faas-cli template list
   faas-cli new --list
   ```
   It should show:
    ```
    Languages available as templates:
   - golang-flow
   - golang-http
   - golang-middleware
    ```
4. **K6** - For performance testing and comparison

   Before running tests, ensure K6 is installed

    ```bash
    # Linux
    sudo gpg -k
    sudo gpg --no-default-keyring --keyring /usr/share/keyrings/k6-archive-keyring.gpg --keyserver hkp://keyserver.ubuntu.com:80 --recv-keys C5AD17C747E3415A3642D57D77C6C491D6AC1D69
    echo "deb [signed-by=/usr/share/keyrings/k6-archive-keyring.gpg] https://dl.k6.io/deb stable main" | sudo tee /etc/apt/sources.list.d/k6.list
    sudo apt-get update
    sudo apt-get install k6
    
    # macOS
    brew install k6
    
    # Windows
    choco install k6
    ```
   
## Setup and Deployment

### 1. Clone the repository

```bash
git clone https://github.com/danenherdi/Non-FaasFlows-Demo-Functions.git
cd Non-FaasFlows-Demo-Functions
```

### 2. Verify and Configure OpenFaaS Gateway
Ensure your OpenFaaS gateway is running and accessible.
```bash
# Set your OpenFaaS gateway URL
export OPENFAAS_URL=http://127.0.0.1:8080

# Login if authentication is enabled
faas-cli login --password <your-password>

# Verify connection
faas-cli list
```

### 3. Build all functions
To build all functions defined in the .yml files:
```bash
faas-cli build -f user-info.yml
faas-cli build -f ride-history.yml
faas-cli build -f ride-recommend.yml
faas-cli build -f last-ride.yml
faas-cli build -f homepage.yml
```

Or build them all at once (if combined into one YAML later):
```bash
faas-cli build -f stack.yml
```

### 4. Push (optional, if you use private registry)
If needed, push to your Docker registry:
```bash
faas-cli push -f user-info.yml
faas-cli push -f ride-history.yml
faas-cli push -f ride-recommend.yml
faas-cli push -f last-ride.yml
faas-cli push -f homepage.yml
```
(If using local Docker and Kubernetes integration, you may skip push.)

### 5. Deploy the functions
Deploy to your OpenFaaS Gateway:
```bash
faas-cli deploy -f user-info.yml
faas-cli deploy -f ride-history.yml
faas-cli deploy -f ride-recommend.yml
faas-cli deploy -f last-ride.yml
faas-cli deploy -f homepage.yml
faas-cli deploy -f hello-world.yml
```
Each function will be available through the OpenFaaS Gateway.

## Function Overview
The demo functions implement a simple ride-sharing application workflow:
1. **homepage** - Entry point function that aggregates data from other functions by making direct HTTP requests 
2. **user-info** - Provides user profile information 
3. **last-ride** - Retrieves information about the user's most recent ride 
4. **ride-history** - Shows a history of past rides 
5. **ride-recommend** - Recommends rides based on user history and location 
6. **friends** - Shows friends who have used the service 
7. **hello-world** - A simple test function

## Testing

The repository includes several K6 test scripts to verify FaasFlows functionality and measure performance improvements. [K6](https://k6.io/) is a modern load testing tool that allows us to simulate various traffic patterns and measure key performance metrics.

### Running Migration Tests
Migration tests verify the effectiveness of the function cluster migration method in reducing vendor lock-in. These tests simulate the migration of functions between environments while measuring performance impact.

```bash
# Run migration tests
k6 run non-faasflows-migration-test.js
k6 run non-faasflows-migration-homepage-test.js
k6 run non-faasflows-migration-friends-test.js
k6 run non-faasflows-migration-ride-history-test.js
```

### Running Request Redirection Tests
Request redirection tests measure the effectiveness of the redirecting requests to the nearest function method in reducing response time.

```bash
# Run redirection test
k6 run non-faasflows-redirection-test.js
```

## Notes on Implementation
1. In this standard implementation, each function needs to explicitly handle request routing and response parsing 
2. Error handling is more complex as failures in dependent functions must be managed manually 
3. Deployments and updates require careful coordination when function interfaces change 
4. The code for each function is typically longer and more complex due to the additional HTTP request management