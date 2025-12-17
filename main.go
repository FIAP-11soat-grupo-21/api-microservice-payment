// @title API Microservice Payment
// @version 1.0
// @description This is an API for a payments on tech challenge.
// @license.name MIT
// @license.url https://opensource.org/licenses/MIT
// @host localhost:8080
// @BasePath /v1
// @schemes http
//
//go:debug x509negativeserial=1
package main

import "payment_microservice/internal/common/infra/api"

func main() {
	api.Init()
}
