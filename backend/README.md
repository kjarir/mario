# Dr. Mario Backend API

A comprehensive Go backend API for diabetic retinopathy detection through retinal imaging analysis.

## 🏥 Features

- **User Authentication & Authorization**: JWT-based authentication with role-based access control
- **Patient Management**: Complete patient profile and medical history management
- **Doctor Management**: Doctor profiles, specializations, and credentials
- **Image Upload & Processing**: Secure retinal image upload with AI detection
- **AI Detection**: Diabetic retinopathy detection with confidence scoring
- **Appointment Scheduling**: Patient-doctor appointment management
- **Analytics & Reporting**: System statistics and detection analytics
- **File Management**: Secure image storage and retrieval
- **In-Memory Storage**: Fast, stateless data storage (no database required)

## 🚀 Quick Start

### Prerequisites

- Go 1.21 or higher
- Git

### Installation

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd backend
   ```

2. **Install dependencies**
   ```bash
   go mod tidy
   ```

3. **Set up environment variables**
   ```bash
   cp env.example .env
   # Edit .env with your configuration
   ```

4. **Run the application**
   ```bash
   go run main.go
   ```

The server will start on `http://localhost:8080`

## 📋 API Endpoints

### Authentication
- `POST /api/v1/auth/register` - User registration
- `POST /api/v1/auth/login` - User login

### User Profile
- `GET /api/v1/profile` - Get current user profile
- `PUT /api/v1/profile` - Update user profile

### Patients
- `GET /api/v1/patients/profile` - Get current patient profile
- `PUT /api/v1/patients/profile` - Update patient profile
- `GET /api/v1/patients` - Get all patients (doctors only)
- `GET /api/v1/patients/:id` - Get specific patient
- `GET /api/v1/patients/:id/images` - Get patient images

### Doctors
- `GET /api/v1/doctors` - Get all doctors
- `GET /api/v1/doctors/:id` - Get specific doctor

### Images
- `POST /api/v1/images/upload` - Upload retinal image
- `POST /api/v1/images/detect` - Perform AI detection
- `GET /api/v1/images` - Get user images
- `GET /api/v1/images/:id` - Get specific image
- `GET /api/v1/images/:id/file` - Serve image file

### Appointments
- `POST /api/v1/appointments` - Create appointment
- `GET /api/v1/appointments` - Get user appointments
- `GET /api/v1/appointments/:id` - Get specific appointment
- `PUT /api/v1/appointments/:id` - Update appointment
- `DELETE /api/v1/appointments/:id` - Cancel appointment

### Analytics
- `GET /api/v1/analytics/stats` - Get system statistics (doctors only)

## 🔐 Authentication

All protected endpoints require a JWT token in the Authorization header:

```
Authorization: Bearer <your-jwt-token>
```

### User Roles

- **patient**: Can manage their own profile, upload images, view results, schedule appointments
- **doctor**: Can view patients, manage appointments, review detection results
- **admin**: Full system access (future implementation)

## 💾 Data Storage

This API uses **in-memory storage** for simplicity and fast development:

- **Thread-safe operations** with mutex locks
- **Automatic UUID generation** for all entities
- **Relationship management** between users, patients, doctors, images, and appointments
- **Data persistence** only during server runtime
- **No database setup required**

### Storage Structure

- **Users**: Base user information with authentication
- **Patients**: Patient-specific medical information
- **Doctors**: Doctor credentials and specializations
- **Retinal Images**: Uploaded retinal images with metadata
- **Detection Results**: AI detection results and analysis
- **Appointments**: Patient-doctor scheduling data

## 🔧 Configuration

### Environment Variables

```env
# Server Configuration
PORT=8080
ENV=development

# JWT Configuration
JWT_SECRET=your-super-secret-jwt-key-here
JWT_EXPIRY=24h

# File Upload Configuration
MAX_FILE_SIZE=10485760
UPLOAD_DIR=./uploads
ALLOWED_EXTENSIONS=jpg,jpeg,png,tiff,bmp

# AI Model Configuration
MODEL_PATH=./models/dr_detection_model
CONFIDENCE_THRESHOLD=0.7

# CORS Configuration
ALLOWED_ORIGINS=http://localhost:3000,http://localhost:5173
```

## 🤖 AI Detection

The system includes a placeholder AI detection service that simulates diabetic retinopathy analysis. In production, this would be replaced with:

- TensorFlow/PyTorch models
- Pre-trained CNN architectures
- Real-time image processing
- Confidence scoring algorithms

### Detection Features

- **DR Stages**: No DR, Mild, Moderate, Severe, Proliferative
- **Additional Findings**: Macular edema, hemorrhages, exudates, microaneurysms
- **Confidence Scoring**: 0-1 confidence levels
- **Processing Time**: Performance metrics

## 📁 Project Structure

```
backend/
├── config/          # Configuration management
├── storage/         # In-memory data storage
├── handlers/        # HTTP request handlers
├── middleware/      # Authentication and authorization
├── routes/          # API route definitions
├── services/        # Business logic and AI detection
├── main.go          # Application entry point
├── go.mod           # Go module file
├── env.example      # Environment variables template
└── README.md        # This file
```

## 🧪 Testing

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run specific test
go test ./handlers -v
```

## 🚀 Deployment

### Docker (Recommended)

```dockerfile
FROM golang:1.21-alpine

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o main .

EXPOSE 8080
CMD ["./main"]
```

### Production Considerations

- Use environment-specific configurations
- Set up proper logging and monitoring
- Implement rate limiting
- Use HTTPS in production
- Configure proper CORS origins
- Use production-grade JWT secrets
- Consider implementing persistent storage for production use

## ⚠️ Important Notes

### Data Persistence
- **Data is lost on server restart** - this is a development/demo setup
- For production, consider implementing:
  - Database integration (PostgreSQL, MySQL)
  - File-based persistence
  - Redis for caching
  - Cloud storage for images

### Scalability
- In-memory storage is suitable for:
  - Development and testing
  - Small-scale deployments
  - Prototypes and demos
- For production scaling, implement proper database storage

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Submit a pull request

## 📄 License

This project is licensed under the MIT License.

## 🆘 Support

For support and questions:
- Create an issue in the repository
- Contact the development team
- Check the documentation

## 🔮 Future Enhancements

- Database integration for data persistence
- Real-time WebSocket notifications
- Advanced analytics and reporting
- Integration with EMR systems
- Mobile app API endpoints
- Advanced AI model integration
- Multi-language support
- Audit logging
- Advanced security features 