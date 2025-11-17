
package middleware

import (
    "net/http"
    "strings"

    "github.com/gin-gonic/gin"
    "github.com/golang-jwt/jwt/v5"
)

// AuthMiddleware vérifie l'authentification JWT
func AuthMiddleware(jwtSecret string) gin.HandlerFunc {
    return func(c *gin.Context) {
        // Récupération du token depuis l'en-tête Authorization
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            c.JSON(http.StatusUnauthorized, gin.H{
                "error": "Token d'authentification requis",
            })
            c.Abort()
            return
        }

        // Vérification du format Bearer
        tokenString := ""
        if strings.HasPrefix(authHeader, "Bearer ") {
            tokenString = strings.TrimPrefix(authHeader, "Bearer ")
        } else {
            c.JSON(http.StatusUnauthorized, gin.H{
                "error": "Format d'authentification invalide (utilisez Bearer <token>)",
            })
            c.Abort()
            return
        }

        // Parsing et validation du token
        token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
            // Vérification de l'algorithme de signature
            if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
                return nil, jwt.ErrSignatureInvalid
            }
            return []byte(jwtSecret), nil
        })

        if err != nil {
            c.JSON(http.StatusUnauthorized, gin.H{
                "error": "Token invalide",
                "details": err.Error(),
            })
            c.Abort()
            return
        }

        // Extraction des claims
        if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
            // Ajout des informations utilisateur au contexte
            c.Set("user_id", claims["user_id"])
            c.Set("username", claims["sub"])
            c.Set("role", claims["role"])
        } else {
            c.JSON(http.StatusUnauthorized, gin.H{
                "error": "Token invalide",
            })
            c.Abort()
            return
        }

        c.Next()
    }
}

// RequireRole vérifie que l'utilisateur a le rôle requis
func RequireRole(allowedRoles ...string) gin.HandlerFunc {
    return func(c *gin.Context) {
        userRole, exists := c.Get("role")
        if !exists {
            c.JSON(http.StatusForbidden, gin.H{
                "error": "Rôle utilisateur non trouvé",
            })
            c.Abort()
            return
        }

        roleStr, ok := userRole.(string)
        if !ok {
            c.JSON(http.StatusForbidden, gin.H{
                "error": "Rôle utilisateur invalide",
            })
            c.Abort()
            return
        }

        // L'admin a tous les droits
        if roleStr == "admin" {
            c.Next()
            return
        }

        // Vérification des rôles autorisés
        for _, allowedRole := range allowedRoles {
            if roleStr == allowedRole {
                c.Next()
                return
            }
        }

        c.JSON(http.StatusForbidden, gin.H{
            "error": "Permissions insuffisantes",
            "required_roles": allowedRoles,
            "user_role": roleStr,
        })
        c.Abort()
    }
}
