# Assetra Frontend

A modern and responsive React application for the Assetra system, built with development best practices.

## 🚀 Technologies

- **React 19** - JavaScript library for user interfaces
- **React Router Dom** - Routing for SPAs
- **Tailwind CSS** - Utility-first CSS framework
- **Axios** - HTTP client for API requests
- **React Hook Form** - Form management
- **React Hot Toast** - Elegant notifications
- **Lucide React** - Modern icons
- **JWT** - Token-based authentication

## 🎨 Features

### ✅ Implemented

- **Homepage** - Home page with system presentation
- **Authentication** - User login and registration
- **User Management** - List, view and management
- **Profile** - User profile page
- **Responsive Layout** - Adaptive design for all devices
- **Notification System** - Visual feedback for actions
- **Protected Routing** - Authentication-based access control
- **HTTP Interceptors** - Automatic token management

### 🎯 Security Features

- Secure JWT token storage
- Automatic authentication interceptors
- Form validation
- Route protection
- Automatic logout on expired token

## 🛠️ Setup

1. **Install dependencies:**

```bash
npm install
```

2. **Configure environment variables:**

```bash
cp .env.example .env
```

Edit the `.env` file with your backend settings:

```
REACT_APP_API_URL=http://localhost:9000
```

3. **Start development server:**

```bash
npm start
```

The application will be available at `http://localhost:3000`

### UI Components

- **Button**: Buttons with variants (primary, secondary, outline, ghost, danger)
- **Input**: Input fields with validation and error states
- **Card**: Content containers with header, content and footer
- **Loading**: Loading indicators for different contexts

## 🔐 Authentication

The system uses JWT (JSON Web Tokens) for authentication:

1. **Login**: User provides credentials
2. **Token**: Backend returns JWT token
3. **Storage**: Token stored in cookies

## 🚀 Deploy

To build the application:

```bash
npm run build
```

The optimized files will be generated in the `build/` folder

## 🔄 Backend Integration

The application connects to the following APIs:

- `POST /signin` - User login
- `POST /signup` - User registration
- `GET /users` - User list (protected)
- `GET /users/{id}` - Specific user (protected)
- `PUT /users/{id}` - Update user (protected)
- `DELETE /users/{id}` - Delete user (protected)
