# Uploading Files to S3 using Go and React

This is an example repo for uploading files to S3 using Go and React.

## Setup

### Adding in the .env variables

```bash
cp .env.example .env
```

Get all the necessary variables from the AWS console.

### Running the server

```bash
go mod download
go run main.go
```

### Running the frontend

```bash
cd frontend
npm install
npm run dev
```

The frontend will be running on `localhost:5173` and the backend will be running on `localhost:8080`.
