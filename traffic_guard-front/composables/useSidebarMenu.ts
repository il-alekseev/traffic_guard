// @ts-ignore
import HouseIcon from '~/assets/img/home.svg?component';
// @ts-ignore
import NetworkShield from '~/assets/img/security-shield.svg?component';
// @ts-ignore
import ChartRose from '~/assets/img/chart-rose.svg?component';
// @ts-ignore
import CodeIcon from '~/assets/img/code.svg?component';
// @ts-ignore
import UsersGroup from '~/assets/img/user-group.svg?component';
// @ts-ignore
import LogsIcon from '~/assets/img/logs.svg?component';


export function useSidebarMenu() {
  const config = useRuntimeConfig();
  const grafanaUrl = `${config.public.grafana.host}:${config.public.grafana.port}`

  const getMenuForRole = (role: string) => {
    const fullMenu = [
      {
        id: "main",
        label: "Dashboard",
        to: "/dashboard",
        iconComponent: HouseIcon,
        showFor: ["KSU", "KSF"],
      },
      {
        id: "detections",
        label: "Выявления",
        to: "/detections",
        iconComponent: NetworkShield,
        showFor: ["KSU", "KSF"],
      },
      {
        id: "reports",
        label: "Отчеты",
        to: "/reports",
        iconComponent: ChartRose,
        showFor: ["KSU", "KSF"],
      },
      {
        id: "sessions",
        label: "Сессии",
        to: "/sessions",
        iconComponent: CodeIcon,
        showFor: ["KSU", "KSF"],
      },
      {
        id: "users",
        label: "Пользователи",
        to: "/users",
        iconComponent: UsersGroup,
        showFor: ["KSU", "KSF"],
      },
      {
        id: "logs",
        label: "Логи",
        to: "/logs",
        iconComponent: LogsIcon,
        showFor: ["KSU", "KSF"],
      },
      // {
      //   id: "grafana",
      //   label: "Статистика",
      //   href: grafanaUrl,
      //   iconComponent: ChartBar,
      //   external: true,
      //   showFor: ["KSU", "KSF"],
      // },
    ]

    return fullMenu.filter((item) => item.showFor.includes(role))
  }

  return {
    getMenuForRole,
    grafanaUrl,
  }
}

