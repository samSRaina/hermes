# API Gateway & Trip Service Interaction Report

## Executive Summary
This report details the inner workings and inter-service communication between the `api-gateway` and `trip-service` within Hermes. It breaks down the current data flow for the `/trip/preview` endpoint and highlights critical configuration and architectural issues that will prevent these services from functioning gracefully in production.

---

## 1. Current Architecture & Inner Workings

### The Communication Flow
When a user requests a trip preview, the following sequence of events takes place between the gateway and the trip service:

1. **Client Request**: The client sends a `POST /trip/preview` request to the [api-gateway](file:///home/sam/git/hermes/services/api-gateway/http.go). It includes a JSON body with a `userID`, `pickup` coordinates, and `destination` coordinates.
2. **Artificial Delay**: Once the request hits the gateway, the handler execution is immediately halted for 9 seconds due to a diagnostic sleep `time.Sleep(9 * time.Second)`. 
3. **Synchronous HTTP Call**: The [api-gateway](file:///home/sam/git/hermes/services/api-gateway/http.go) constructs a new JSON request natively and issues a direct, synchronous POST request to `http://localhost:8083/preview`. (The code indicates a future migration to gRPC is planned, via a `TODO` comment). 
4. **Trip Service Processing**: The request is routed to the `trip-service`'s `HandleTripPreview` method in its `http` handler.
5. **External Routing Engine**: The [trip-service](file:///home/sam/git/hermes/services/trip-service/internal/service/service.go) `TripService.GetRoute` contacts an external OSRM API engine (`router.project-osrm.org`) to calculate the actual driving route geometries bounding the coordinates. 
6. **Data Discarding**: The `trip-service` responds cleanly with the fetched route calculations. In response, however, the `api-gateway` reads the JSON payload, unmarshals it into memory, completely discards the contents, and statically returns `{"data": "ok"}` to the client.

### WebSocket Connections (Gateway)
The gateway also manages local Websocket connections over `/ws/riders` and `/ws/drivers` in [ws.go](file:///home/sam/git/hermes/services/api-gateway/ws.go). These functions currently generate stub/dummy profiles for drivers interacting via websockets rather than pushing data to an external service or broker. 

---

## 2. Production Rollout Issues 
When translated strictly into standard Kubernetes production environments via their existing definition files, these services will fail catastrophically down the line. Here is how they pan out during real deployments:

> [!CAUTION]
> **Container Localhost Isolation**
> In `api-gateway`, the HTTP POST for the trip preview is hardcoded to target `http://localhost:8083`. Because each Kubernetes Pod has its own loopback network bounds, the API gateway will attempt to resolve port 8083 internally inside its own container. It will never reach the `trip-service` pod running in Kubernetes cluster IP space, guaranteeing connection timeouts/refusals. 

> [!WARNING]
> **Silenced Application Ports**
> Inside `trip-service/cmd/main.go`, the application is configured to listen statically on port `:8181` by default. 
> However, the `infra/production/k8s/trip-service-deployment.yaml` exposes `containerPort: 8083` and the service binds `targetPort: 8083`. Kubernetes LoadBalancers will forward traffic towards `8083` on the pod, but the actual daemon is listening blindly on `8181` because the K8s manifest never declares the mapped `HTTP_ADDR` environment variable. The service is permanently inaccessible. 

> [!NOTE]
> **Unresolved ConfigMap Variables**
> In the gateway's deployment `infra/production/k8s/api-gateway-deployment.yaml`, an environment variable `GATEWAY_HTTP_ADDR` is injected. However, the application logic directly expects the `HTTP_ADDR` flag, missing the mapped variable name. Luckily, this defaults to `:8081` which aligns seamlessly with the current load balancer configurations, avoiding a complete initialize sequence failure.

---

## 3. Recommended Breaking Changes
Based on the systemic review, the following code-breaking transitions must be planned and applied. Do not implement these directly yet, but add them to your sprint schedule:

1. **Replace hardcoded Localhosts with Service Discovery**
   * Change `http.Post("http://localhost:8083/preview")` to point to a configurable environment variable (e.g. `TRIP_SERVICE_URL`).
   * For production, set this environment variable to the internal Kubernetes DNS name (`http://trip-service:8083`).
2. **Correct Port Misalignments in Kubernetes Defaults**
   * Pre-define the `HTTP_ADDR` environment variable dynamically in `trip-service-deployment.yaml` mapping it clearly to `:8083` -- or alternatively default the `trip-service` HTTP address consistently to `:8083` internally inside `main.go`.
3. **Propagate Functional Payloads** 
   * The `api-gateway` must stop spoofing HTTP strings like `{"data": "ok"}`. Instead, it must proxy the unmarshaled raw geo-data payload received from `trip-service` and proxy it properly to the client response.
4. **Remove Hardcoded Yields**
   * Remove `time.Sleep(9 * time.Second)` in the gateway’s [http handler](file:///home/sam/git/hermes/services/api-gateway/http.go).
5. **Commit to Transport Agnostic Comm (gRPC)**
   * Address the `//TODO: call the trip service via grpc later` mapped comment by transitioning `trip-service` to host a standard listener exposing a gRPC prototype.
