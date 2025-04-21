package docs

import "github.com/swaggo/swag"

// @title Package Tracking API
// @version 1.0
// @description A microservice for tracking packages with MongoDB backend
// @host localhost:8080
// @BasePath /api/v1
// @schemes http
// @contact.name API Support
// @contact.email support@example.com
// @license.name MIT
// @license.url https://opensource.org/licenses/MIT
func SwaggerInfo() {
	swag.Register(swag.Name, &swag.Spec{})
}
