package middleware

import (
	"strings"

	"github.com/evrone/go-clean-template/internal/entity"
	"github.com/evrone/go-clean-template/pkg/jwt"
	"github.com/gofiber/fiber/v2"
)

func headerToken(ctx *fiber.Ctx) (string, error) {
	header := ctx.Get(fiber.HeaderAuthorization)
	if header == "" {
		return "", fiber.ErrUnauthorized
	}

	parts := strings.Fields(header)
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return "", fiber.ErrUnauthorized
	}

	return parts[1], nil
}

func UserAuth(jwtService *jwt.Service) fiber.Handler {
	return authMiddleware(jwtService, false)
}

func AdminAuth(jwtService *jwt.Service) fiber.Handler {
	return authMiddleware(jwtService, true)
}

func authMiddleware(jwtService *jwt.Service, requireAdmin bool) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		token, err := headerToken(ctx)
		if err != nil {
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "missing or invalid token"})
		}

		claims, err := jwtService.Parse(token)
		if err != nil {
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid token"})
		}

		if requireAdmin {
			if claims.Role != entity.UserRoleAdmin && claims.Role != entity.UserRoleSuperAdmin {
				return ctx.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "insufficient role"})
			}
		}

		ctx.Locals("userID", claims.UserID)
		ctx.Locals("role", claims.Role)

		return ctx.Next()
	}
}
