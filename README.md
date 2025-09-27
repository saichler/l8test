# Layer 8 Quality

Comprehensive quality assurance and testing infrastructure for the Layer 8 ecosystem. This project provides professional-grade testing capabilities for both unit and integration testing, featuring advanced service orchestration, distributed testing topology, and comprehensive quality metrics.

## Overview

Layer 8 Quality is a state-of-the-art Go-based testing framework engineered to ensure the highest quality standards for Layer 8 networking infrastructure. It delivers a complete quality assurance environment with distributed service testing, transaction validation, performance benchmarking, and system integrity verification.

## Features

- **Test Topology Management**: Create and manage test network topologies with multiple nodes
- **Service Testing**: Comprehensive service handlers for testing REST operations (GET, POST, PUT, PATCH, DELETE)
- **Plugin System**: Extensible plugin architecture for registry and service components
- **Transaction Support**: Transaction management with replication capabilities
- **Distributed Caching**: Built-in distributed cache testing
- **Security Testing**: Integrated security testing framework
- **System Information**: Machine fingerprinting and environment detection for reproducible testing
- **Coverage Reporting**: Built-in test coverage reporting with HTML output
- **Performance Metrics**: Real-time performance monitoring and benchmarking

## Project Structure

```
go/
├── infra/
│   ├── sysinfo/           # System information and environment detection
│   ├── t_plugin/          # Plugin system components
│   │   ├── registry/      # Registry plugin for service discovery
│   │   └── service/       # Service plugin implementation
│   ├── t_resources/       # Resource management for tests
│   ├── t_service/         # Core test service implementations
│   └── t_topology/        # Network topology management
├── tests/                 # Test implementations
│   ├── TestInit.go        # Test initialization and setup
│   └── *_test.go          # Unit test files
├── test.sh               # Main test runner script
├── cover.html            # Coverage report output
└── go.mod                # Go module dependencies
```

## Dependencies

This project depends on several Layer 8 ecosystem modules:

- `github.com/saichler/l8services` - Service management and distributed caching
- `github.com/saichler/l8types` - Core types and interfaces
- `github.com/saichler/l8srlz` - Serialization framework
- `github.com/saichler/l8utils` - Utility functions and helpers
- `github.com/saichler/layer8` - Core Layer 8 networking
- `github.com/saichler/reflect` - Reflection utilities

## Getting Started

### Prerequisites

- Go 1.23.8 or later
- Access to the Layer 8 ecosystem repositories

### Running Tests

Execute the test suite using the provided test script:

```bash
cd go/
./test.sh
```

This script will:
1. Initialize Go modules and fetch dependencies
2. Set up security testing environment
3. Run unit tests with coverage analysis
4. Generate HTML coverage report
5. Open coverage report in browser

### Manual Testing

For more granular control, you can run tests manually:

```bash
cd go/
go mod init
GOPROXY=direct GOPRIVATE=github.com go mod tidy
go mod vendor
go test -tags=unit -v -coverpkg=./infra/... -coverprofile=cover.html ./...
```

## Quality Assurance Components

### System Information Module

The `sysinfo` module provides machine fingerprinting capabilities for ensuring test reproducibility across different environments. It generates unique machine hashes based on:
- Hostname and network interfaces
- Operating system and architecture
- CPU configuration
- Network MAC addresses

### Test Services

The framework provides three types of service handlers:

1. **TestServiceHandler**: Basic service operations without transactions
2. **TestServiceTransactionHandler**: Service operations with transaction support
3. **TestServiceReplicationHandler**: Service operations with replication and distributed caching

### Test Topology

The test topology manager (`TestTopology`) creates isolated network environments for testing:

- Configurable number of nodes
- Port allocation for services
- Network isolation and cleanup
- Logging and monitoring

### Plugin System

Extensible plugin architecture supports:

- **Registry Plugins**: Service discovery and registration
- **Service Plugins**: Custom service implementations
- Dynamic loading and unloading
- Build scripts for plugin compilation

## Configuration

### Test Parameters

Key configuration options:

- **Node Count**: Number of test nodes in topology (default: 4)
- **Port Ranges**: Service port allocations [20000, 30000, 40000]
- **Log Levels**: Configurable logging (Trace, Debug, Info, etc.)
- **Replication Count**: Number of replicas for distributed services (default: 2)

### Environment Variables

The test framework respects standard Go environment variables:

- `GOPROXY`: Go module proxy configuration
- `GOPRIVATE`: Private module configuration

## Recent Updates

### Version 2025.01
- Added system information module for environment detection
- Enhanced test reproducibility with machine fingerprinting
- Improved coverage reporting and metrics
- Updated dependencies to latest versions
- Performance optimizations in test topology management

## Quality Standards

Layer 8 Quality enforces strict quality standards:
- Minimum 80% test coverage for all components
- All tests must pass before deployment
- Performance regression testing
- Security vulnerability scanning
- Code quality metrics and linting

## Contributing

1. Ensure all tests pass before submitting changes
2. Add appropriate test coverage for new features (minimum 80%)
3. Follow Go best practices and coding standards
4. Update documentation as needed
5. Run performance benchmarks for critical paths

## License

This project is part of the Layer 8 ecosystem. See the LICENSE file for details.

---

**Layer 8 Quality** - Ensuring Excellence in Network Infrastructure Testing
