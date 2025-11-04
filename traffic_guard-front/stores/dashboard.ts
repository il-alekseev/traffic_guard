import { defineStore } from "pinia";
import { useNuxtApp } from "#app";
import type { defaultResponse } from "~/types/api";
import { useUserStore } from "./user";
import { getTokenHeaders } from "~/helpers";
import type { DashboardState, DashboardTraffic, TopCategories, TopDetections } from "~/types/dashboard";


export const useDashboardStore = defineStore("dashboard", {
  state: (): DashboardState => {
    return {
        topCategories: null,
        topDetections: null,
        trafic: null,
        requests: null,
        anomalies: null,
        events: null,
        proh_activity: null,
        devicesState: null
    };
  },

  getters: {
    // devices: (state) => state.devices,
  },

  actions: {
    async fetchTopCategories(from: string = 'now-10m', to: string = 'now', count: number = 5, hostname?: string): Promise<TopCategories> {

      const userStore = useUserStore();
      try {
        await userStore.ensureValidToken();
      } catch {
        userStore.clearToken();
        throw new Error(
          "Не удалось получить топ категорий. Пользователь неавторизован",
        );
      }

      const token = useCookie('auth_token').value;

      if (!token) {
        userStore.clearToken();
        throw new Error(
          "Не удалось получить топ категорий. Пользователь неавторизован",
        );
      }

      try {
        const { $api } = useNuxtApp();

        const params: Record<string, string | number> = {
          from,
          to,
          count,
          ...(hostname ? { hostname } : {}),
        };

        const result = await $api.get<TopCategories>('/dashboards/top-categories', {
          params,
          // ...getTokenHeaders(token)
        });


        if (result) {
          this.topCategories = result;
          const data: TopCategories = [
            {
              category: 'Агрессия, расизм, терроризм',
              count: 3044
            },
            {
              category: 'Наркотики',
              count: 2450
            },
            {
              category: 'Порнография и секс',
              count: 1300
            },
            {
              category: 'Торренты и P2P-сети',
              count: 1199
            },
            {
              category: 'Прокси и анонимайзеры',
              count: 485
            },
          ];
          this.topCategories = data;
          return result;
        } else {
          throw new Error("Не удалось получить топ категорий");
        }
      } catch (error: any) {
        throw new Error(error.message || "Ошибка при получении топа категорий");
      }
    },

    async fetchDetections(from: string = 'now-10m', to: string = 'now', count: number = 5, hostname?: string): Promise<TopDetections> {

      const userStore = useUserStore();
      try {
        await userStore.ensureValidToken();
      } catch {
        userStore.clearToken();
        throw new Error(
          "Не удалось получить топ выявлений. Пользователь неавторизован",
        );
      }

      const token = useCookie('auth_token').value;

      if (!token) {
        userStore.clearToken();
        throw new Error(
          "Не удалось получить топ выявлений. Пользователь неавторизован",
        );
      }

      try {
        const { $api } = useNuxtApp();

        const params: Record<string, string | number> = {
          from,
          to,
          count,
          ...(hostname ? { hostname } : {}),
        };

        const result = await $api.get<TopDetections>('/dashboards/top-unresolved_detections', {
          params,
          ...getTokenHeaders(token)
        });


        if (result) {
          this.topDetections = result;
          const data = [
            {
              domain: "twitter",
              requests_after: 600,
              requests_all: 1000,
              requests_before: 400
            },
            {
              domain: "youtube",
              requests_after: 123,
              requests_all: 443,
              requests_before: 320
            },
            {
              domain: "facebook",
              requests_after: 190,
              requests_all: 200,
              requests_before: 10
            },
            {
              domain: "telegram",
              requests_after: 302,
              requests_all: 512,
              requests_before: 210
            },
            {
              domain: "discord",
              requests_after: 9,
              requests_all: 10,
              requests_before: 1
            },
          ];
          this.topDetections = data;
          return result;
        } else {
          throw new Error("Не удалось получить топ выявлений");
        }
      } catch (error: any) {
        throw new Error(error.message || "Ошибка при получении топа выявлений");
      }
    },

    async fetchTraffic(from: string = 'now-10m', to: string = 'now', count: number = 5, hostname?: string): Promise<DashboardTraffic> {

      const userStore = useUserStore();
      try {
        await userStore.ensureValidToken();
      } catch {
        userStore.clearToken();
        throw new Error(
          "Не удалось получить статистику трафика. Пользователь неавторизован",
        );
      }

      const token = useCookie('auth_token').value;

      if (!token) {
        userStore.clearToken();
        throw new Error(
          "Не удалось получить статистику трафика. Пользователь неавторизован",
        );
      }

      try {
        const { $api } = useNuxtApp();

        const params: Record<string, string | number> = {
          from,
          to,
          count,
          ...(hostname ? { hostname } : {}),
        };

        const result = await $api.get<DashboardTraffic>('/dashboards/traffic', {
          params,
          ...getTokenHeaders(token)
        });


        if (result) {
          this.trafic = result;
          return result;
        } else {
          throw new Error("Не удалось получить статистику трафика");
        }
      } catch (error: any) {
        throw new Error(error.message || "Ошибка при получении статистику трафика");
      }
    },
  },
});
