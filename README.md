# Library Management System

This project is a Library Management System designed for local libraries, providing features for users, agents, and administrators. The backend is implemented in Go language using the Gin web framework, PostgreSQL for data storage, Redis for caching, Cron for scheduled tasks, Razorpay for payment processing, SMTP for email notifications, Golint for code linting, and GORM as the ORM.

## Live Demo

Access the live demo of the Library Management System at [golib.online](https://golib.online).

## Features

## Changelog

### Version 2 (2023-11-26)

- Added email notifications for membership expiry and fines.
- Added return feature for books.
- Enhanced book viewing:
  - View books by category, author.
  - Sort by order count and rating.
- Added order management table.
- Added image field in the book table.
- Introduced membership invoice download feature.
- Admin routes for viewing books checked out, history, and orders.


### User Features

- **Authentication:**
  - User registration with OTP verification for enhanced security.
  - Login with credentials.

- **Book Management:**
  - View available books.
  - Search and view books.
  - Pagination for a better user experience.

- **Feedback:**
  - Provide feedback on books or the library.

- **Profile:**
  - View and update user profiles.
  - Change passwords securely.

- **Membership and Borrowing:**
  - Take membership to borrow books.
  - Book delivery by agents.

### Admin Features

- **User and Role Management:**
  - Manage users.
  - Role-based access control (Admin).

- **Library Management:**
  - CRUD operations on authors, categories, publications, and books.

- **Review and Events:**
  - Remove abusive/indecent reviews.
  - Add and view library events.

- **Feedback Management:**
  - View user feedback.

### Agent Features

- **Order Management:**
  - View pending orders.
  - Get order details and update order status.

- **Fine Calculation:**
  - Automatically calculate fines for overdue books.

## Technologies Used

- **Backend:**
  - Go language
  - Gin web framework
  - PostgreSQL for data storage
  - Redis for caching
  - Cron for scheduled tasks
  - Razorpay for payment processing
  - SMTP for email notifications
  - GORM as the ORM
  - Golint for code linting

## API Documentation

For detailed API documentation, refer to [API Documentation](https://documenter.getpostman.com/view/30219361/2s9YeEdCnd).


## Setup and Installation

1. Clone the repository:

   ```bash
   git clone https://github.com/your-username/library-management.git

2.Install dependencies:

    cd library-management
    go get -u github.com/gin-gonic/gin
    go get -u github.com/razorpay/razorpay-go
    go get -u github.com/go-redis/redis
    go get -u gorm.io/driver/postgres

3.Set up the database:

    CREATE DATABASE library_management

4.Configure environment variables:

    DB_Config="host=localhost user=##### password=***** dbname=library_management port=0000 sslmode=disable"  


    Email="youremail@email.com"
    Password="__ __ __ __(use app password)"
    
    
    RAZORPAY_KEY_ID="________________(apikey)"
    RAZORPAY_SECRET="_______________(api secret)"

5.Run the application:

    make run
