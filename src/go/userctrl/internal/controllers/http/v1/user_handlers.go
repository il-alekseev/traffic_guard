package v1

import (
	"errors"
	"math"
	"net/http"
	"strconv"
	"userctrl/internal/controllers/http/v1/dto"
	"userctrl/internal/controllers/http/v1/utils"
	"userctrl/internal/models"
	"userctrl/pkg/keycloakclient"
	"userctrl/pkg/slogger/wsl"

	"github.com/gin-gonic/gin"
)

// getProfile godoc
// @Summary Получение информации о текущем пользователе
// @Description Получение информации об авторизованном в системе пользователе
// @Tags users
// @Security BearerAuth
// @Success 200 {object} models.User
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 403 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /v1/users/profile [get]
func (s *Server) getProfile(c *gin.Context) {
	userMeta, err := utils.GetUserMeta(c)
	if err != nil {
		s.ErrorResponse(c, http.StatusBadRequest, "utils.GetUserMeta", err)
		return
	}

	if userMeta.UUID == "" {
		s.ErrorResponse(c, http.StatusBadRequest, "token param", errors.New("userID in token is not exists"))
		return
	}

	// Получение информации об авторизованном в системе пользователе
	user, err := s.u.GetByID(c.Request.Context(), userMeta, userMeta.UUID)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if errors.Is(err, models.ErrPermissionDenied) {
			statusCode = http.StatusForbidden
		}
		s.ErrorResponse(c, statusCode, "GetByID", err)
		return
	}

	c.JSON(http.StatusOK, user)
}

// updatePassword godoc
// @Summary Изменение пароля текущего пользователя системы
// @Description Пользователю необходимо отправить в теле запроса новый пароль.
// @Tags users
// @Accept json
// @Produce json
// @Param input body dto.UserUpdatePass true "Новый пароль"
// @Security BearerAuth
// @Success 200 {object} dto.SuccessResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 403 {object} dto.ErrorResponse
// @Failure 422 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /v1/users/profile/pass [put]
func (s *Server) updatePassword(c *gin.Context) {
	userMeta, err := utils.GetUserMeta(c)
	if err != nil {
		s.ErrorResponse(c, http.StatusBadRequest, "utils.GetUserMeta", err)
		return
	}

	if userMeta.UUID == "" {
		s.ErrorResponse(c, http.StatusBadRequest, "token param", errors.New("userID in token is not exists"))
		return
	}

	// валидация входящего нового пароля
	var input dto.UserUpdatePass
	if err := c.ShouldBindJSON(&input); err != nil {
		s.ErrorResponse(c, http.StatusBadRequest, "ShouldBindJSON", err)
		return
	}

	// Изменение пароля текущего пользователя системы
	err = s.u.UpdatePassword(c.Request.Context(), userMeta, userMeta.UUID, input.OldPass, input.NewPass)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if errors.Is(err, keycloakclient.ErrInvalidCredentials) {
			statusCode = http.StatusUnprocessableEntity
		} else if errors.Is(err, models.ErrPermissionDenied) {
			statusCode = http.StatusForbidden
		}

		s.ErrorResponse(c, statusCode, "UpdatePassword", err)
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse{Message: "Pass updated successfully"})
}

// updateUser godoc
// @Summary Обновление информации о пользователе системы
// @Description Запрос для обновления существующего пользователя системы.
// @Tags users
// @Accept json
// @Param userId path string true "ID пользователя"
// @Param input body dto.UserUpdateData true "Данные для обновления"
// @Security BearerAuth
// @Success 200 {object} dto.SuccessResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 403 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /v1/users/{userId} [put]
func (s *Server) updateUser(c *gin.Context) {
	userMeta, err := utils.GetUserMeta(c)
	if err != nil {
		s.ErrorResponse(c, http.StatusBadRequest, "utils.GetUserMeta", err)
		return
	}

	userID := c.Param("userId")
	if userID == "" {
		s.ErrorResponse(c, http.StatusBadRequest, "param", errors.New("userID is empty"))
		return
	}

	// валидация входящих данных
	var input dto.UserUpdateData
	if err := c.ShouldBindJSON(&input); err != nil {
		s.ErrorResponse(c, http.StatusBadRequest, "ShouldBindJSON", err)
		return
	}

	// Обновление информации о пользователе системы
	if err := s.u.Update(c.Request.Context(), userMeta, userID, input); err != nil {
		statusCode := http.StatusInternalServerError
		if errors.Is(err, models.ErrPermissionDenied) {
			statusCode = http.StatusForbidden
		} else if errors.Is(err, keycloakclient.ErrUserNotFound) {
			statusCode = http.StatusNotFound
		}
		s.ErrorResponse(c, statusCode, "Update", err)
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse{Message: "User updated successfully"})
}

// deleteUser godoc
// @Summary Удаление пользователя из системы
// @Description Запрос для удаления существующего пользователя из систему.
// @Tags users
// @Param userId path string true "ID пользователя"
// @Security BearerAuth
// @Success 200 {object} dto.SuccessResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 403 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /v1/users/{userId} [delete]
func (s *Server) deleteUser(c *gin.Context) {
	userMeta, err := utils.GetUserMeta(c)
	if err != nil {
		s.ErrorResponse(c, http.StatusBadRequest, "utils.GetUserMeta", err)
		return
	}

	userID := c.Param("userId")
	if userID == "" {
		s.ErrorResponse(c, http.StatusBadRequest, "param", errors.New("userID is empty"))
		return
	}

	// Удаление пользователя из системы
	if err := s.u.Delete(c.Request.Context(), userMeta, userID); err != nil {
		statusCode := http.StatusInternalServerError
		if errors.Is(err, models.ErrPermissionDenied) {
			statusCode = http.StatusForbidden
		} else if errors.Is(err, keycloakclient.ErrUserNotFound) {
			statusCode = http.StatusNotFound
		}
		s.ErrorResponse(c, statusCode, "Delete", err)
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse{Message: "User deleted successfully"})
}

// resetPassword godoc
// @Summary Сброс пользовательского пароля
// @Description Запрос для изменения пароля пользователя.
// @Tags users
// @Accept json
// @Param userId path string true "ID пользователя"
// @Param input body dto.UserResetPass true "Новый пароль"
// @Security BearerAuth
// @Success 200 {object} dto.SuccessResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 403 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /v1/users/{userId}/pass/otp [put]
func (s *Server) resetPassword(c *gin.Context) {
	userMeta, err := utils.GetUserMeta(c)
	if err != nil {
		s.ErrorResponse(c, http.StatusBadRequest, "utils.GetUserMeta", err)
		return
	}

	userID := c.Param("userId")
	if userID == "" {
		s.ErrorResponse(c, http.StatusBadRequest, "param", errors.New("userID is empty"))
		return
	}

	// валидация входящих данных
	var input dto.UserResetPass
	if err := c.ShouldBindJSON(&input); err != nil {
		s.ErrorResponse(c, http.StatusBadRequest, "ShouldBindJSON", err)
		return
	}

	// Сброс пользовательского пароля
	if err := s.u.ResetPassword(c.Request.Context(), userMeta, userID, input.NewPass); err != nil {
		statusCode := http.StatusInternalServerError
		if errors.Is(err, models.ErrPermissionDenied) {
			statusCode = http.StatusForbidden
		} else if errors.Is(err, keycloakclient.ErrUserNotFound) {
			statusCode = http.StatusNotFound
		}
		s.ErrorResponse(c, statusCode, "ResetPassword", err)
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse{Message: "Password reset successfully"})
}

// getUsersCount godoc
// @Summary Получение числа всех пользователей системы
// @Description Получение числа всех пользователей системы.
// @Tags users
// @Security BearerAuth
// @Success 200 {object} dto.CountResponse "Количество пользователей"
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 403 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /v1/users/count [get]
func (s *Server) getUsersCount(c *gin.Context) {
	userMeta, err := utils.GetUserMeta(c)
	if err != nil {
		s.ErrorResponse(c, http.StatusBadRequest, "utils.GetUserMeta", err)
		return
	}

	// Получение числа всех пользователей системы
	count, err := s.u.Count(c.Request.Context(), userMeta)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if errors.Is(err, models.ErrPermissionDenied) {
			statusCode = http.StatusForbidden
		}
		s.ErrorResponse(c, statusCode, "Count", err)
		return
	}
	res := dto.CountResponse{
		Count: count,
	}
	c.JSON(http.StatusOK, res)
}

// getUsersCountForContext godoc
// @Summary Получение числа CA и CO в контексте
// @Description Получение числа CA и CO в контексте.
// @Tags users,context
// @Param contextID path string true "ID контекста"
// @Security BearerAuth
// @Success 200 {object} dto.CountUserForContextResponse "Количество пользователей"
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 403 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /v1/users/count/context/{contextID} [get]
func (s *Server) getUsersCountForContext(c *gin.Context) {
	userMeta, err := utils.GetUserMeta(c)
	if err != nil {
		s.ErrorResponse(c, http.StatusBadRequest, "utils.GetUserMeta", err)
		return
	}

	contextID := c.Param("contextID")
	if contextID == "" {
		s.ErrorResponse(c, http.StatusBadRequest, "param", errors.New("contextID is empty"))
		return
	}

	// Получение числа CA и CO в контексте
	meta, err := s.u.GetUsersCountForContext(c.Request.Context(), userMeta, contextID)
	if err != nil {
		s.ErrorResponse(c, http.StatusInternalServerError, "GetUsersCountForContext", err)
		return
	}
	res := dto.CountUserForContextResponse{
		CountCA: meta.CountCA,
		CountCO: meta.CountCO,
	}
	c.JSON(http.StatusOK, res)
}

// updateUserRoles godoc
// @Summary Обновление ролей пользователя системы
// @Description Запрос для обновления списка ролей пользователя системы.
// @Tags users,roles
// @Accept json
// @Param userId path string true "ID пользователя"
// @Param input body models.Role true "Список ролей"
// @Security BearerAuth
// @Success 200 {object} dto.SuccessResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 403 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /v1/users/{userId}/roles [put]
func (s *Server) updateUserRoles(c *gin.Context) {
	userMeta, err := utils.GetUserMeta(c)
	if err != nil {
		s.ErrorResponse(c, http.StatusBadRequest, "utils.GetUserMeta", err)
		return
	}

	userID := c.Param("userId")
	if userID == "" {
		s.ErrorResponse(c, http.StatusBadRequest, "param", errors.New("userID is empty"))
		return
	}

	// валидация входящих данных
	var input models.Role
	if err := c.ShouldBindJSON(&input); err != nil {
		s.ErrorResponse(c, http.StatusBadRequest, "ShouldBindJSON", err)
		return
	}

	// Обновление ролей пользователя системы
	if err := s.u.UpdateRole(c.Request.Context(), userMeta, userID, input.Role); err != nil {
		statusCode := http.StatusInternalServerError
		if errors.Is(err, models.ErrPermissionDenied) {
			statusCode = http.StatusForbidden
		} else if errors.Is(err, keycloakclient.ErrUserNotFound) {
			statusCode = http.StatusNotFound
		} else if errors.Is(err, keycloakclient.ErrInvalidRole) {
			statusCode = http.StatusBadRequest
		}
		s.ErrorResponse(c, statusCode, "UpdateRole", err)
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse{Message: "Roles updated successfully"})
}

// deleteUserRoles godoc
// @Summary Отзыв ролей пользователя в системе
// @Description Запрос для отзыва ролей пользователя из системы.
// @Tags users,roles
// @Accept json
// @Param userId path string true "ID пользователя"
// @Param input body models.Role true "Список ролей для отзыва"
// @Security BearerAuth
// @Success 200 {object} dto.SuccessResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 403 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /v1/users/{userId}/roles [delete]
func (s *Server) deleteUserRoles(c *gin.Context) {
	userMeta, err := utils.GetUserMeta(c)
	if err != nil {
		s.ErrorResponse(c, http.StatusBadRequest, "utils.GetUserMeta", err)
		return
	}

	userID := c.Param("userId")
	if userID == "" {
		s.ErrorResponse(c, http.StatusBadRequest, "param", errors.New("userID is empty"))
		return
	}

	// валидация входящих данных
	var input models.Role
	if err := c.ShouldBindJSON(&input); err != nil {
		s.ErrorResponse(c, http.StatusBadRequest, "ShouldBindJSON", err)
		return
	}

	// Отзыв ролей пользователя в системе
	if err := s.u.RemoveRole(c.Request.Context(), userMeta, userID, input.Role); err != nil {
		s.ErrorResponse(c, http.StatusInternalServerError, "RemoveRole", err)
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse{Message: "Roles removed successfully"})
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
// @Success 200 {object} dto.UsersListResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 403 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /v1/users [get]
func (s *Server) getUsers(c *gin.Context) {
	log := s.logger.With(wsl.Label("method", "getUsers"))
	userMeta, err := utils.GetUserMeta(c)
	if err != nil {
		s.ErrorResponse(c, http.StatusBadRequest, "utils.GetUserMeta", err)
		return
	}

	// Получаем параметры запроса
	// TODO: от api gw приходит int64 - нужно проверять это
	page, err := strconv.Atoi(c.DefaultQuery("page", "0"))
	if err != nil {
		s.ErrorResponse(c, http.StatusBadRequest, "page", err)
		return
	}
	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil {
		s.ErrorResponse(c, http.StatusBadRequest, "limit", err)
		return
	}
	role := c.Query("role")
	search := c.Query("search")
	contextID := c.Query("contextID")

	// Валидация параметров
	if page < 1 { // нумерация с 1
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	// Получение списка пользователей с пагинацией и фильтрацией
	users, total, err := s.u.GetUsers(c.Request.Context(), userMeta, page, limit, role, contextID, search)
	if err != nil {
		log.ErrorContext(c.Request.Context(), "Failed to get users", wsl.Err(err))
		statusCode := http.StatusInternalServerError
		if errors.Is(err, models.ErrPermissionDenied) {
			statusCode = http.StatusForbidden
		}

		s.ErrorResponse(c, statusCode, "GetUsers", err)
		return
	}

	// Формируем ответ
	response := dto.UsersListResponse{
		Data: users,
		Meta: dto.PaginationMeta{
			Page:  page,
			Limit: limit,
			Total: total,
			Pages: int(math.Ceil(float64(total) / float64(limit))),
		},
	}

	c.JSON(http.StatusOK, response)
}

// createUser godoc
// @Summary Создание нового пользователя
// @Description Создает нового пользователя с указанными данными
// @Tags users
// @Accept json
// @Produce json
// @Param input body dto.UserCreateRequest true "Данные пользователя"
// @Security BearerAuth
// @Success 201 {object} models.User
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 403 {object} dto.ErrorResponse
// @Failure 409 {object} dto.ErrorResponse "Пользователь с таким email/логином уже существует или роль не найдена"
// @Failure 500 {object} dto.ErrorResponse
// @Router /v1/users [post]
func (s *Server) createUser(c *gin.Context) {
	userMeta, err := utils.GetUserMeta(c)
	if err != nil {
		s.ErrorResponse(c, http.StatusBadRequest, "utils.GetUserMeta", err)
		return
	}

	var req dto.UserCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		s.ErrorResponse(c, http.StatusBadRequest, "ShouldBindJSON", err)
		return
	}

	// Валидация данных
	if err := req.ValidateUserCreateRequest(); err != nil {
		s.ErrorResponse(c, http.StatusBadRequest, "ValidateUserCreateRequest", err)
		return
	}

	// Создаем пользователя через usecase
	user, err := s.u.Create(c.Request.Context(), userMeta, req)
	if err != nil {
		if errors.Is(err, models.ErrUserExists) {
			s.ErrorResponse(c, http.StatusConflict, "Create", err)
			return
		} else if errors.Is(err, models.ErrRoleNotFound) {
			s.ErrorResponse(c, http.StatusConflict, "Create", err)
			return
		}
		s.ErrorResponse(c, http.StatusInternalServerError, "Create", err)
		return
	}

	c.JSON(http.StatusCreated, user)
}

// getUser godoc
// @Summary Получение информации о пользователе
// @Description Возвращает полную информацию о пользователе по его ID
// @Tags users
// @Param userId path string true "ID пользователя"
// @Security BearerAuth
// @Success 200 {object} models.User
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 403 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse "Пользователь не найден"
// @Failure 500 {object} dto.ErrorResponse
// @Router /v1/users/{userId} [get]
func (s *Server) getUser(c *gin.Context) {
	log := s.logger.With(wsl.Label("method", "getUser"))

	userMeta, err := utils.GetUserMeta(c)
	if err != nil {
		s.ErrorResponse(c, http.StatusBadRequest, "utils.GetUserMeta", err)
		return
	}
	userID := c.Param("userId")
	if userID == "" {
		s.ErrorResponse(c, http.StatusBadRequest, "param", errors.New("userID is empty"))
		return
	}

	// Получаем пользователя из usecase
	user, err := s.u.GetByID(c.Request.Context(), userMeta, userID)
	if err != nil {
		if errors.Is(err, models.ErrUserNotFound) {
			s.ErrorResponse(c, http.StatusNotFound, "GetByID", err)
			return
		}
		log.ErrorContext(c.Request.Context(), "Failed to get user", wsl.Err(err))
		s.ErrorResponse(c, http.StatusInternalServerError, "GetByID", err)
		return
	}

	c.JSON(http.StatusOK, user)
}

// getProfile godoc
// @Summary Получение информации о числе пользователй по ролям
// @Description Получение информации сколько в системе SA, CA и CO
// @Tags users
// @Security BearerAuth
// @Success 200 {object} dto.CountUserByRoleResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 403 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /v1/users/count-by-role [get]
func (s *Server) getCountUsersByRole(c *gin.Context) {
	userMeta, err := utils.GetUserMeta(c)
	if err != nil {
		s.ErrorResponse(c, http.StatusBadRequest, "utils.GetUserMeta", err)
		return
	}

	// Получение информации сколько в системе SA, CA и CO
	countUsers, err := s.u.GetCountUsersByRole(c.Request.Context(), userMeta)
	if err != nil {
		s.ErrorResponse(c, http.StatusInternalServerError, "GetByID", err)
		return
	}
	res := dto.CountUserByRoleResponse{
		CountSA: countUsers.CountSA,
		CountCA: countUsers.CountCA,
		CountCO: countUsers.CountCO,
	}
	c.JSON(http.StatusOK, res)
}
