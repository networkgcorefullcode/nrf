<!--
SPDX-FileCopyrightText: 2025 Canonical Ltd
SPDX-FileCopyrightText: 2021 Open Networking Foundation <info@opennetworking.org>
Copyright 2019 free5GC.org

SPDX-License-Identifier: Apache-2.0
-->
[![Go Report Card](https://goreportcard.com/badge/github.com/omec-project/nrf)](https://goreportcard.com/report/github.com/omec-project/nrf)

# NRF

Network Repository Function provides service discovery functionality in the 5G
core network. Each network function registers with NRF with a set of properties
called as NF Profile. Implements 3gpp specification 29.510. NRF Keeps the
profile data for each network function in the MongoDB. Supports Discovery and
registration procedure. When the discovery API is called, NRF fetches a matching
profile from the database and returns it to the caller.


## Repository Structure

Below is a high-level view of the repository and its main components:

```
.
├── accesstoken                           # Access token request handling and routing
│   ├── api_access_token_request.go
│   ├── api_access_token_request_test.go
│   └── routers.go
├── context                               # NRF context management and NF data structures
│   ├── context.go
│   ├── dataconv.go
│   ├── links.go
│   ├── management_data.go
│   ├── search_nf_instance.go
│   └── uri_list.go
├── dbadapter                            # Database abstraction for NF instance storage
│   └── dbadapter.go
├── dev-container.ps1
├── dev-container.sh
├── DEV_README.md
├── discovery                           # NF discovery API logic and routers
│   ├── api_nf_instances_store.go
│   └── routers.go
├── Dockerfile                          # Container build definitions
├── Dockerfile_dev
├── Dockerfile.fast
├── docs                                # Technical documentation and diagrams
│   └── images
│       ├── README-NRF.png
│       └── README-NRF.png.license
├── factory                            # Configuration management and initialization
│   ├── config.go
│   ├── factory.go
│   └── nrf_config_test.go
├── go.mod
├── go.mod.license
├── go.sum
├── go.sum.license
├── LICENSES
│   └── Apache-2.0.txt
├── logger                             # Centralized logging utilities
│   └── logger.go
├── Makefile                           # Build and test automation
├── management                         # NF registration, update, and subscription APIs
│   ├── api_management.go
│   ├── api_nf_instance_id_document.go
│   ├── api_nf_instances_store.go
│   ├── api_subscription_id_document.go
│   ├── api_subscriptions_collection.go
│   └── routers.go
├── metrics                           # Metrics and telemetry collection
│   └── telemetry.go
├── NOTICE.txt
├── nrf.go
├── nrfTest
│   ├── nrfcfg_with_custom_webui_url.yaml
│   └── nrfcfg.yaml
├── polling                           # Configuration polling and synchronization
│   ├── nf_configuration.go
│   └── nf_configuration_test.go
├── producer                          # Core NRF logic: discovery, management, and token generation
│   ├── access_token.go
│   ├── nf_discovery.go
│   ├── nf_management.go
│   └── nf_producer_test.go
├── README.md
├── service                           # NRF service initialization and entry point
│   └── init.go
├── Taskfile.yml
├── test
│   └── test_server2.go
├── test-mirror.txt
├── util                              # Utility functions and test helpers
│   ├── util.go
│   └── util_test.go
├── VERSION
└── VERSION.license

18 directories, 56 files
```

## Configuration and Deployment

**Docker**

To build the container image:
```
task mod-start
task build
task docker-build-fast
```

**Kubernetes**

The standard deployment uses Helm charts from the Aether project. The version of the Chart can be found in the OnRamp repository in the `vars/main.yml` file.


## Quick Navigation


| Directory      | Description                                                                        |
| -------------- | ---------------------------------------------------------------------------------- |
| `management/`  | Implements NF management APIs for registration, status updates, and subscriptions. |
| `discovery/`   | Handles NF discovery requests and filtering logic.                                 |
| `accesstoken/` | Provides OAuth2-based access token handling for secure NF communication.           |
| `context/`     | Manages NF instance context, registration data, and discovery state.               |
| `factory/`     | Loads and manages NRF configuration files.                                         |
| `producer/`    | Core implementation of NF registration, discovery, and access token generation.    |
| `metrics/`     | Exports telemetry data to observability platforms.                                 |
| `polling/`     | Manages periodic NF configuration synchronization.                                 |
| `logger/`      | Centralized logging module for structured output.                                  |
| `dbadapter/`   | Abstracts database interactions for storing NF registration data.                  |




## NRF block diagram
![NRF Block Diagram](/docs/images/README-NRF.png)

## Supported Features
- Registration of Network Functions
- Searching of matching Network functions
- Handling multiple instances of registration from Network Functions
- Supporting keepalive functionality to check the health of network functions


## Upcoming changes in NRF
- Supporting callbacks to send notification when a network function is added/removed/modified.
- Subscription management callbacks to network functions.
- NRF cache library which can be used by modules to avoid frequent queries to NRF

Compliance of the 5G Network functions can be found at [5G Compliance](https://docs.sd-core.opennetworking.org/main/overview/3gpp-compliance-5g.html)

## Dynamic Network configuration (via webconsole)

NRF fetches the latest PLMN configuration from webconsole whenever a network function registers without providing
a list of supported PLMNs.
If a network function does not provide a list of supported PLMNs and NRF is not able to fetch any PLMN from webconsole (or
webconsole is unreachable), registration fails.
If a network function provides a list of supported PLMNs, it is registered without NRF fetching the configuration from webconsole.

### Setting Up PLMN configuration fetch

Include the `webuiUri` of the webconsole in the configuration file
```
configuration:
  ...
  webuiUri: https://webui:5001 # or http://webui:5001
  ...
```
The scheme (http:// or https://) must be explicitly specified.

## Reach out to us through

1. #sdcore-dev channel in [ONF Community Slack](https://aether5g-project.slack.com/)
2. Raise Github issues
