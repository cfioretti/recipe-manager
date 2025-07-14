# Recipe-Manager Service

**Recipe-Manager Service** - Part of PizzaMaker Microservices Architecture: A microservice for recipe management and aggregation based on Domain-Driven Design (DDD) architecture with complete observability and monitoring.

## Main Features

- **Recipe Management**: Store and retrieve pizza recipes with ingredients and steps
- **Recipe Aggregation**: Combine recipe data with calculated ingredients and balanced portions
- **Multi-Service Integration**: Orchestrates calls to Calculator and Ingredients-Balancer services
- **Database Operations**: MySQL-based recipe storage and retrieval
- **Business Metrics**: Collects domain-specific metrics (recipe operations, service calls, database performance)

## Technologies

- **Go** - Primary language
- **Gin** - HTTP web framework
- **MySQL** - Database storage
- **gRPC** - External service communication
- **Prometheus** - Metrics and monitoring
- **OpenTelemetry + Jaeger** - Distributed tracing
- **Logrus** - Structured logging
- **Docker** - Containerization

## Endpoints

### HTTP Endpoints
- **Port**: 8080 (configurable)
- `POST /recipes/:uuid/aggregate` - Aggregate recipe with calculated ingredients
- `GET /metrics` - Prometheus metrics
- `GET /health` - Health check

### External Service Integration
- **Calculator Service** (gRPC): Dough weight calculations
- **Ingredients-Balancer Service** (gRPC): Ingredient balancing and optimization

## Observability

### Structured Logging
- **Correlation ID** for cross-service request tracking
- **Structured JSON** for easy parsing
- **Configurable levels** (Debug, Info, Warn, Error)

### Distributed Tracing
- **OpenTelemetry** for instrumentation
- **Jaeger** for trace visualization
- **Automatic spans** for HTTP and gRPC operations

### Prometheus Metrics
The service exposes both **business** and **technical** metrics:

#### Business Metrics
- `recipe_manager_recipe_retrievals_total` - Total recipe retrievals by UUID
- `recipe_manager_recipe_aggregations_total` - Total recipe aggregations by type
- `recipe_manager_calculator_service_calls_total` - Calculator service calls
- `recipe_manager_balancer_service_calls_total` - Balancer service calls
- `recipe_manager_database_operations_total` - Database operations by type
- `recipe_manager_recipes_by_author_total` - Recipes count by author
- `recipe_manager_recipe_complexity` - Recipe complexity score
- `recipe_manager_ingredient_variations` - Ingredient variations per recipe

#### Technical Metrics
- `recipe_manager_http_requests_total` - Total HTTP requests
- `recipe_manager_http_request_duration_seconds` - HTTP request duration
- `recipe_manager_active_http_connections` - Active HTTP connections
- `recipe_manager_database_operation_duration_seconds` - Database operation duration

## Business Logic

### Recipe Aggregation Flow
1. **Receive Request**: `POST /recipes/:uuid/aggregate` with pan specifications
2. **Calculate Dough**: Call Calculator service for dough weight calculations
3. **Retrieve Recipe**: Fetch recipe data from MySQL database
4. **Balance Ingredients**: Call Ingredients-Balancer for optimal ingredient distribution
5. **Return Aggregate**: Combined recipe with calculated and balanced ingredients

### Service Dependencies
- **Calculator Service**: `TotalDoughWeightByPans(context, Pans) -> Pans`
- **Ingredients-Balancer Service**: `Balance(context, Recipe, Pans) -> RecipeAggregate`

### Database Schema
The service uses MySQL for recipe storage with the following main entities:
- **recipes**: Core recipe information
- **ingredients**: Recipe ingredients with quantities
- **steps**: Recipe preparation steps
- **pans**: Pan specifications for calculations
