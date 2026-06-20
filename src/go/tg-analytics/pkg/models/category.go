package models

import (
	"fmt"
)

// CategoryType представляет тип категории
type CategoryType string

const (
	CategoryTypeNegative CategoryType = "negative" // Отрицательная
	CategoryTypePositive CategoryType = "positive" // Положительная
	CategoryTypeNeutral  CategoryType = "neutral"  // Нейтральная
)

func (t CategoryType) String() string {
	switch t {
	case CategoryTypeNegative:
		return "Отрицательная"
	case CategoryTypePositive:
		return "Положительная"
	case CategoryTypeNeutral:
		return "Нейтральная"
	default:
		return "Неизвестный тип"
	}
}

// Category представляет категорию в базе данных
type Category struct {
	ID   uint         `gorm:"primaryKey;autoIncrement" json:"id"`
	Name string       `gorm:"type:varchar(255);not null;uniqueIndex" json:"name"`
	Type CategoryType `gorm:"type:varchar(20);not null;default:'neutral'" json:"type"`
}

// Предопределенные категории с уникальными ID
var (
	Unknown                   = Category{ID: 1, Name: "Неизвестный класс", Type: CategoryTypeNeutral}
	AggressionRacismTerrorism = Category{ID: 2, Name: "Агрессия, расизм, терроризм", Type: CategoryTypeNegative}
	Botnets                   = Category{ID: 3, Name: "Ботнеты", Type: CategoryTypeNegative}
	WebMail                   = Category{ID: 4, Name: "Веб-почта", Type: CategoryTypePositive}
	LeisureEntertainment      = Category{ID: 5, Name: "Досуг и развлечения", Type: CategoryTypeNegative}
	OnlineStores              = Category{ID: 6, Name: "Интернет магазины", Type: CategoryTypeNegative}
	ComputerGames             = Category{ID: 7, Name: "Компьютерные игры", Type: CategoryTypeNegative}
	Cryptomining              = Category{ID: 8, Name: "Криптомайнинг", Type: CategoryTypeNegative}
	Drugs                     = Category{ID: 9, Name: "Наркотики", Type: CategoryTypeNegative}
	PornographySex            = Category{ID: 10, Name: "Порнография и секс", Type: CategoryTypeNegative}
	ProxyAnonymizers          = Category{ID: 11, Name: "Прокси и анонимайзеры", Type: CategoryTypeNegative}
	BannedSitesRegistry       = Category{ID: 12, Name: "Реестр запрещенных сайтов", Type: CategoryTypeNegative}
	AdultSites                = Category{ID: 13, Name: "Сайты для взрослых", Type: CategoryTypeNegative}
	VirusDistributionSites    = Category{ID: 14, Name: "Сайты, распространяющие вирусы", Type: CategoryTypeNegative}
	SocialNetworks            = Category{ID: 15, Name: "Социальные сети", Type: CategoryTypeNegative}
	TorrentsP2P               = Category{ID: 16, Name: "Торренты и Р2Р-сети", Type: CategoryTypeNegative}
	FileArchives              = Category{ID: 17, Name: "Файловые архивы", Type: CategoryTypeNegative}
	MoviesVideoOnline         = Category{ID: 18, Name: "Фильмы и видео онлайн", Type: CategoryTypeNegative}
	Phishing                  = Category{ID: 19, Name: "Фишинг", Type: CategoryTypeNegative}
	ChatsMessengers           = Category{ID: 20, Name: "Чаты и мессенджеры", Type: CategoryTypeNegative}
	Cryptojacking             = Category{ID: 21, Name: "Криптоджекинг", Type: CategoryTypeNegative}
	Advertising               = Category{ID: 22, Name: "Реклама", Type: CategoryTypeNegative}
	OnlineGames               = Category{ID: 23, Name: "Онлайн-игры", Type: CategoryTypeNegative}
	GamingPlatforms           = Category{ID: 24, Name: "Игровые платформы", Type: CategoryTypeNegative}
	Malware                   = Category{ID: 25, Name: "Вредоносное ПО", Type: CategoryTypeNegative}
	Gambling                  = Category{ID: 26, Name: "Азартные игры", Type: CategoryTypeNegative}
	DepressiveContentSuicide  = Category{ID: 27, Name: "Депрессивный контент", Type: CategoryTypeNegative}
	AlcoholTobacco            = Category{ID: 28, Name: "Алкоголь и табак", Type: CategoryTypeNegative}
	PositiveCategory          = Category{ID: 29, Name: "Положительная категория", Type: CategoryTypePositive}
	AllowedCategory           = Category{ID: 30, Name: "Разрешенный ресурс", Type: CategoryTypePositive}
)

// Все категории для удобного доступа
var PredefinedCategories = []Category{
	Unknown,
	AggressionRacismTerrorism,
	Botnets,
	WebMail,
	LeisureEntertainment,
	OnlineStores,
	ComputerGames,
	Cryptomining,
	Drugs,
	PornographySex,
	ProxyAnonymizers,
	BannedSitesRegistry,
	AdultSites,
	VirusDistributionSites,
	SocialNetworks,
	TorrentsP2P,
	FileArchives,
	MoviesVideoOnline,
	Phishing,
	ChatsMessengers,
	Cryptojacking,
	Advertising,
	OnlineGames,
	GamingPlatforms,
	Malware,
	Gambling,
	DepressiveContentSuicide,
	AlcoholTobacco,
	PositiveCategory,
	AllowedCategory,
}

// String возвращает строковое представление категории
func (c Category) String() string {
	return c.Name
}

// IsValid проверяет, является ли категория допустимой
func (c Category) IsValid() bool {
	for _, cat := range PredefinedCategories {
		if cat.ID == c.ID {
			return true
		}
	}
	return false
}

// IsNegative проверяет, является ли категория отрицательной
func (c Category) IsNegative() bool {
	return c.Type == CategoryTypeNegative
}

// IsPositive проверяет, является ли категория положительной
func (c Category) IsPositive() bool {
	return c.Type == CategoryTypePositive
}

// IsNeutral проверяет, является ли категория нейтральной
func (c Category) IsNeutral() bool {
	return c.Type == CategoryTypeNeutral
}

// GetType возвращает тип категории в виде строки
func (c Category) GetType() string {
	return c.Type.String()
}

// ParseContentCategory парсит строку в категорию
func ParseContentCategory(input string) (Category, error) {
	for _, cat := range PredefinedCategories {
		if cat.Name == input {
			return cat, nil
		}
	}
	return Category{}, fmt.Errorf("недопустимое значение категории: %q", input)
}

// ParseContentCategoryByID парсит ID в категорию
func ParseContentCategoryByID(id uint) (Category, error) {
	for _, cat := range PredefinedCategories {
		if cat.ID == id {
			return cat, nil
		}
	}
	return Category{}, fmt.Errorf("недопустимый ID категории: %d", id)
}
