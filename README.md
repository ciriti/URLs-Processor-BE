[![Lint](https://github.com/ciriti/URLs-Processor-BE/actions/workflows/lint.yaml/badge.svg)](https://github.com/ciriti/URLs-Processor-BE/actions/workflows/lint.yaml)
[![Run Go Tests](https://github.com/ciriti/URLs-Processor-BE/actions/workflows/go-tests.yml/badge.svg)](https://github.com/ciriti/URLs-Processor-BE/actions/workflows/go-tests.yml)

# URLs Processor Backend

**URLs Processor Backend** is the server-side implementation of the URLs Processor application. This Go-based backend manages URL processing tasks, including authentication, URL management, and task queue handling. It interacts with the frontend to provide a seamless experience for users who need to manage and analyze URLs.

## Table of Contents

1. [Installation](#installation)
2. [Running the Application](#running-the-application)
3. [Project Structure](#project-structure)
4. [API Endpoints](#api-endpoints)
   - [Public Endpoints](#public-endpoints)
   - [Protected Endpoints](#protected-endpoints)
5. [Contributing](#contributing)
6. [License](#license)

## Installation

1. Clone the repository:

   ```sh
   git clone https://github.com/ciriti/URLs-Processor-BE.git
   cd URLs-Processor-BE
   ```

2. Set up environment variables. Create a `.env` file in the root directory of the project and add the following environment variables:

   ```env
   JWT_SECRET=your_jwt_secret
   ALLOWED_ORIGIN=http://your-frontend-domain.com
   PORT=8080
   WORKER_COUNT=5
   ```

3. Load the environment variables and dependencies:

   ```sh
   go mod tidy
   ```

## Running the Application

1. Build and run the Go server:

   ```sh
   go run ./cmd/api
   ```

   or

   ```sh
   go build -o myserver ./cmd/api
   ./myserver
   ```

2. The application will start on the port specified in the `.env` file (default is 8080). You can access it via `http://localhost:8080`.

## Project Structure

The project is organized into several key components:

- **cmd/api**: Contains the entry point for the application, including routing and handler functions.
- **internal/auth**: Manages authentication logic, including JWT generation and validation.
- **internal/middleware**: Provides middleware for the application, such as CORS handling.
- **internal/services**: Contains the core services, such as URL management, task queue handling, and page analysis.
- **internal/utils**: Utility functions for the application, including environment variable loading and graceful shutdown.

## API Endpoints

### Public Endpoints

- **POST /authenticate**: Authenticates the user and returns a JWT token.

### Protected Endpoints (require JWT token)

- **GET /logout**: Logs out the user by invalidating the JWT token.
- **GET /api/urls**: Retrieves all processed URLs.
- **POST /api/urls**: Adds URLs for processing.
- **GET /api/url**: Retrieves information about a specific URL.
- **POST /api/start**: Starts the computation for a specific URL.
- **POST /api/stop**: Stops the computation for a specific URL.

## Contributing

Contributions are welcome! Please read the contributing guidelines to get started.

## License

This project is licensed under the MIT License. See the LICENSE file for more details.
