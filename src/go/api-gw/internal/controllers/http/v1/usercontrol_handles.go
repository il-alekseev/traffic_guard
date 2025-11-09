package v1

import (
	// userctrl_models "api-gateway/pkg/usercontrol/users"

	"api-gateway/internal/controllers/http/v1/utils"
	"api-gateway/pkg/models"
	"api-gateway/pkg/usercontrol/auth"
	"api-gateway/pkg/usercontrol/roles"
	"api-gateway/pkg/usercontrol/users"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// signIn godoc
// @Summary Авторизация пользователя в системе
// @Description Для авторизации требуется логин и пароль. В случае успешной авторизации пользователю выдается токен OAuth. Время жизни токена 24 часа.
// @Tags auth
// @Accept json
// @Produce json
// @Param input body models.DtoLogUser true "Данные для авторизации"
// @Success 200 {object} models.ModelsToken
// @Failure 400 {object} models.DtoErrorResponse
// @Failure 401 {object} models.DtoErrorResponse
// @Failure 403 {object} models.DtoErrorResponse
// @Failure 500 {object} models.DtoErrorResponse
// @Router /api/v1/auth/sign-in [post]
func (s *Server) signIn(c *gin.Context) {
	var input models.DtoLogUser
	if err := c.ShouldBindJSON(&input); err != nil {
		s.ErrorResponse(c, http.StatusBadRequest, "ShouldBindJSON", err)
		return
	}

	resp, err := s.userCl.Auth.PostV1AuthSignIn(
		&auth.PostV1AuthSignInParams{
			Input:   &input,
			Context: c,
		},
	)

	if err != nil {
		if conflictErr, ok := err.(ResponseErrorInterface); ok {
			c.JSON(conflictErr.Code(), conflictErr.GetPayload())
			return
		}

		s.ErrorResponse(c, http.StatusBadRequest, "PostV1AuthSignIn", err)
		return
	}
	// grafana_session
	grafana_session := resp.Payload.GrafanaSession
	grafana_session_expiry := resp.Payload.GrafanaSessionExpiry
	// Установить куки
	exp := int(time.Until(time.Now().Add(time.Hour * 24 * 30)).Seconds())
	c.SetCookie(
		"grafana_session", // имя cookie
		grafana_session,   // значение
		exp,               // время жизни в секундах
		"/",               // путь
		s.domain,          // домен (пустая строка для текущего домена)
		false,             // secure (true для HTTPS)
		true,              // httpOnly (защита от XSS)
	)

	c.SetCookie(
		"grafana_session_expiry", // имя cookie
		grafana_session_expiry,   // значение
		exp,                      // время жизни в секундах
		"/",                      // путь
		s.domain,                 // домен (пустая строка для текущего домена)
		false,                    // secure (true для HTTPS)
		true,                     // httpOnly (защита от XSS)
	)

	c.JSON(resp.Code(), resp.GetPayload())
}

// signOut godoc
// @Summary Выход пользователя из системы
// @Description Для успешного завершения запроса пользователю необходимо передать валидный токен OAuth в заголовке Authorization и refresh-токен в заголовке X-Refresh-Token.
// @Tags auth
// @Security BearerAuth
// @Param X-Refresh-Token header string true "Refresh токен"
// @Success 200 {object} models.DtoSuccessResponse
// @Failure 400 {object} models.DtoErrorResponse
// @Failure 401 {object} models.DtoErrorResponse
// @Failure 403 {object} models.DtoErrorResponse
// @Failure 500 {object} models.DtoErrorResponse
// @Router /api/v1/auth/sign-out [post]
func (s *Server) signOut(c *gin.Context) {
	refreshToken := c.GetHeader("X-Refresh-Token")
	if refreshToken == "" {
		s.ErrorResponse(c, http.StatusBadRequest, "Refresh token is required", nil)
		return
	}

	// Создаем authInfoWriter для передачи токена
	authInfo, err := utils.GetAuthInfo(c)
	if err != nil {
		s.ErrorResponse(c, http.StatusBadRequest, "utils.GetAuthInfo(c)", err)
		return
	}

	resp, err := s.userCl.Auth.PostV1AuthSignOut(
		&auth.PostV1AuthSignOutParams{
			XRefreshToken: refreshToken,
			Context:       c,
		},
		authInfo,
	)

	if err != nil {
		if conflictErr, ok := err.(ResponseErrorInterface); ok {
			c.JSON(conflictErr.Code(), conflictErr.GetPayload())
			return
		}

		s.ErrorResponse(c, http.StatusBadRequest, "PostV1AuthSignOut", err)
		return
	}

	// Удаляем куки grafana_session и grafana_session_expiry
	c.SetCookie(
		"grafana_session", // имя cookie
		"",                // пустое значение
		-1,                // время жизни -1 (удаление)
		"/",               // путь
		s.domain,          // домен
		false,             // secure (true для HTTPS)
		true,              // httpOnly (защита от XSS)
	)

	c.SetCookie(
		"grafana_session_expiry", // имя cookie
		"",                       // пустое значение
		-1,                       // время жизни -1 (удаление)
		"/",                      // путь
		s.domain,                 // домен
		false,                    // secure (true для HTTPS)
		true,                     // httpOnly (защита от XSS)
	)

	c.JSON(resp.Code(), resp.GetPayload())
}

// signOut godoc
// @Summary Обновление токена
// @Description Для успешного обновления токена пользователю необходимо передать валидный токен OAuth в заголовке Authorization и refresh-токен в заголовке X-Refresh-Token.
// @Tags auth
// @Param X-Refresh-Token header string true "Refresh токен"
// @Success 200 {object} models.ModelsToken
// @Failure 400 {object} models.DtoErrorResponse
// @Failure 401 {object} models.DtoErrorResponse
// @Failure 403 {object} models.DtoErrorResponse
// @Failure 500 {object} models.DtoErrorResponse
// @Router /api/v1/auth/refresh [post]
func (s *Server) refresh(c *gin.Context) {
	refreshToken := c.GetHeader("X-Refresh-Token")
	if refreshToken == "" {
		s.ErrorResponse(c, http.StatusBadRequest, "Refresh token is required", nil)
		return
	}

	resp, err := s.userCl.Auth.PostV1AuthRefresh(
		&auth.PostV1AuthRefreshParams{
			XRefreshToken: refreshToken,
			Context:       c,
		},
	)

	if err != nil {
		if conflictErr, ok := err.(ResponseErrorInterface); ok {
			c.JSON(conflictErr.Code(), conflictErr.GetPayload())
			return
		}

		s.ErrorResponse(c, http.StatusBadRequest, "get users", err)
		return
	}
	c.JSON(resp.Code(), resp.GetPayload())
}

// getProfile godoc
// @Summary Получение информации о текущем пользователе
// @Description Получение информации об авторизованном в системе пользователе
// @Tags users
// @Security BearerAuth
// @Success 200 {object} models.ModelsUser
// @Failure 400 {object} models.DtoErrorResponse
// @Failure 401 {object} models.DtoErrorResponse
// @Failure 403 {object} models.DtoErrorResponse
// @Failure 500 {object} models.DtoErrorResponse
// @Router /api/v1/users/profile [get]
func (s *Server) getProfile(c *gin.Context) {
	authInfo, err := utils.GetAuthInfo(c)
	if err != nil {
		s.ErrorResponse(c, http.StatusBadRequest, "utils.GetAuthInfo(c)", err)
		return
	}

	resp, err := s.userCl.Users.GetV1UsersProfile(
		&users.GetV1UsersProfileParams{
			Context: c,
		},
		authInfo,
	)

	if err != nil {
		if conflictErr, ok := err.(ResponseErrorInterface); ok {
			c.JSON(conflictErr.Code(), conflictErr.GetPayload())
			return
		}

		s.ErrorResponse(c, http.StatusBadRequest, "GetV1UsersProfile", err)
		return
	}
	c.JSON(resp.Code(), resp.GetPayload())
}

// updatePassword godoc
// @Summary Изменение пароля текущего пользователя системы
// @Description Пользователю необходимо отправить в теле запроса новый пароль.
// @Tags users
// @Accept json
// @Produce json
// @Param password body models.DtoUserUpdatePass true "Старый и новый пароль"
// @Security BearerAuth
// @Success 200 {object} models.DtoSuccessResponse
// @Failure 400 {object} models.DtoErrorResponse
// @Failure 401 {object} models.DtoErrorResponse
// @Failure 403 {object} models.DtoErrorResponse
// @Failure 500 {object} models.DtoErrorResponse
// @Router /api/v1/users/profile/pass [put]
func (s *Server) updatePassword(c *gin.Context) {
	var input models.DtoUserUpdatePass
	if err := c.ShouldBindJSON(&input); err != nil {
		s.ErrorResponse(c, http.StatusBadRequest, "ShouldBindJSON", err)
		return
	}

	authInfo, err := utils.GetAuthInfo(c)
	if err != nil {
		s.ErrorResponse(c, http.StatusBadRequest, "utils.GetAuthInfo(c)", err)
		return
	}

	resp, err := s.userCl.Users.PutV1UsersProfilePass(
		&users.PutV1UsersProfilePassParams{
			Input:   &input,
			Context: c,
		},
		authInfo,
	)

	if err != nil {
		if conflictErr, ok := err.(ResponseErrorInterface); ok {
			c.JSON(conflictErr.Code(), conflictErr.GetPayload())
			return
		}

		s.ErrorResponse(c, http.StatusBadRequest, "PutV1UsersProfilePass", err)
		return
	}
	c.JSON(resp.Code(), resp.GetPayload())
}

// updateUser godoc
// @Summary Обновление информации о пользователе системы
// @Description Запрос для обновления существующего пользователя системы.
// @Tags users
// @Accept json
// @Param userId path string true "ID пользователя"
// @Param input body models.DtoUserUpdateData true "Данные для обновления"
// @Security BearerAuth
// @Success 200 {object} models.DtoSuccessResponse
// @Failure 400 {object} models.DtoErrorResponse
// @Failure 401 {object} models.DtoErrorResponse
// @Failure 403 {object} models.DtoErrorResponse
// @Failure 404 {object} models.DtoErrorResponse
// @Failure 500 {object} models.DtoErrorResponse
// @Router /api/v1/users/{userId} [put]
func (s *Server) updateUser(c *gin.Context) {
	id := c.Param("userId")
	if id == "" {
		s.ErrorResponse(c, http.StatusBadRequest, "id is empty", nil)
		return
	}
	var input models.DtoUserUpdateData
	if err := c.ShouldBindJSON(&input); err != nil {
		s.ErrorResponse(c, http.StatusBadRequest, "ShouldBindJSON", err)
		return
	}

	authInfo, err := utils.GetAuthInfo(c)
	if err != nil {
		s.ErrorResponse(c, http.StatusBadRequest, "utils.GetAuthInfo(c)", err)
		return
	}

	resp, err := s.userCl.Users.PutV1UsersUserID(
		&users.PutV1UsersUserIDParams{
			Input:   &input,
			UserID:  id,
			Context: c,
		},
		authInfo,
	)

	if err != nil {
		if conflictErr, ok := err.(ResponseErrorInterface); ok {
			c.JSON(conflictErr.Code(), conflictErr.GetPayload())
			return
		}

		s.ErrorResponse(c, http.StatusBadRequest, "PutV1UsersUserID", err)
		return
	}
	c.JSON(resp.Code(), resp.GetPayload())
}

// deleteUser godoc
// @Summary Удаление пользователя из системы
// @Description Запрос для удаления существущего пользователя из систему.
// @Tags users
// @Param userId path string true "ID пользователя"
// @Security BearerAuth
// @Success 200 {object} models.DtoSuccessResponse
// @Failure 400 {object} models.DtoErrorResponse
// @Failure 401 {object} models.DtoErrorResponse
// @Failure 403 {object} models.DtoErrorResponse
// @Failure 404 {object} models.DtoErrorResponse
// @Failure 500 {object} models.DtoErrorResponse
// @Router /api/v1/users/{userId} [delete]
func (s *Server) deleteUser(c *gin.Context) {
	id := c.Param("userId")
	if id == "" {
		s.ErrorResponse(c, http.StatusBadRequest, "id is empty", nil)
		return
	}

	authInfo, err := utils.GetAuthInfo(c)
	if err != nil {
		s.ErrorResponse(c, http.StatusBadRequest, "utils.GetAuthInfo(c)", err)
		return
	}

	resp, err := s.userCl.Users.DeleteV1UsersUserID(
		&users.DeleteV1UsersUserIDParams{
			UserID:  id,
			Context: c,
		},
		authInfo,
	)

	if err != nil {
		if conflictErr, ok := err.(ResponseErrorInterface); ok {
			c.JSON(conflictErr.Code(), conflictErr.GetPayload())
			return
		}

		s.ErrorResponse(c, http.StatusBadRequest, "DeleteV1UsersUserID", err)
		return
	}
	c.JSON(resp.Code(), resp.GetPayload())
}

// resetPassword godoc
// @Summary Сброс пользовательского пароля
// @Description Запрос для изменения пароля пользователя.
// @Tags users
// @Accept json
// @Param userId path string true "ID пользователя"
// @Param password body models.DtoUserResetPass true "Новый пароль"
// @Security BearerAuth
// @Success 200 {object} models.DtoSuccessResponse
// @Failure 400 {object} models.DtoErrorResponse
// @Failure 401 {object} models.DtoErrorResponse
// @Failure 403 {object} models.DtoErrorResponse
// @Failure 404 {object} models.DtoErrorResponse
// @Failure 500 {object} models.DtoErrorResponse
// @Router /api/v1/users/{userId}/pass/otp [put]
func (s *Server) resetPassword(c *gin.Context) {
	id := c.Param("userId")
	if id == "" {
		s.ErrorResponse(c, http.StatusBadRequest, "Param", errors.New("userId is empty"))
		return
	}

	var input models.DtoUserResetPass
	if err := c.ShouldBindJSON(&input); err != nil {
		s.ErrorResponse(c, http.StatusBadRequest, "ShouldBindJSON", err)
		return
	}

	authInfo, err := utils.GetAuthInfo(c)
	if err != nil {
		s.ErrorResponse(c, http.StatusBadRequest, "utils.GetAuthInfo(c)", err)
		return
	}

	resp, err := s.userCl.Users.PutV1UsersUserIDPassOtp(
		&users.PutV1UsersUserIDPassOtpParams{
			Input:   &input,
			UserID:  id,
			Context: c,
		},
		authInfo,
	)

	if err != nil {
		if conflictErr, ok := err.(ResponseErrorInterface); ok {
			c.JSON(conflictErr.Code(), conflictErr.GetPayload())
			return
		}

		s.ErrorResponse(c, http.StatusBadRequest, "GetV1UsersProfile", err)
		return
	}
	c.JSON(resp.Code(), resp.GetPayload())
}

// getUsersCount godoc
// @Summary Получение числа всех пользователей системы
// @Description Получение числа всех пользователей системы.
// @Tags users
// @Security BearerAuth
// @Success 200 {integer} integer "Количество пользователей"
// @Failure 400 {object} models.DtoErrorResponse
// @Failure 401 {object} models.DtoErrorResponse
// @Failure 403 {object} models.DtoErrorResponse
// @Failure 500 {object} models.DtoErrorResponse
// @Router /api/v1/users/count [get]
func (s *Server) getUsersCount(c *gin.Context) {
	authInfo, err := utils.GetAuthInfo(c)
	if err != nil {
		s.ErrorResponse(c, http.StatusBadRequest, "utils.GetAuthInfo(c)", err)
		return
	}

	resp, err := s.userCl.Users.GetV1UsersCount(
		&users.GetV1UsersCountParams{
			Context: c,
		},
		authInfo,
	)

	if err != nil {
		if conflictErr, ok := err.(ResponseErrorInterface); ok {
			c.JSON(conflictErr.Code(), conflictErr.GetPayload())
			return
		}

		s.ErrorResponse(c, http.StatusBadRequest, "GetV1UsersCount", err)
		return
	}
	c.JSON(resp.Code(), resp.GetPayload())
}

// updateUserRoles godoc
// @Summary Обновление ролей пользователя системы
// @Description Запрос для обновления списка ролей пользователя системы.
// @Tags users,roles
// @Accept json
// @Param userId path string true "ID пользователя"
// @Param input body models.ModelsRole true "Список ролей"
// @Security BearerAuth
// @Success 200 {object} models.DtoSuccessResponse
// @Failure 400 {object} models.DtoErrorResponse
// @Failure 401 {object} models.DtoErrorResponse
// @Failure 403 {object} models.DtoErrorResponse
// @Failure 404 {object} models.DtoErrorResponse
// @Failure 500 {object} models.DtoErrorResponse
// @Router /api/v1/users/{userId}/roles [put]
func (s *Server) updateUserRoles(c *gin.Context) {
	id := c.Param("userId")
	if id == "" {
		s.ErrorResponse(c, http.StatusBadRequest, "Param", errors.New("userId is empty"))
		return
	}

	var input models.ModelsRole
	if err := c.ShouldBindJSON(&input); err != nil {
		s.ErrorResponse(c, http.StatusBadRequest, "ShouldBindJSON", err)
		return
	}

	authInfo, err := utils.GetAuthInfo(c)
	if err != nil {
		s.ErrorResponse(c, http.StatusBadRequest, "utils.GetAuthInfo(c)", err)
		return
	}

	resp, err := s.userCl.Users.PutV1UsersUserIDRoles(
		&users.PutV1UsersUserIDRolesParams{
			Input:   &input,
			UserID:  id,
			Context: c,
		},
		authInfo,
	)

	if err != nil {
		if conflictErr, ok := err.(ResponseErrorInterface); ok {
			c.JSON(conflictErr.Code(), conflictErr.GetPayload())
			return
		}

		s.ErrorResponse(c, http.StatusBadRequest, "GetV1UsersProfile", err)
		return
	}
	c.JSON(resp.Code(), resp.GetPayload())
}

// deleteUserRoles godoc
// @Summary Отзыв ролей пользователя в системе
// @Description Запрос для отзыва ролей пользователя из системы.
// @Tags users,roles
// @Accept json
// @Param userId path string true "ID пользователя"
// @Param input body models.ModelsRole true "Список ролей для отзыва"
// @Security BearerAuth
// @Success 200 {object} models.DtoSuccessResponse
// @Failure 400 {object} models.DtoErrorResponse
// @Failure 401 {object} models.DtoErrorResponse
// @Failure 403 {object} models.DtoErrorResponse
// @Failure 404 {object} models.DtoErrorResponse
// @Failure 500 {object} models.DtoErrorResponse
// @Router /api/v1/users/{userId}/roles [delete]
func (s *Server) deleteUserRoles(c *gin.Context) {
	id := c.Param("userId")
	if id == "" {
		s.ErrorResponse(c, http.StatusBadRequest, "Param", errors.New("userId is empty"))
		return
	}

	var input models.ModelsRole
	if err := c.ShouldBindJSON(&input); err != nil {
		s.ErrorResponse(c, http.StatusBadRequest, "ShouldBindJSON", err)
		return
	}

	authInfo, err := utils.GetAuthInfo(c)
	if err != nil {
		s.ErrorResponse(c, http.StatusBadRequest, "utils.GetAuthInfo(c)", err)
		return
	}

	resp, err := s.userCl.Users.DeleteV1UsersUserIDRoles(
		&users.DeleteV1UsersUserIDRolesParams{
			Input:   &input,
			UserID:  id,
			Context: c,
		},
		authInfo,
	)

	if err != nil {
		if conflictErr, ok := err.(ResponseErrorInterface); ok {
			c.JSON(conflictErr.Code(), conflictErr.GetPayload())
			return
		}

		s.ErrorResponse(c, http.StatusBadRequest, "DeleteV1UsersUserIDRoles", err)
		return
	}
	c.JSON(resp.Code(), resp.GetPayload())
}

// getUsers godoc
// @Summary Получение списка пользователей с пагинацией и фильтрацией
// @Description Возвращает список пользователей с возможностью пагинации и фильтрации по ролям
// @Tags users
// @Param page query int false "Номер страницы с 1" default(1)
// @Param limit query int false "Количество элементов на странице" default(10)
// @Param role query string false "Фильтр по роли" Enums(SA, CA, CO)
// @Param contextID query string false "Фильтр по contextID"
// @Param search query string false "Поиск по имени/email/login"
// @Security BearerAuth
// @Success 200 {object} models.DtoUsersListResponse
// @Failure 400 {object} models.DtoErrorResponse
// @Failure 401 {object} models.DtoErrorResponse
// @Failure 403 {object} models.DtoErrorResponse
// @Failure 500 {object} models.DtoErrorResponse
// @Router /api/v1/users [get]
func (s *Server) getUsers(c *gin.Context) {
	page, err := strconv.ParseInt(c.Query("page"), 10, 64)
	if err != nil {
		s.ErrorResponse(c, http.StatusBadRequest, "'page' value must be int64", err)
		return
	}
	limit, err := strconv.ParseInt(c.Query("limit"), 10, 64)
	if err != nil {
		s.ErrorResponse(c, http.StatusBadRequest, "'limit' value must be int64", err)
		return
	}
	role := c.Query("role")
	search := c.Query("search")
	contextID := c.Query("contextID")

	// Создаем authInfoWriter для передачи токена
	authInfo, err := utils.GetAuthInfo(c)
	if err != nil {
		s.ErrorResponse(c, http.StatusBadRequest, "utils.GetAuthInfo(c)", err)
		return
	}

	resp, err := s.userCl.Users.GetV1Users(
		&users.GetV1UsersParams{
			Page:      &page,
			Limit:     &limit,
			Search:    &search,
			Role:      &role,
			ContextID: &contextID,
			Context:   c,
		},
		authInfo,
	)

	if err != nil {
		if conflictErr, ok := err.(ResponseErrorInterface); ok {
			c.JSON(conflictErr.Code(), conflictErr.GetPayload())
			return
		}

		s.ErrorResponse(c, http.StatusBadRequest, "get users", err)
		return
	}
	c.JSON(resp.Code(), resp.GetPayload())
}

// createUser godoc
// @Summary Создание нового пользователя
// @Description Создает нового пользователя с указанными данными
// @Tags users
// @Accept json
// @Produce json
// @Param input body models.DtoUserCreateRequest true "Данные пользователя"
// @Security BearerAuth
// @Success 201 {object} models.ModelsUser
// @Failure 400 {object} models.DtoErrorResponse
// @Failure 401 {object} models.DtoErrorResponse
// @Failure 403 {object} models.DtoErrorResponse
// @Failure 409 {object} models.DtoErrorResponse "Пользователь с таким email/логином уже существует или роль не найдена"
// @Failure 500 {object} models.DtoErrorResponse
// @Router /api/v1/users [post]
func (s *Server) createUser(c *gin.Context) {
	var input models.DtoUserCreateRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		s.ErrorResponse(c, http.StatusBadRequest, "ShouldBindJSON", err)
		return
	}

	authInfo, err := utils.GetAuthInfo(c)
	if err != nil {
		s.ErrorResponse(c, http.StatusBadRequest, "utils.GetAuthInfo(c)", err)
		return
	}

	resp, err := s.userCl.Users.PostV1Users(
		&users.PostV1UsersParams{
			Input:   &input,
			Context: c,
		},
		authInfo,
	)

	if err != nil {
		if conflictErr, ok := err.(ResponseErrorInterface); ok {
			c.JSON(conflictErr.Code(), conflictErr.GetPayload())
			return
		}

		s.ErrorResponse(c, http.StatusInternalServerError, "PostV1Users", err)
		return
	}

	c.JSON(resp.Code(), resp.GetPayload())
}

// getUser godoc
// @Summary Получение информации о пользователе
// @Description Возвращает полную информацию о пользователе по его ID
// @Tags users
// @Param userId path string true "ID пользователя"
// @Security BearerAuth
// @Success 200 {object} models.ModelsUser
// @Failure 400 {object} models.DtoErrorResponse
// @Failure 401 {object} models.DtoErrorResponse
// @Failure 403 {object} models.DtoErrorResponse
// @Failure 404 {object} models.DtoErrorResponse "Пользователь не найден"
// @Failure 500 {object} models.DtoErrorResponse
// @Router /api/v1/users/{userId} [get]
func (s *Server) getUser(c *gin.Context) {
	id := c.Param("userId")
	if id == "" {
		s.ErrorResponse(c, http.StatusBadRequest, "Param", errors.New("userId is empty"))
		return
	}

	authInfo, err := utils.GetAuthInfo(c)
	if err != nil {
		s.ErrorResponse(c, http.StatusBadRequest, "utils.GetAuthInfo(c)", err)
		return
	}

	resp, err := s.userCl.Users.GetV1UsersUserID(
		&users.GetV1UsersUserIDParams{
			UserID:  id,
			Context: c,
		},
		authInfo,
	)

	if err != nil {
		if conflictErr, ok := err.(ResponseErrorInterface); ok {
			c.JSON(conflictErr.Code(), conflictErr.GetPayload())
			return
		}

		s.ErrorResponse(c, http.StatusBadRequest, "GetV1UsersUserID", err)
		return
	}
	c.JSON(resp.Code(), resp.GetPayload())
}

// getUsers godoc
// @Summary Получение списка ролей с пагинацией и фильтрацией
// @Description Возвращает список ролей с возможностью пагинации и фильтрации по названию
// @Tags roles
// @Param page query int false "Номер страницы" default(1)
// @Param limit query int false "Количество элементов на странице" default(10)
// @Param search query string false "Поиск по названию роли"
// @Security BearerAuth
// @Success 200 {object} models.DtoRolesListResponse
// @Failure 400 {object} models.DtoErrorResponse
// @Failure 401 {object} models.DtoErrorResponse
// @Failure 403 {object} models.DtoErrorResponse
// @Failure 500 {object} models.DtoErrorResponse
// @Router /api/v1/roles [get]
func (s *Server) getRoles(c *gin.Context) {
	page, err := strconv.ParseInt(c.Query("page"), 10, 64)
	if err != nil {
		s.ErrorResponse(c, http.StatusBadRequest, "'page' value must be int64", err)
		return
	}
	limit, err := strconv.ParseInt(c.Query("limit"), 10, 64)
	if err != nil {
		s.ErrorResponse(c, http.StatusBadRequest, "'limit' value must be int64", err)
		return
	}
	search := c.Query("search")

	// Создаем authInfoWriter для передачи токена
	authInfo, err := utils.GetAuthInfo(c)
	if err != nil {
		s.ErrorResponse(c, http.StatusBadRequest, "utils.GetAuthInfo(c)", err)
		return
	}

	resp, err := s.userCl.Roles.GetV1Roles(
		&roles.GetV1RolesParams{
			Page:    &page,
			Limit:   &limit,
			Search:  &search,
			Context: c,
		},
		authInfo,
	)

	if err != nil {
		if conflictErr, ok := err.(ResponseErrorInterface); ok {
			c.JSON(conflictErr.Code(), conflictErr.GetPayload())
			return
		}

		s.ErrorResponse(c, http.StatusBadRequest, "GetV1Roles", err)
		return
	}
	c.JSON(resp.Code(), resp.GetPayload())
}

// getRolesCount godoc
// @Summary Получение числа всех ролей системы
// @Description Получение числа всех ролей системы.
// @Tags roles
// @Security BearerAuth
// @Success 200 {object} models.DtoCountResponse "Количество ролей"
// @Failure 400 {object} models.DtoErrorResponse
// @Failure 401 {object} models.DtoErrorResponse
// @Failure 403 {object} models.DtoErrorResponse
// @Failure 500 {object} models.DtoErrorResponse
// @Router /api/v1/roles/count [get]
func (s *Server) getRolesCount(c *gin.Context) {
	authInfo, err := utils.GetAuthInfo(c)
	if err != nil {
		s.ErrorResponse(c, http.StatusBadRequest, "utils.GetAuthInfo(c)", err)
		return
	}

	resp, err := s.userCl.Roles.GetV1RolesCount(
		&roles.GetV1RolesCountParams{
			Context: c,
		},
		authInfo,
	)

	if err != nil {
		if conflictErr, ok := err.(ResponseErrorInterface); ok {
			c.JSON(conflictErr.Code(), conflictErr.GetPayload())
			return
		}

		s.ErrorResponse(c, http.StatusBadRequest, "GetV1RolesCount", err)
		return
	}
	c.JSON(resp.Code(), resp.GetPayload())
}

// getProfile godoc
// @Summary Получение информации о числе пользователй по ролям
// @Description Получение информации сколько в системе SA, CA и CO
// @Tags users
// @Security BearerAuth
// @Success 200 {object} models.DtoCountUserByRoleResponse
// @Failure 400 {object} models.DtoErrorResponse
// @Failure 401 {object} models.DtoErrorResponse
// @Failure 403 {object} models.DtoErrorResponse
// @Failure 500 {object} models.DtoErrorResponse
// @Router /api/v1/users/count-by-role [get]
func (s *Server) getCountUsersByRole(c *gin.Context) {
	authInfo, err := utils.GetAuthInfo(c)
	if err != nil {
		s.ErrorResponse(c, http.StatusBadRequest, "utils.GetAuthInfo(c)", err)
		return
	}

	resp, err := s.userCl.Users.GetV1UsersCountByRole(
		&users.GetV1UsersCountByRoleParams{
			Context: c,
		},
		authInfo,
	)

	if err != nil {
		if conflictErr, ok := err.(ResponseErrorInterface); ok {
			c.JSON(conflictErr.Code(), conflictErr.GetPayload())
			return
		}

		s.ErrorResponse(c, http.StatusBadRequest, "GetV1RolesCount", err)
		return
	}
	c.JSON(resp.Code(), resp.GetPayload())
}
