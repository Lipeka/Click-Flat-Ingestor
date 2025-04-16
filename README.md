# ClickHouse Data Ingestion Tool

## Overview

This web application provides a complete solution for bidirectional data transfer between ClickHouse databases and CSV flat files. The tool features a responsive user interface with secure JWT authentication, efficient data processing, and comprehensive progress tracking.

## Key Features

### Core Functionality
- **Bidirectional Data Transfer**:
  - ClickHouse → CSV file export
  - CSV file → ClickHouse import
- **Multiple Source Support**:
  - ClickHouse database connections
  - Local CSV file uploads
- **Flexible Configuration**:
  - Custom column selection
  - Configurable CSV delimiters
  - Table joins for ClickHouse sources

### Security Features
- **JWT Authentication**:
  - 72-hour token expiration
  - Secure token signing with HMAC
- **Middleware Protection**:
  - All operations require valid JWT
  - Secure header parsing
- **CORS Configuration**:
  - Whitelisted origins
  - Restricted HTTP methods

### User Interface
- **Intuitive Workflow**:
  - Source selection
  - Connection configuration
  - Data preview
  - Ingestion execution
- **Responsive Design**:
  - Works on desktop and mobile
  - Adaptive layout
- **Visual Feedback**:
  - Progress tracking
  - Status notifications
  - Result summaries

## Technical Architecture

### System Components

```
Frontend (HTML/JS)           Backend (Go)
┌──────────────────────┐     ┌──────────────────────┐
│                      │     │                      │
│  User Interface      │◄───►│  API Endpoints      │
│  - Source selection  │     │  - /ingest          │
│  - Column selection  │     │  - /preview         │
│  - Progress tracking │     │  - /schema          │
│                      │     │  - /tables          │
└──────────────────────┘     └──────────┬──────────┘
                                        │
                                        ▼
                                ┌──────────────────────┐
                                │  ClickHouse Database │
                                └──────────────────────┘
```

### Backend Structure

```
.
├── handlers/        # HTTP request handlers
├── models/          # Data structures
├── services/        # Business logic
├── utils/           # Shared utilities
└── main.go          # Application entry point
```

## Installation & Setup

### Requirements
- Modern web browser (Chrome, Firefox, Safari, Edge)
- ClickHouse server access
- Go 1.16+ (for backend)
- Node.js (for frontend development)

### Configuration
Set these environment variables:
```bash
export JWT_SECRET_KEY="your-secret-key"
export CLICKHOUSE_DSN="clickhouse://user:password@host:port/database"
```

## Usage Instructions

### 1. Connect to Data Source
- Select source type (ClickHouse or CSV)
- For ClickHouse:
  - Enter host, port, database, username
  - Provide JWT token
- For CSV:
  - Upload file
  - Specify delimiter

### 2. Select Data
- For ClickHouse: choose table
- Select specific columns
- Preview data (optional)

### 3. Configure Transfer
- Choose direction:
  - ClickHouse → CSV
  - CSV → ClickHouse
- Set output parameters:
  - CSV filename
  - Target table name

### 4. Execute Transfer
- Monitor progress
- View completion summary

## API Reference

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/api/hello` | GET | Test endpoint |
| `/ingest` | POST | Main ingestion endpoint |
| `/preview` | GET | Data preview |
| `/schema` | POST | Table schema |
| `/tables` | POST | List tables |

## Example Requests

**ClickHouse Connection:**
```http
POST /tables
Authorization: Bearer <JWT_TOKEN>
Content-Type: application/json

{
  "source_type": "clickhouse",
  "host": "localhost",
  "port": 9000,
  "database": "default",
  "username": "user",
  "jwt_token": "<JWT_TOKEN>"
}
```

**CSV Ingestion:**
```http
POST /ingest
Authorization: Bearer <JWT_TOKEN>
Content-Type: multipart/form-data

{
  "source_type": "flatfile",
  "target_direction": "flatfile_to_clickhouse",
  "selected_tables": ["target_table"],
  "delimiter": ",",
  "jwt_token": "<JWT_TOKEN>"
}
```

## Advanced Features

### Multi-Table Joins
1. Select multiple ClickHouse tables
2. Configure JOIN conditions:
   - Join type (INNER, LEFT, RIGHT)
   - Left/right tables and columns
3. Execute joined query

### Data Preview
- View first 100 rows
- Verify column selection
- Check data formatting

## Best Practices

### Security
1. Rotate JWT secrets regularly
2. Limit file upload sizes
3. Use HTTPS in production
4. Restrict database permissions

### Performance
1. Use batch processing for large files
2. Monitor connection pooling
3. Optimize ClickHouse queries

## Troubleshooting


| Error Type | Status Code | Resolution |
|------------|-------------|------------|
| Authentication | 401 | Check JWT token |
| Validation | 400 | Verify request format |
| Database | 500 | Check connection details |
| File Processing | 400/500 | Verify file format/size |
