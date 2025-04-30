package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/Nerzal/gocloak/v12"
	"github.com/golang-jwt/jwt/v4"
	"github.com/rs/cors"
)

// Middleware для проверки токена Keycloak
func keycloakMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		client := gocloak.NewClient("http://keycloak:8080")
		ctx := r.Context()

		// Извлечение токена из заголовка Authorization
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Missing Authorization header", http.StatusUnauthorized)
			return
		}

		token := authHeader[len("Bearer "):]
		_, claims, err := client.DecodeAccessToken(ctx, token, "reports-realm")
		if err != nil || claims == nil {
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

		// Проверка ролей (например, доступ только для "prothetic_user")
		if !hasRole(*claims, "prothetic_user") {
			http.Error(w, "Forbidden: Insufficient permissions", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// Проверка наличия роли в токене
func hasRole(claims jwt.MapClaims, role string) bool {
	// Получаем доступ к полю "realm_access.roles" из claims
	if realmAccess, ok := claims["realm_access"].(map[string]interface{}); ok {
		if roles, ok := realmAccess["roles"].([]interface{}); ok {
			for _, r := range roles {
				if r == role {
					return true
				}
			}
		}
	}
	return false
}

// Генерация случайных чисел
func generateRandomNumbers(count int) []int {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	numbers := make([]int, count)
	for i := 0; i < count; i++ {
		numbers[i] = r.Intn(100)
	}
	return numbers
}

// Обработчик для /reports
func reportsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	numbers := generateRandomNumbers(10)
	json.NewEncoder(w).Encode(numbers)
}

func main() {
	// Настройка маршрутов
	mux := http.NewServeMux()
	mux.Handle("/reports", keycloakMiddleware(http.HandlerFunc(reportsHandler)))

	// Настройка CORS
	corsHandler := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:3000"}, // Разрешаем запросы только с localhost:3000
		AllowedMethods: []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders: []string{"Authorization", "Content-Type"},
	})

	// Применяем middleware для CORS
	handler := corsHandler.Handler(mux)

	// Запуск сервера
	log.Println("Starting server on :8000")
	log.Fatal(http.ListenAndServe(":8000", handler))
}
