# Kokomed Finance

Kokomed Finance is a comprehensive lending platform designed to support Small and Medium-sized Enterprises (SMEs) by streamlining loan management processes. The application offers features such as customer management, loan processing, user roles for loan officers, and detailed reporting capabilities.

## Features

-   **Customer Management**: Easily add and manage customer information.
-   **Loan Processing**: Efficiently handle loan applications and track their statuses.
-   **User Roles**: Assign and manage roles for loan officers and other users.
-   **Reporting**: Generate comprehensive reports to gain insights into lending activities.
-   **Email Notifications**: Async reminders and status updates

## Technologies Used

### Frontend

-   **React** + **TypeScript**
-   **shadcn/ui**: Modern component library
-   **TanStack Query**: Server state management
-   **TanStack Table**: Table state management and rendering
-   **Axios**: HTTP client
-   **Zod**: Schema validation

### Backend

-   **Golang**: High-performance API
-   **Gin Framework**: RESTful routing
-   **GORM**: MySQL ORM layer
-   **Redis**: Caching & rate limiting
-   **Asynq**: Background job processing

### Infrastructure

-   **MySQL**: Primary database
-   **Docker**: Containerization
-   **NGINX**: Reverse proxy

## Installation

To set up the Kokomed Finance project locally, follow these steps:

### Prerequisites

-   Go 1.21+
-   Node.js 18+
-   MySQL 8+
-   Redis 7+

### Clone the Repository

```bash
git clone https://github.com/EmilioCliff/kokomed-fin.git
cd kokomed-fin
```

### Backend Setup

-   Navigate to the backend directory:

    ```bash
    cd backend
    ```

-   Install dependencies:

    ```bash
    go mod download
    ```

-   Set up environment variables as required (e.g., database credentials, Redis configuration).

```bash
 cp ./.envs/.local/config.env ./.envs/.local/.env
```

-   Run the backend server:

    ```bash
    go run main.go
    ```

### Frontend Setup

-   Navigate to the frontend directory:

    ```bash
    cd ../frontend
    ```

-   Install dependencies:

    ```bash
    npm install
    ```

-   Start the frontend development server:

    ```bash
    npm run dev
    ```

### Access the Application

-   Open your browser and navigate to `http://localhost:3000` to access the Kokomed Finance application.

## Usage

Once the application is running, you can:

-   **Add Customers**: Navigate to the "Customers" section to add and manage customer details.
-   **Process Loans**: In the "Loans" section, initiate new loan applications and monitor their progress.
-   **Manage Users**: Assign roles and manage loan officers in the "Users" section.
-   **Generate Reports**: Access the "Reports" section to create and view detailed reports on lending activities.

## Work To Be Done

-   [x] Implement core lending functionalities
-   [x] Integrate Redis for caching
-   [x] Set up user authentication and role-based access control
-   [ ] Send SMS notification on loan status and deadlines
-   [ ] Write tests for Golang backend
-   [ ] Use Testcontainers for integrated testing
-   [ ] Implement Zerolog for structured logging
-   [ ] Add observability using OpenTelemetry

## License

This project is licensed under the MIT License. See the [Apache License](LICENSE) file for more details.
