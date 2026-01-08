Project Structure
.
├── adapter/                # Adapter Layer - External system interactions
│   ├── amqp/               # Message queue adapters
│   ├── dependency/         # Dependency injection configuration
│   │   └── wire.go         # Wire DI setup with interface bindings
│   ├── job/                # Scheduled task adapters
│   └── repository/         # Data repository adapters
│       ├── mysql/          # MySQL implementation
│       │   └── entity/     # Database entities and repo implementations
│       ├── postgre/        # PostgreSQL implementation
│       ├── mongo/          # MongoDB implementation
│       └── redis/          # Redis implementation
│           └── enhanced_cache.go  # Enhanced cache with advanced features
├── api/                    # API Layer - HTTP requests and responses
│   ├── dto/                # Data Transfer Objects for API
│   ├── error_code/         # Error code definitions
│   ├── grpc/               # gRPC API handlers
│   ├── middleware/         # Global middleware including metrics collection
│   └── http/               # HTTP API handlers
│       ├── handle/         # Request handlers using domain interfaces
│       ├── middleware/     # HTTP middleware
│       ├── paginate/       # Pagination handling
│       └── validator/      # Request validation
├── application/            # Application Layer - Use cases coordinating domain objects
│   ├── core/               # Core interfaces and base implementations
│   │   └── interfaces.go   # UseCase and UseCaseHandler interfaces
│   └── example/            # Example use case implementations
│       ├── create_example.go     # Create example use case
│       ├── delete_example.go     # Delete example use case
│       ├── get_example.go        # Get example use case
│       ├── update_example.go     # Update example use case
│       └── find_example_by_name.go # Find example by name use case
├── cmd/                    # Command-line entry points
│   └── main.go             # Main application entry point
├── config/                 # Configuration management
│   ├── config.go           # Configuration structure and loading
│   └── config.yaml         # Configuration file
├── domain/                 # Domain Layer - Core business logic
│   ├── aggregate/          # Domain aggregates
│   ├── dto/                # Domain Data Transfer Objects
│   ├── event/              # Domain events
│   ├── model/              # Domain models
│   ├── repo/               # Repository interfaces
│   ├── service/            # Domain services
│   └── vo/                 # Value Objects
└── tests/                  # Test utilities and examples
    ├── migrations/         # Database migrations for testing
    ├── mysql.go            # MySQL test utilities
    ├── postgresql.go       # PostgreSQL test utilities
    └── redis.go            # Redis test utilities
