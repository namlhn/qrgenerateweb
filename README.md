# QR Code Generator

A modern, web-based QR code generator with customizable styling options built with Go and HTML/CSS/JavaScript.

## Features

- **Customizable QR Codes**: Choose from different dot shapes (square, circle, diamond), eye shapes, and colors
- **Logo Support**: Upload and embed logos in the center of QR codes
- **Modern UI**: Clean, responsive design with Tailwind CSS
- **RESTful API**: Simple `/api/qr` endpoint for generating QR codes

## Project Structure

```
├── main.go                     # Application entry point
├── internal/
│   ├── handlers/              # HTTP request handlers
│   │   └── qr.go
│   └── services/              # Business logic layer
│       └── qr.go
├── pkg/
│   └── qr/                    # Reusable QR generation utilities
│       └── generator.go
├── web/
│   ├── static/                # Static assets
│   │   ├── css/              # Stylesheets
│   │   ├── js/script.js       # JavaScript
│   │   └── img/gopher.png     # Images
│   └── templates/index.html    # HTML template
├── go.mod
├── go.sum
└── README.md
```

## Getting Started

### Prerequisites

- Go 1.24 or later
- Modern web browser

### Installation

1. Clone the repository:
```bash
git clone <repository-url>
cd qrgenerateweb
```

2. Install dependencies:
```bash
go mod download
```

3. Run the server:
```bash
go run main.go
```

4. Open your browser and navigate to `http://localhost:8080`

## API Usage

### Generate QR Code

**Endpoint**: `POST /api/qr`

**Form Parameters**:
- `text` (required): The text or URL to encode in the QR code
- `dotColor` (optional): Hex color for dots (default: #000000)
- `bgColor` (optional): Hex color for background (default: #ffffff)
- `eyeColor` (optional): Hex color for eyes (default: same as dotColor)
- `eyeShape` (optional): Shape for eyes (square, circle, rounded)
- `dotShape` (optional): Shape for dots (square, circle, diamond)
- `padding` (optional): Padding around QR code (0-10, default: 1)
- `logo` (optional): Image file to embed in center

**Example Request**:
```bash
curl -X POST -F "text=https://example.com" \
     -F "dotColor=#ff0000" \
     -F "bgColor=#ffffff" \
     -F "logo=@logo.png" \
     http://localhost:8080/api/qr
```

## Architecture

This project follows Go best practices with a clean architecture:

- **main.go**: Application entry point (simplified)
- **internal/**: Private application and library code
- **pkg/**: Library code that can be used by external applications
- **web/**: Web-related files (templates, static assets)

## License

MIT License - see LICENSE file for details.
