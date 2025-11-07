import type { SessionStatus, SessionTypes } from "~/types/session";
import type { ActionType } from "~/types/detections";
import type { Categories } from "~/types/categories";

export const getTokenHeaders = (token: string): {
  headers: {
    Authorization: string;
  };} => {
  return {
    headers: {
      Authorization: `Bearer ${token}`,
    },
  };
}

export const getLogoutTokenHeaders = ( token: string, refreshToken: string): {
  headers: {
    Authorization: string;
    'X-Refresh-Token': string;
  };
} => ({
  headers: {
    Authorization: `Bearer ${token}`,
    'X-Refresh-Token': refreshToken,
  },
});


export const getRefreshTokenHeaders = (refreshToken: string): {
  headers: {
    'X-Refresh-Token': string;
  };
} => ({
  headers: {
    'X-Refresh-Token': refreshToken,
  },
});
export const getTempPassword = (): string => {
  const length = 12;

  const sets = {
    lower: 'abcdefghijklmnopqrstuvwxyz',
    upper: 'ABCDEFGHIJKLMNOPQRSTUVWXYZ',
    number: '0123456789',
    special: '!@#$%^&*()-_=+[]{};:,.<>?',
  };

  const all = Object.values(sets).join('');
  const getRandomChar = (str: string) =>
    str[crypto.getRandomValues(new Uint32Array(1))[0] % str.length];

  const base = Object.values(sets).map(getRandomChar);

  const rest = Array.from({ length: length - base.length }, () => getRandomChar(all));

  return [...base, ...rest].sort(() => crypto.getRandomValues(new Uint8Array(1))[0] - 128).join('');
};
export const getFormattedDatetime = (datetime: string | undefined): string => {
  if (datetime === undefined) return '–';
  const date = new Date(datetime);

  return date.toLocaleString('ru-RU', {
    day: '2-digit',
    month: '2-digit',
    year: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
    hour12: false,
  });
};

export function getJwtExpMs(token: string | null | undefined): number | null {
  if (!token) return null;
  try {
    const [, payload] = token.split(".");
    if (!payload) return null;
    const json = JSON.parse(atob(payload.replace(/-/g, "+").replace(/_/g, "/")));
    return typeof json?.exp === "number" ? json.exp * 1000 : null;
  } catch {
    return null;
  }
}

export function getModificatorByCategory(category: Categories): string {
  const categoryMap: Record<Categories, string> = {
    'Неизвестный класс': 'unknown',
    'Положительная категория': 'positive',
    'Агрессия, расизм, терроризм': 'aggression',
    'Ботнеты': 'botnets',
    'Веб-почта': 'webmail',
    'Досуг и развлечения': 'leisure',
    'Интернет магазины': 'shopping',
    'Компьютерные игры': 'games',
    'Криптомайнинг': 'cryptomining',
    'Наркотики': 'drugs',
    'Порнография и секс': 'pornography',
    'Прокси и анонимайзеры': 'proxy',
    'Реестр запрещенных сайтов': 'blocked',
    'Сайты для взрослых': 'adult',
    'Сайты, распространяющие вирусы': 'malware',
    'Социальные сети': 'social',
    'Торренты и P2P-сети': 'torrents',
    'Файловые архивы': 'files',
    'Фильмы и видео онлайн': 'movies',
    'Фишинг': 'phishing',
    'Чаты и мессенджеры': 'chats',
    'Дополнительно': 'additional',
    'Криптоджекинг': 'cryptojacking',
    'Реклама': 'advertising',
    'Онлайн-игры': 'online-games',
    'Игровые платформы': 'gaming-platforms',
    'Вредоносное ПО': 'malicious-software',
    'Азартные игры': 'gambling',
    'Депрессивный контент': 'depressive',
    'Алкоголь и табак': 'alcohol'
  };

  return categoryMap[category] || 'unknown';
}

export function getBadgeClassByStatus(status: SessionStatus): string {
  switch (status) {
    case ('Разрешен'):
      return 'session-status-badge_greeen'
    case ('Запрещен'):
      return 'session-status-badge_red'
    case ('Ожидает'):
      return 'session-status-badge_yellow'
    default:
      return 'session-status-badge_yellow'
  }
}

export const getNgfwBadgeClass = (ngfw: number) => {
  return `sessions__table-cell__badge`;
};

export const getCurrentDateWithOffset = (
  offset: number = 0,
  unit: 'ms' | 's' | 'm' | 'h' | 'd' = 'ms'
): Date => {
  const multipliers = {
    ms: 1,
    s: 1000,
    m: 60 * 1000,
    h: 60 * 60 * 1000,
    d: 24 * 60 * 60 * 1000,
  } as const;

  return new Date(Date.now() + offset * multipliers[unit]);
};

const CATEGORIES_LIST = [
  'Неизвестный класс',
  'Положительная категория',
  'Агрессия, расизм, терроризм',
  'Ботнеты',
  'Веб-почта',
  'Досуг и развлечения',
  'Интернет магазины',
  'Компьютерные игры',
  'Криптомайнинг',
  'Наркотики',
  'Порнография и секс',
  'Прокси и анонимайзеры',
  'Реестр запрещенных сайтов',
  'Сайты для взрослых',
  'Сайты, распространяющие вирусы',
  'Социальные сети',
  'Торренты и P2P-сети',
  'Файловые архивы',
  'Фильмы и видео онлайн',
  'Фишинг',
  'Чаты и мессенджеры',
  'Дополнительно',
  'Криптоджекинг',
  'Реклама',
  'Онлайн-игры',
  'Игровые платформы',
  'Вредоносное ПО',
  'Азартные игры',
  'Депрессивный контент',
  'Алкоголь и табак',
] as const satisfies readonly Categories[];

export function isCategory(value: unknown): value is Categories {
  return CATEGORIES_LIST.includes(value as Categories);
}

const SESSION_STATUSES_LIST = [
  "Разрешен",
  "Запрещен",
  "Ожидает",
  "Аномалия"
]

export function isSessionStatus(value: unknown): value is SessionStatus {
  return SESSION_STATUSES_LIST.includes(value as SessionStatus);
}

const SESSION_TYPES_LIST = [
  "Разрешен",
  "Запрещен",
  "VPN"
]

export function isSessionTypes(value: unknown): value is SessionTypes {
  return SESSION_TYPES_LIST.includes(value as SessionTypes);
}

export const calculatePercent = (value: number, total: number): string => {
  if (total === 0) return "0%";
  
  const percent = (value / total) * 100;
  
  if (percent > 0 && percent < 1) {
    return "<1%";
  }
  
  return `${Math.round(percent)}%`;
};

export const formatTraffic = (bytes: number): string => {
  if (bytes >= 1000000) {
    return `${(bytes / 1000000).toFixed(0)}M`;
  }
  if (bytes >= 1000) {
    return `${(bytes / 1000).toFixed(0)}k`;
  }
  return `${bytes}`;
};

export const getRoleDisplayName = (role: string): string => {
  const onlyRole = getOnlyRole(role);

  const roles: Record<string, string> = {
    'SA': 'Системный администратор',
    'CA': 'Администратор узла NGFW',
  };
  return roles[onlyRole] || role;
};

export const getOnlyRole = (role: string): string => {
  const [onlyRole] = role?.split('-', 1);
  return onlyRole;
};

export const getOnlyDeviceName = (role: string): string => {
  if (role === 'SA') {
    return ''
  }
  return role.split("-")[1];
};

export function makeUTCDate(year: number, month: number, day: number): Date {
  return new Date(Date.UTC(year, month, day))
}

export const isValidDateString = (val: string | null | undefined) => {
  if (!val || typeof val !== 'string') return false;
  const d = new Date(val);
  return !isNaN(d.getTime());
};
