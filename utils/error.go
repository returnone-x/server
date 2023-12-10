package untils

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
)

// this is for mutiple values check the value is exists
func RequestDataRequired(data map[string]string, data_name []string) fiber.Map {
	for i := 0; i < len(data_name); i++ {
		if data[data_name[i]] == "" {
			return fiber.Map{
				"status": "error",
				"message": fmt.Sprintf("%v is required", data_name[i]),
			}
		}
	}
	return nil
}
// this is for single values check the value is exists
func RequestValueRequired(value string) fiber.Map {
	return fiber.Map{
		"status": "error",
		"message": fmt.Sprintf("%v is required", value),
	}
}

func InvalidRequest() fiber.Map{
	return fiber.Map{
		"status": "error",
		"message": "Invalid request",
	}
}
func RequestValueValid(value string) fiber.Map {
	return fiber.Map{
		"status": "error",
		"message": fmt.Sprintf("This %v is not valid", value),
	}
}

func RequestValueInUse(value string) fiber.Map {
	return fiber.Map{
		"status": "error",
		"message": fmt.Sprintf("This %v is already in use", value),
		"inuse": true,
	}
}

func ErrorMessage(message string, err error) fiber.Map {
	return fiber.Map{
		"status": "error",
		"message": message,
		"error":   err,
	}
}
