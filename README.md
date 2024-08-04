[![Lint](https://github.com/ciriti/URLs-Processor-BE/actions/workflows/lint.yaml/badge.svg)](https://github.com/ciriti/URLs-Processor-BE/actions/workflows/lint.yaml)
[![Run Go Tests](https://github.com/ciriti/URLs-Processor-BE/actions/workflows/go-tests.yml/badge.svg)](https://github.com/ciriti/URLs-Processor-BE/actions/workflows/go-tests.yml)

# URLs Processor

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

2. The application will start on the port specified in the `.env` file (default is 8080). You can access it via `http://localhost:8080`.

## API Endpoints

### Public Endpoints

#### `POST /authenticate`

**Description:** Authenticate the user and receive a JWT token.

**Request**

- **Headers:**
  - `Content-Type`: `application/json`
- **Body:**
  - `user` (string): Username
  - `pass` (string): Password

**Response**

- **200 OK**
  - **Fields:**
    - `status` (string): "success"
    - `token` (string): JWT token
- **400 Bad Request**
  - **Fields:**
    - `status` (string): "error"
    - `message` (string): Error message
- **401 Unauthorized**
  - **Fields:**
    - `status` (string): "error"
    - `message` (string): "Unauthorized"

### Protected Endpoints (require JWT token)

#### `GET /logout`

**Description:** Logout the user by invalidating the JWT token.

**Request**

- **Headers:**
  - `Content-Type`: `application/json`

**Response**

- **200 OK**
  - **Fields:**
    - `message` (string): "logout successful"
- **500 Internal Server Error**
  - **Fields:**
    - `status` (string): "error"
    - `message` (string): "Internal Server Error"

### Protected Endpoints (require JWT token)

#### `GET /api/urls`

**Description:** Get all processed URLs.

**Request**

- **Headers:**
  - `Authorization`: `Bearer {token}`

**Response**

- **200 OK**
  - **Fields:**
    - `status` (string): "success"
    - `data` (array of objects): List of processed URLs
- **401 Unauthorized**
  - **Fields:**
    - `status` (string): "error"
    - `message` (string): "Unauthorized"

#### `POST /api/urls`

**Description:** Add URLs for processing.

**Request**

- **Headers:**
  - `Authorization`: `Bearer {token}`
  - `Content-Type`: `application/json`
- **Body:**
  - `urls` (array of strings): List of URLs to be processed

**Response**

- **200 OK**
  - **Fields:**
    - `status` (string): "success"
    - `failed` (array of strings): List of URLs that failed to be processed
- **400 Bad Request**
  - **Fields:**
    - `status` (string): "error"
    - `message` (string): Error message
- **401 Unauthorized**
  - **Fields:**
    - `status` (string): "error"
    - `message` (string): "Unauthorized"

#### `GET /api/url`

**Description:** Get information about a specific URL.

**Request**

- **Headers:**
  - `Authorization`: `Bearer {token}`
- **Query Parameters:**
  - `id` (int): The ID of the URL to retrieve information for

**Response**

- **200 OK**
  - **Fields:**
    - `status` (string): "success"
    - `data` (object): URL information
- **400 Bad Request**
  - **Fields:**
    - `status` (string): "error"
    - `message` (string): "Invalid URL ID"
- **401 Unauthorized**
  - **Fields:**
    - `status` (string): "error"
    - `message` (string): "Unauthorized"
- **404 Not Found**
  - **Fields:**
    - `status` (string): "error"
    - `message` (string): "URL not found"

#### `POST /api/start`

**Description:** Start the computation for a specific URL.

**Request**

- **Headers:**
  - `Authorization`: `Bearer {token}`
  - `Content-Type`: `application/json`
- **Body:**
  - `id` (int): The ID of the URL to start processing

**Response**

- **200 OK**
  - **Fields:**
    - `status` (string): "success"
    - `id` (int): The ID of the URL
    - `state` (string): "pending"
- **400 Bad Request**
  - **Fields:**
    - `status` (string): "error"
    - `message` (string): "Invalid request payload"
- **401 Unauthorized**
  - **Fields:**
    - `status` (string): "error"
    - `message` (string): "Unauthorized"
- **404 Not Found**
  - **Fields:**
    - `status` (string): "error"
    - `message` (string): "URL not found"

#### `POST /api/stop`

**Description:** Stop the computation for a specific URL.

**Request**

- **Headers:**
  - `Authorization`: `Bearer {token}`
  - `Content-Type`: `application/json`
- **Body:**
  - `id` (int): The ID of the URL to stop processing

**Response**

- **200 OK**
  - **Fields:**
    - `status` (string): "success"
    - `id` (int): The ID of the URL
    - `state` (string): "stopped"
- **400 Bad Request**
  - **Fields:**
    - `status` (string): "error"
    - `message` (string): "Invalid request payload"
- **401 Unauthorized**
  - **Fields:**
    - `status` (string): "error"
    - `message` (string): "Unauthorized"
- **404 Not Found**
  - **Fields:**
    - `status` (string): "error"
    - `message` (string): "URL not found"
