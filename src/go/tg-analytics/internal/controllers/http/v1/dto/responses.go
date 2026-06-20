package dto

import (
	"tg-an/internal/models"
	"time"
)

type SuccessResponse struct {
	Message string `json:"message"`
}

type ErrorResponse struct {
	Error string `json:"error"` // Описание ошибки
}

// PaginationMeta содержит метаданные пагинации
// Используется вместе со списками данных для навигации по страницам
type PaginationMeta struct {
	Page  int   `json:"page"`  // Текущая страница (начинается с 1)
	Limit int   `json:"limit"` // Количество элементов на странице
	Total int64 `json:"total"` // Общее количество элементов
	Pages int   `json:"pages"` // Общее количество страниц
}

type ListResponse struct {
	Data interface{}    `json:"data"` // Список однотипных данных
	Meta PaginationMeta `json:"meta"` // Метаданные пагинации
}

type RequestStatResponse struct {
	Type  string             `json:"type"`  // Тип данных
	Data  models.RequestStat `json:"data"`  // Список однотипных данных
	Count uint               `json:"count"` // Число точек с данными
}

type TrafficStatResponse struct {
	Data  models.TrafficStat `json:"data"`
	Count uint               `json:"count"`
}

type DeviceStatResponse struct {
	Time  []time.Time                `json:"time"`
	Data  []models.DeviceRequestStat `json:"data"`
	Count uint                       `json:"count"`
}

type GetAnomaliesResponse struct {
	HostCount     uint                   `json:"host_count"`
	BlockCount    uint                   `json:"block_count"`
	HostAnomalies []models.HostAnomalies `json:"host_anomalies"`
}

type GetProhActivityResponse struct {
	TimeSince time.Time `json:"time_since"`
	Count     uint      `json:"count"`
	Data      []uint    `json:"data"`
}

// GetDetectionsResponseV2 представляет ответ со списком выявлений (версия 2)
type GetDetectionsResponseV2 struct {
	Data []Detection_v2 `json:"data"`
	Meta PaginationMeta `json:"meta"`
}
