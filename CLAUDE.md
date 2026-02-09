# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

CryptAply is an enterprise-grade cryptographic key management and governance platform. It provides secure team-based key operations with multi-layer approval workflows, comprehensive audit trails, and integration with cloud HSM providers (Azure Key Vault, AWS KMS, Google Cloud KMS).

**Technology Stack:**
- Backend: C# / .NET 8.0, ASP.NET Core, Entity Framework Core 6.0
- Database: PostgreSQL (primary)
- Frontend Web: Next.js 15.5.3 (React 19), TypeScript, TailwindCSS
- Mobile: .NET MAUI (iOS/Android targeting net8.0)
- Infrastructure: Docker, Kubernetes, RabbitMQ, Serilog + Graylog

## Architecture

### Clean/Onion Architecture (Layered)

```
Controllers → Application Layer → Infrastructure Layer → Domain Layer
```

**Domain Layer** (`CryptAply.Domain`):
- Core entities, interfaces, enums, and value objects
- Domain-driven design with aggregate roots
- Key aggregates: Application, Audit, Ceremony, Cloud, Device, Organization, KeyManagement, Compliance, Emergency, Identity

**Application Layer** (`CryptAply.Api.Application`):
- Business logic and service implementations
- FluentValidation validators
- Service interfaces and models
- Dependency injection configuration

**Infrastructure Layer** (`CryptAply.Api.Infrastructure`):
- Entity Framework DbContext implementations
- Repository pattern implementations (IBaseRepository, IUnitOfWork)
- Database migrations
- Data access layer organized by domain context

**API Layer** (`CryptAply.Api`, `CryptAply.PublicApi`):
- REST API controllers (versioned V1/)
- Middleware for rate limiting, exception handling, timing
- Swagger/OpenAPI documentation
- OAuth 2.0 + DPoP support

### Authentication & Authorization

**Current architecture (as of Oct 2024 commits):**
- Primary identity provider: **PayEz IDP** (external)
- JWT Bearer tokens issued by PayEz (CryptAply no longer issues JWTs)
- Secrets managed via Azure Key Vault (loaded at startup)
- Scope-based authorization: `cryptaply:keys:read`, `cryptaply:keys:manage`, `cryptaply:keys:rotate`, `cryptaply:org:info`
- Role-based access control: Owner, Admin, SecurityOfficer, ComplianceOfficer, Auditor, TechnicalLead, Member, Viewer

**DPoP Support:**
- Demonstration of Proof-of-Possession for enhanced token security
- Custom implementation in `CryptAply.PublicApi`

### Repository Pattern

All data access uses the repository pattern with:
- `IBaseRepository<T>` for CRUD operations
- `IUnitOfWork` for transaction management
- Aggregate root interfaces: `IAggregateRoot`, `IOrganizationalAggregate`, `ISecurityAggregate`, `IWorkflowAggregate`
- Specification pattern for complex queries
- Domain-specific repositories in `CryptAply.Api.Infrastructure/Repositories/`

### Key Domain Contexts

1. **Organization**: Teams, members, departments, roles
2. **KeyManagement**: Encryption keys, rotation policies, metadata
3. **Ceremony**: Key ceremony workflows with quorum-based voting
4. **Audit**: Compliance logging (AuditLog, KeyOperationLog, AccessControlLog)
5. **Device**: HSM device management and monitoring
6. **Cloud**: Cloud provider integrations (Azure, AWS, GCP)
7. **Compliance**: Policy enforcement, frameworks, violation tracking
8. **Emergency**: Emergency protocols and incident response
9. **Workflow**: State machines for key operations

## Common Development Commands

### Backend API

```bash
# Restore dependencies
dotnet restore CryptAply.sln

# Build entire solution
dotnet build CryptAply.sln

# Build specific project
dotnet build Api/CryptAply.Api/CryptAply.Api.csproj

# Run main API locally (development)
dotnet run --project Api/CryptAply.Api/CryptAply.Api.csproj

# Run PublicApi (partner integration API)
dotnet run --project Api/CryptAply.PublicApi/CryptAply.PublicApi.csproj

# Run tests (if test projects exist)
dotnet test

# Entity Framework migrations (from Infrastructure project directory)
cd Api/CryptAply.Api.Infrastructure
dotnet ef migrations add MigrationName --startup-project ../CryptAply.Api/CryptAply.Api.csproj
dotnet ef database update --startup-project ../CryptAply.Api/CryptAply.Api.csproj
```

### Frontend Web (Next.js)

```bash
cd website/cryptaply-web

# Install dependencies
npm install

# Run development server (http://localhost:3000)
npm run dev

# Build for production
npm run build

# Start production server
npm start

# Run linter
npm run lint
```

### Mobile App (MAUI)

```bash
# Build for iOS
dotnet build -f net8.0-ios App/CryptAply.App/CryptAply.App.csproj

# Build for Android
dotnet build -f net8.0-android App/CryptAply.App/CryptAply.App.csproj

# Run on iOS simulator (macOS only)
dotnet build -f net8.0-ios -t:Run App/CryptAply.App/CryptAply.App.csproj

# Run on Android emulator
dotnet build -f net8.0-android -t:Run App/CryptAply.App/CryptAply.App.csproj
```

### Docker

```bash
# Build and run with docker-compose (local development)
docker-compose up --build

# Build specific service
docker build -t cryptaply-api -f docker/Dockerfile .

# Run container
docker run -p 5000:80 cryptaply-api
```

### Database

```bash
# Run PostgreSQL migrations
dotnet ef database update --project Api/CryptAply.Api.Infrastructure --startup-project Api/CryptAply.Api

# Generate migration script
dotnet ef migrations script --project Api/CryptAply.Api.Infrastructure --startup-project Api/CryptAply.Api --output Database/Scripts/migration.sql

# Database scripts are in Database/Scripts/ and Database/SQL/
```

## Project Structure

```
E:\repos\cryptaply/
├── Api/                          # Backend APIs
│   ├── CryptAply.Api/           # Main REST API
│   ├── CryptAply.Api.Application/   # Business logic layer
│   ├── CryptAply.Api.Infrastructure/ # Data access layer
│   ├── CryptAply.PublicApi/     # Partner integration API (OAuth)
│   └── CryptAply.Infrastructure/    # Shared infrastructure
│
├── App/                          # Mobile applications
│   ├── CryptAply.App/           # MAUI cross-platform app
│   ├── CryptAply.App.Application/   # Mobile business logic
│   └── CryptAply.App.Infrastructure/ # Mobile data access
│
├── Domain/                       # Domain models
│   ├── CryptAply.Domain/        # Entities, interfaces, aggregates
│   └── CryptAply.DTOs/          # Data Transfer Objects
│
├── website/cryptaply-web/       # Next.js web application
│
├── Database/                     # Database scripts and migrations
│   ├── Scripts/                 # SQL scripts, migration scripts
│   └── SQL/                     # Schema backups
│
├── docs/                        # Documentation
│   ├── CryptoTeamUserGuide.md
│   ├── CryptoTeamAdminGuide.md
│   ├── cryptaply-repository-list-and-patterns.md
│   └── backing-website-docs/
│
└── docker/                      # Docker configurations
```

## Important Files

**Entry Points:**
- `Api/CryptAply.Api/Program.cs` - Main API startup
- `Api/CryptAply.PublicApi/Program.cs` - PublicApi startup
- `website/cryptaply-web/app/page.tsx` - Next.js home page

**Configuration:**
- `Api/CryptAply.Api/appsettings.json` - API configuration
- `Api/CryptAply.PublicApi/appsettings.json` - PublicApi configuration
- `Directory.Build.props` - Global .NET SDK settings (8.0.100)
- `docker-compose.yml` - Local development orchestration

**Database:**
- `Api/CryptAply.Api.Infrastructure/Context/CryptAplyDbContext.cs` - Main EF DbContext
- `Api/CryptAply.Api.Infrastructure/Migrations/` - EF migrations

## Development Workflow

### Git Branching

- **Main branch**: `master`
- **Current development**: `ten-seventeen-stage`
- Use pull requests for all changes
- Follow conventional commit messages

### Recent Development Focus (Oct 2024)

1. **Authentication Refactoring**: Migrating JWT issuance to PayEz IDP
2. **Secret Management**: Loading all secrets from Azure Key Vault at startup
3. **OAuth/DPoP**: Enhanced PublicApi with OAuth 2.0 Client Credentials + DPoP
4. **Compliance**: Enhanced audit trail and reporting capabilities

### Code Style

- Follow .editorconfig settings in the root
- Use C# naming conventions (PascalCase for public members, camelCase for private)
- TypeScript/React: Follow ESLint configuration in Next.js project
- Comprehensive XML documentation comments for public APIs

## Security Considerations

1. **Multi-factor Authentication**: Enforced via PayEz IDP
2. **Quorum-Based Approvals**: Multi-person approval workflows for sensitive key operations
3. **Comprehensive Audit Logging**: All operations logged to AuditLog, KeyOperationLog, AccessControlLog
4. **HSM Integration**: Keys managed via Azure Key Vault, AWS KMS, Google Cloud KMS
5. **Emergency Protocols**: Override procedures with mandatory post-incident review
6. **Data Protection**: ASP.NET Core data protection API, TLS 1.2+ enforced

## Database Schema

**Core Tables:**
- Users, Teams, Organizations, Members (identity/organization)
- Keys, KeyRotation, KeyMetadata (key management)
- AuditLogs, KeyOperationLogs, AccessControlLogs (compliance)
- VoteSessions, Quorum, Certificates (ceremony/approval)
- HSMDevices, CloudProviders (infrastructure)
- Policies, ComplianceFrameworks (governance)

**EF Core Considerations:**
- Code-first migrations approach
- Multi-context architecture (separate DbContext per aggregate)
- Navigation properties for complex relationships
- Shadow properties for audit tracking (CreatedAt, UpdatedAt, CreatedBy, UpdatedBy)

## API Endpoints

### Main API (CryptAply.Api)

Controller namespaces in `Api/CryptAply.Api/Controllers/V1/`:
- `Application/` - Application management
- `Cloud/` - Cloud provider integrations
- `Device/` - HSM device management
- `KeyCeremony/` - Key rotation and ceremony workflows
- `Organization/` - Team and organization management

### PublicApi (Partner Integration)

- OAuth 2.0 endpoints: `/oauth/token`
- Partner key inventory: `/api/partner/keys`
- Organization info: `/api/partner/organization`
- Swagger UI available in production (configurable)

## Testing

- Test scripts in `Tests/` directory
- Integration tests in `Tests/EmergencyProtocols/`
- Development test scripts in `Tests/dev/`

## Monitoring & Logging

**Structured Logging (Serilog):**
- Console sink for development
- Graylog sink for production log aggregation
- Configured in `CryptAply.Api.Application/Configuration/SerilogSetup.cs`

**Custom Middleware:**
- Request timing middleware for performance metrics
- Exception handling middleware with structured error responses
- Rate limiting middleware for API protection

## Documentation

**User Documentation:**
- `Readme.md` - Getting started guide for team members
- `docs/CryptoTeamUserGuide.md` - Detailed user operations
- `docs/CryptoTeamAdminGuide.md` - Administrator procedures

**Technical Documentation:**
- `docs/cryptaply-repository-list-and-patterns.md` - Architecture deep-dive
- `docs/database-schema-organization.md` - Database design
- `docs/breaking-circular-dependencies-EF-considerations.md` - EF Core patterns
- `docs/backing-website-docs/` - Web platform documentation

**Strategic Documentation:**
- `docs/cryptaply-escrow-vision.md` - Product vision and roadmap
- `docs/Phase1-Reporting-Strategy-Julius.md` - Phase 1 reporting strategy

## Solution Files

- `CryptAply.sln` - Main development solution (all projects)
- `CryptAply-McrSrvs.sln` - Microservices solution variant
- `App/CryptAply-AppDevs.sln` - Mobile app development solution

## Key Implementation Patterns

### Service Registration

Services registered in `CryptAply.Api.Application/Extensions/` using extension methods:
- `AddApplicationServices()` - Business services
- `AddInfrastructureServices()` - Repositories, DbContext
- `AddAuthenticationServices()` - JWT, OAuth, DPoP configuration

### Validation

FluentValidation validators in `CryptAply.Api.Application/Validators/` for request DTOs.

### Dependency Injection

All services use constructor injection. Register in `Program.cs` or service collection extensions.

### Error Handling

- Custom exception handling middleware in PublicApi
- Structured error responses with appropriate HTTP status codes
- Comprehensive logging of all exceptions

## Performance Considerations

- Use async/await for all I/O operations
- Repository pattern supports IQueryable for deferred execution
- Specification pattern for complex queries
- Consider pagination for large result sets
- Entity Framework query optimization (include eager loading where appropriate)

## Cloud Provider Integration

**Azure Key Vault:**
- Connection setup via Azure SDK
- Secrets loaded at startup from Key Vault
- Certificate management

**AWS KMS:**
- IAM configuration required
- Key policy setup in AWS console

**Google Cloud KMS:**
- Service account authentication
- Project and location configuration
