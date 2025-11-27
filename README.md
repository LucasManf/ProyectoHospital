# Medical Appointments Management System

A simple medical appointment management system built with Go, PostgreSQL, and BoltDB. This CLI application allows managing patients, doctors, appointments, scheduling, and monthly billing for medical practices.

## Technologies Used

- **Go**: Command-line application for interacting with the database
- **PostgreSQL**: Relational database for storing medical appointments data
- **BoltDB (NoSQL)**: Alternative storage system using JSON-based key-value pairs

## Prerequisites

Before running this project, make sure you have the following installed:

- **Go** (version 1.16 or higher) - [Download Go](https://golang.org/dl/)
- **Docker** and **Docker Compose** - [Download Docker](https://www.docker.com/get-started)

## Project Structure

```
root
│
├── hospital.sql
├── docker-compose.yml
│
├── aplicacion-de-bases-de-datos/
│   └── main.go
│
└── hospital-json/
    ├── app-boltdb.go
    └── hospital.db
```

## Installation & Setup

### 1. Clone the Repository

```bash
git clone <your-repository-url>
cd <project-directory>
```

### 2. Start PostgreSQL with Docker

Start the PostgreSQL database using Docker Compose:

```bash
docker-compose up -d
```

This will start a PostgreSQL container on port `5432` with the following credentials:
- **User**: `admin`
- **Password**: `1234`
- **Database**: `postgres` (default)

To verify the container is running:

```bash
docker-compose ps
```

To stop the database:

```bash
docker-compose down
```

To stop and remove all data:

```bash
docker-compose down -v
```

### 3. Run the Application

Navigate to the application directory and run:

```bash
cd aplicacion-de-bases-de-datos
go run main.go
```

## Features

### 1. Patient Management
- Register patients with medical history number, name, date of birth, health insurance, and contact information

### 2. Doctor Management
- Register doctors with ID, specialty, consultation fees, and availability

### 3. Doctor Schedule Management
- Set up doctor availability for specific days and time slots, including appointment duration

### 4. Appointment Booking
- Patients can book appointments with available doctors
- System validates doctor availability, patient credentials, and scheduling conflicts

### 5. Appointment Cancellation & Rescheduling
- Cancel appointments and register them for rescheduling
- Notification system for affected patients

### 6. Appointment Attendance
- Mark appointments as attended
- Validates that appointments are reserved and scheduled for the current day

### 7. Health Insurance Billing
- Generate monthly billing statements for health insurance companies
- Detailed reports of services provided and amounts due

### 8. Automated Email Notifications
The system generates automatic email notifications for:
- Appointment confirmation
- Appointment cancellation
- Appointment reminders (for upcoming appointments)
- Missed appointment notifications

### 9. NoSQL Model
- Alternative storage using BoltDB for comparison between relational and non-relational approaches

### 10. CLI Interface
- Simple command-line interface for all system operations

## Database Features

The project includes stored procedures and triggers for:
- Automatic generation of available appointment slots
- Appointment reservation validation
- Appointment cancellation and rescheduling
- Attendance tracking
- Health insurance billing
- Automated email generation

## Docker Compose Configuration

The `docker-compose.yml` file is configured as follows:

```yaml
version: '3.8'
services:
  postgres:
    image: postgres:15
    environment:
      POSTGRES_USER: admin
      POSTGRES_PASSWORD: 1234
      POSTGRES_DB: postgres
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data

volumes:
  postgres_data:
```

## Usage

1. Start Docker container: `docker-compose up -d`
2. Run the application: `go run main.go`
3. Follow the CLI prompts to manage appointments
4. Stop the database when done: `docker-compose down`

## Notes

- The application will create the `hospital` database automatically on first run
- The timezone is set to `America/Argentina/Buenos_Aires` by default
- Data persists in Docker volumes between container restarts
