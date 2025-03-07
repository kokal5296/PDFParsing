# PDFParsing

## Project Description

PDFParsing is a project designed to parse PDF files and store the parsed data in a PostgreSQL database. The project provides a web server that allows users to upload PDF files, which are then processed and stored in the database. The project also includes functionality for managing users and their uploaded files.

## Purpose

The purpose of this project is to provide a simple and efficient way to parse and store PDF files. It is designed to be easy to set up and use, with a focus on providing a robust and reliable solution for managing PDF files.

## How to Use

1. Clone the repository:
   ```sh
   git clone https://github.com/kokal5296/PDFParsing.git
   cd PDFParsing
   ```

2. Set up the environment variables:
   Create a `.env` file in the root directory of the project and add the following variables:
   ```sh
   DB_USER=your_db_user
   DB_PASSWORD=your_db_password
   DB_NAME=your_db_name
   POSTGRESQL_URI=your_postgresql_uri
   PORT=your_port
   ```

3. Build and run the project using Docker:
   ```sh
   docker-compose up --build
   ```

4. The server will be running on the specified port. You can now upload PDF files and manage users through the provided API endpoints.

## Project Structure

The project is structured as follows:

- `database/`: Contains the code for managing the PostgreSQL database connection and operations.
- `error/`: Contains error handling code.
- `models/`: Contains the data models used in the project.
- `service/`: Contains the business logic for handling file uploads, user management, and queue management.
- `web/`: Contains the web server code and API handlers.
- `Dockerfile`: Defines the Docker image for the project.
- `docker-compose.yml`: Defines the Docker services for the project.
- `go.mod` and `go.sum`: Define the Go module dependencies.

## Main Components

- **Database Service**: Manages the connection to the PostgreSQL database and performs database operations.
- **File Service**: Handles file uploads, deletions, and imports.
- **Queue Service**: Manages the queue of files waiting to be processed.
- **User Service**: Manages user creation and retrieval of user files.
- **Web Server**: Provides the API endpoints for interacting with the project.

## Examples and Usage Scenarios

### Uploading a PDF File

To upload a PDF file, send a POST request to the `/file/:id` endpoint with the file and user ID as parameters. The file will be uploaded, processed, and stored in the database.

### Deleting a PDF File

To delete a PDF file, send a DELETE request to the `/file/:user_id/file_id/delete` endpoint with the user ID and file ID as parameters. The file will be deleted from the database.

### Importing a PDF File

To import a PDF file, send a POST request to the `/file/:user_id/:file_id/import` endpoint with the user ID and file ID as parameters. The file will be imported and its status updated in the database.

### Retrieving User Files

To retrieve all files uploaded by a user, send a GET request to the `/user/:id` endpoint with the user ID as a parameter. The response will contain a list of files uploaded by the user.
