package category

import "fmt"

// ContentCategory представляет категории контента
type ContentCategory int

const (
	AggressionRacismTerrorism ContentCategory = iota + 1
	Botnets
	WebMail
	LeisureEntertainment
	OnlineStores
	ComputerGames
	Cryptomining
	Drugs
	PornographySex
	ProxyAnonymizers
	BannedSitesRegistry
	AdultSites
	VirusDistributionSites
	SocialNetworks
	TorrentsP2P
	FileArchives
	MoviesVideoOnline
	Phishing
	ChatsMessengers
	Cryptojacking
	Advertising
	OnlineGames
	GamingPlatforms
	Malware
	Gambling
	DepressiveContentSuicide
	AlcoholTobacco
)

// String возвращает строковое представление категории
func (c ContentCategory) String() string {
	switch c {
	case AggressionRacismTerrorism:
		return "Агрессия, расизм, терроризм"
	case Botnets:
		return "Ботнеты"
	case WebMail:
		return "Веб-почта"
	case LeisureEntertainment:
		return "Досуг и развлечения"
	case OnlineStores:
		return "Интернет-магазины"
	case ComputerGames:
		return "Компьютерные игры"
	case Cryptomining:
		return "Криптомайнинг"
	case Drugs:
		return "Наркотики"
	case PornographySex:
		return "Порнография и секс"
	case ProxyAnonymizers:
		return "Прокси и анонимайзеры"
	case BannedSitesRegistry:
		return "Реестр запрещенных сайтов"
	case AdultSites:
		return "Сайты для взрослых"
	case VirusDistributionSites:
		return "Сайты, распространяющие вирусы"
	case SocialNetworks:
		return "Социальные сети"
	case TorrentsP2P:
		return "Торренты и Р2Р-сети"
	case FileArchives:
		return "Файловые архивы"
	case MoviesVideoOnline:
		return "Фильмы и видео онлайн"
	case Phishing:
		return "Фишинг"
	case ChatsMessengers:
		return "Чаты и мессенджеры"
	case Cryptojacking:
		return "Криптоджекинг"
	case Advertising:
		return "Реклама"
	case OnlineGames:
		return "Онлайн-игры"
	case GamingPlatforms:
		return "Игровые платформы"
	case Malware:
		return "Вредоносное ПО"
	case Gambling:
		return "Азартные игры"
	case DepressiveContentSuicide:
		return "Депрессивный контент и суицид"
	case AlcoholTobacco:
		return "Алкоголь, табак"
	default:
		return "Неизвестная категория"
	}
}

// IsValid проверяет, является ли значение категории допустимым
func (c ContentCategory) IsValid() bool {
	switch c {
	case AggressionRacismTerrorism, Botnets, WebMail, LeisureEntertainment,
		OnlineStores, ComputerGames, Cryptomining, Drugs, PornographySex,
		ProxyAnonymizers, BannedSitesRegistry, AdultSites, VirusDistributionSites,
		SocialNetworks, TorrentsP2P, FileArchives, MoviesVideoOnline, Phishing,
		ChatsMessengers, Cryptojacking, Advertising, OnlineGames, GamingPlatforms,
		Malware, Gambling, DepressiveContentSuicide, AlcoholTobacco:
		return true
	default:
		return false
	}
}

func ParseContentCategory(input string) (ContentCategory, error) {
	switch input {
	case "Агрессия, расизм, терроризм":
		return AggressionRacismTerrorism, nil
	case "Ботнеты":
		return Botnets, nil
	case "Веб-почта":
		return WebMail, nil
	case "Досуг и развлечения":
		return LeisureEntertainment, nil
	case "Интернет-магазины":
		return OnlineStores, nil // Используем OnlineStores из String() метода
	case "Компьютерные игры":
		return ComputerGames, nil
	case "Криптомайнинг":
		return Cryptomining, nil
	case "Наркотики":
		return Drugs, nil
	case "Порнография и секс":
		return PornographySex, nil
	case "Прокси и анонимайзеры":
		return ProxyAnonymizers, nil
	case "Реестр запрещенных сайтов":
		return BannedSitesRegistry, nil
	case "Сайты для взрослых":
		return AdultSites, nil
	case "Сайты, распространяющие вирусы":
		return VirusDistributionSites, nil
	case "Социальные сети":
		return SocialNetworks, nil
	case "Торренты и Р2Р-сети":
		return TorrentsP2P, nil
	case "Файловые архивы":
		return FileArchives, nil
	case "Фильмы и видео онлайн":
		return MoviesVideoOnline, nil
	case "Фишинг":
		return Phishing, nil
	case "Чаты и мессенджеры":
		return ChatsMessengers, nil
	case "Криптоджекинг":
		return Cryptojacking, nil
	case "Реклама":
		return Advertising, nil
	case "Онлайн-игры":
		return OnlineGames, nil
	case "Игровые платформы":
		return GamingPlatforms, nil
	case "Вредоносное ПО":
		return Malware, nil
	case "Азартные игры":
		return Gambling, nil
	case "Депрессивный контент и суицид":
		return DepressiveContentSuicide, nil
	case "Алкоголь, табак":
		return AlcoholTobacco, nil
	default:
		return 0, fmt.Errorf("недопустимое значение категории: %q (ожидается: %s)", input, getExpectedCategories())
	}
}

// Вспомогательная функция для получения списка ожидаемых категорий
func getExpectedCategories() string {
	categories := []string{
		"Агрессия, расизм, терроризм",
		"Ботнеты",
		"Веб-почта",
		"Досуг и развлечения",
		"Интернет-магазины",
		"Компьютерные игры",
		"Криптомайнинг",
		"Наркотики",
		"Порнография и секс",
		"Прокси и анонимайзеры",
		"Реестр запрещенных сайтов",
		"Сайты для взрослых",
		"Сайты, распространяющие вирусы",
		"Социальные сети",
		"Торренты и Р2Р-сети",
		"Файловые архивы",
		"Фильмы и видео онлайн",
		"Фишинг",
		"Чаты и мессенджеры",
		"Криптоджекинг",
		"Реклама",
		"Онлайн-игры",
		"Игровые платформы",
		"Вредоносное ПО",
		"Азартные игры",
		"Депрессивный контент и суицид",
		"Алкоголь, табак",
	}

	// Формируем строку с перечислением категорий
	result := ""
	for i, cat := range categories {
		if i > 0 {
			result += ", "
		}
		result += cat
	}
	return result
}

// AllCategoryStrings возвращает слайс всех категорий в строковом представлении
func AllCategoryStrings() []string {
	categories := []ContentCategory{
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
	}

	result := make([]string, len(categories))
	for i, cat := range categories {
		result[i] = cat.String()
	}
	return result
}
