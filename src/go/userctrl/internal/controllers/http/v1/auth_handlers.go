package v1

import (
	"errors"
	"userctrl/internal/controllers/http/v1/dto"
	"net/http"

	"github.com/gin-gonic/gin"
)

// signIn godoc
// @Summary Авторизация пользователя в системе
// @Description Для авторизации требуется логин и пароль. В случае успешной авторизации пользователю выдается токен OAuth. Время жизни токена 24 часа.
// @Tags auth
// @Accept json
// @Produce json
// @Param input body dto.LogUser true "Данные для авторизации"
// @Success 200 {object} models.Token
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 403 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /v1/auth/sign-in [post]
func (s *Server) signIn(c *gin.Context) {
	var input dto.LogUser
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	// при успешной аутентификации возвращает токен доступа
	token, err := s.u.SignIn(c.Request.Context(), input.Login, input.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, token)
}

// signOut godoc
// @Summary Выход пользователя из системы
// @Description Для успешного завершения запроса пользователю необходимо передать валидный токен OAuth в заголовке Authorization и refresh-токен в заголовке X-Refresh-Token.
// @Tags auth
// @Security BearerAuth
// @Param X-Refresh-Token header string true "Refresh токен"
// @Success 200 {object} dto.SuccessResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 403 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /v1/auth/sign-out [post]
func (s *Server) signOut(c *gin.Context) {
	// требует наличие заголовков Authorization (access token) и X-Refresh-Token.
	token := c.GetHeader("Authorization")
	if token == "" {
		s.ErrorResponse(c, http.StatusUnauthorized, "Authorization", errors.New("authorization token is required"))
		return
	}

	refreshToken := c.GetHeader("X-Refresh-Token")
	if refreshToken == "" {
		s.ErrorResponse(c, http.StatusBadRequest, "X-Refresh-Token", errors.New("refresh token is required"))
		return
	}

	// передаёт refresh token в бизнес-логику для удаления или аннулирования
	if err := s.u.SignOut(c.Request.Context(), refreshToken); err != nil {
		// при отсутствии токенов возвращает ошибку 401 или 400
		s.ErrorResponse(c, http.StatusInternalServerError, "SignOut", err)
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse{Message: "Successfully signed out"})
}

// refresh godoc
// @Summary Обновление токена
// @Description Для успешного обновления токена пользователю необходимо передать валидный токен OAuth в заголовке Authorization и refresh-токен в заголовке X-Refresh-Token.
// @Tags auth
// @Param X-Refresh-Token header string true "Refresh токен"
// @Success 200 {object} models.Token
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 403 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /v1/auth/refresh [post]
func (s *Server) refresh(c *gin.Context) {
	// требует наличие заголовков X-Refresh-Token.
	// получает refreshToken токен из заголовка
	refreshToken := c.GetHeader("X-Refresh-Token")
	if refreshToken == "" {
		s.ErrorResponse(c, http.StatusBadRequest, "X-Refresh-Token", errors.New("refresh token is required"))
		return
	}

	// передаёт refresh token в бизнес-логику для обновления токена
	newToken, err := s.u.RefreshToken(c.Request.Context(), refreshToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, newToken)
}
