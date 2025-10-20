package category

import "strings"

// ContentCategory представляет категории контента
type ContentCategory int

const (
	AggressionRacismTerrorism ContentCategory = 1 << iota
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
	InternetStores
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
	case InternetStores:
		return "Интернет-магазины"
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

// CategorySet представляет набор категорий
type CategorySet ContentCategory

// Add добавляет категорию в набор
func (cs *CategorySet) Add(category ContentCategory) {
	*cs |= CategorySet(category)
}

// Remove удаляет категорию из набора
func (cs *CategorySet) Remove(category ContentCategory) {
	*cs &^= CategorySet(category)
}

// Has проверяет наличие категории в наборе
func (cs CategorySet) Has(category ContentCategory) bool {
	return (cs & CategorySet(category)) != 0
}

// String возвращает строковое представление набора категорий
func (cs CategorySet) String() string {
	var categories []string
	for i := 0; i < 28; i++ {
		cat := ContentCategory(1 << i)
		if cs.Has(cat) {
			categories = append(categories, cat.String())
		}
	}
	return strings.Join(categories, " | ")
}
