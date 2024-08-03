[![Lint](https://github.com/ciriti/URLs-Processor-BE/actions/workflows/lint.yaml/badge.svg)](https://github.com/ciriti/URLs-Processor-BE/actions/workflows/lint.yaml)
[![Run Go Tests](https://github.com/ciriti/URLs-Processor-BE/actions/workflows/go-tests.yml/badge.svg)](https://github.com/ciriti/URLs-Processor-BE/actions/workflows/go-tests.yml)

# URLs Processor

## Installation

1. Clone the repository:

   ```sh
   git clone https://github.com/yourusername/urls-processor.git
   cd urls-processor
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
   go run main.go
   ```

2. The application will start on the port specified in the `.env` file (default is 8080). You can access it via `http://localhost:8080`.

## API Endpoints

### Public Endpoints

- `GET /`: Home endpoint to check the status of the application.
- `POST /authenticate`: Authenticate the user and receive a JWT token.

### Protected Endpoints (require JWT token)

- `GET /api/urls`: Get all processed URLs.
- `POST /api/urls`: Add URLs for processing.
- `GET /api/url?id={id}`: Get information about a specific URL.
- `POST /api/start`: Start the computation for a specific URL.
- `POST /api/stop`: Stop the computation for a specific URL.

## Example .env File

Here is an example of what your `.env` file might look like:

```env
JWT_SECRET=mysecretkey
ALLOWED_ORIGIN=http://localhost:3000
PORT=8080
WORKER_COUNT=5
```
