import { defineStore } from "pinia";
import { useNuxtApp } from "#app";
import type { defaultResponse } from "~/types/api";
import { useUserStore } from "./user";
import { getTokenHeaders } from "~/helpers";
import type { DashboardAnomalies, DashboardProhActivity, DashboardProhActivityResponse, DashboardRequestObj, DashboardRequests, DashboardRequestsType, DashboardState, DashboardTraffic, TopCategories, TopDetections } from "~/types/dashboard";


export const useDashboardStore = defineStore("dashboard", {
  state: (): DashboardState => {
    return {
        topCategories: null,
        topDetections: null,
        trafic: null,
        requests: {
          allowed: null,
          blocked: null,
          before_block: null,
          pending: null
        },
        anomalies: null,
        logs: null,
        proh_activity: null,
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

        const result = await $api.get<TopCategories>('/analytics/dashboards/top-categories', {
          params,
          ...getTokenHeaders(token)
        });


        if (result) {
          this.topCategories = result;
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

        const result = await $api.get<TopDetections>('/analytics/dashboards/top-unresolved_detections', {
          params,
          ...getTokenHeaders(token)
        });


        if (result) {
          this.topDetections = result;
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

        const result = await $api.get<DashboardTraffic>('/analytics/dashboards/traffic', {
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
        throw new Error(error.message || "Ошибка при получении статистики трафика");
      }
    },

    async fetchRequest(from: string = 'now-10m', to: string = 'now', request_type: DashboardRequestsType = 'allowed', count: number = 10, hostname?: string, setToStore: boolean = true): Promise<DashboardRequestObj> {
      const userStore = useUserStore();
      try {
        await userStore.ensureValidToken();
      } catch {
        userStore.clearToken();
        throw new Error(
          "Не удалось получить статистику запросов. Пользователь неавторизован",
        );
      }

      const token = useCookie('auth_token').value;

      if (!token) {
        userStore.clearToken();
        throw new Error(
          "Не удалось получить статистику запросов. Пользователь неавторизован",
        );
      }

      try {
        const { $api } = useNuxtApp();

        const params: Record<string, string | number> = {
          from,
          to,
          count,
          request_type,
          ...(hostname ? { hostname } : {}),
        };

        const result = await $api.get<DashboardRequestObj>('/analytics/dashboards/requests', {
          params,
          ...getTokenHeaders(token)
        });


        if (result) {
          if (request_type === 'allowed' && setToStore) {
            this.requests.allowed = result;
          } else if (request_type === 'blocked' && setToStore) {
            this.requests.blocked = result;
          } else if (request_type === 'before_block' && setToStore) {
            this.requests.before_block = result
          } else if (request_type === 'pending' && setToStore) {
            this.requests.pending = result;
          }
          return result;
        } else {
          throw new Error("Не удалось получить статистику запросов");
        }
      } catch (error: any) {
        throw new Error(error.message || "Ошибка при получении статистики запросов");
      }
    },

    async fetchRequests(from: string = 'now-10m', to: string = 'now', count: number = 10, hostname?: string): Promise<DashboardRequests> {
      const userStore = useUserStore();
      try {
        await userStore.ensureValidToken();
      } catch {
        userStore.clearToken();
        throw new Error(
          "Не удалось получить статистику запросов. Пользователь неавторизован",
        );
      }

      const token = useCookie('auth_token').value;

      if (!token) {
        userStore.clearToken();
        throw new Error(
          "Не удалось получить статистику запросов. Пользователь неавторизован",
        );
      }

      try {
        await Promise.all([
          this.fetchRequest(from, to, 'allowed', count, hostname),
          this.fetchRequest(from, to, 'blocked', count, hostname),
          this.fetchRequest(from, to, 'before_block', count, hostname),
          this.fetchRequest(from, to, 'pending', count, hostname),
        ]);
        return this.requests;
      } catch (error: any) {
        throw new Error(error.message || "Ошибка при получении статистики запросов");
      }
    },

    async fetchAnomalies(from: string = 'now-10m', to: string = 'now', hostname?: string): Promise<DashboardAnomalies> {

      const userStore = useUserStore();
      try {
        await userStore.ensureValidToken();
      } catch {
        userStore.clearToken();
        throw new Error(
          "Не удалось получить топ аномалий. Пользователь неавторизован",
        );
      }

      const token = useCookie('auth_token').value;

      if (!token) {
        userStore.clearToken();
        throw new Error(
          "Не удалось получить топ аномалий. Пользователь неавторизован",
        );
      }

      try {
        const { $api } = useNuxtApp();

        const params: Record<string, string | number> = {
          from,
          to,
          ...(hostname ? { hostname } : {}),
        };

        const result = await $api.get<DashboardAnomalies>('/analytics/dashboards/anomalies', {
          params,
          ...getTokenHeaders(token)
        });


        if (result) {
          this.anomalies = result;
          return result;
        } else {
          throw new Error("Не удалось получить топ аномалий");
        }
      } catch (error: any) {
        throw new Error(error.message || "Ошибка при получении топа аномалий");
      }
    },

    async fetchProhActivity(year: string, hostname?: string): Promise<DashboardProhActivityResponse> {
      const userStore = useUserStore();
      try {
        await userStore.ensureValidToken();
      } catch {
        userStore.clearToken();
        throw new Error(
          "Не удалось получить график запрещенной активности. Пользователь неавторизован",
        );
      }

      const token = useCookie('auth_token').value;

      if (!token) {
        userStore.clearToken();
        throw new Error(
          "Не удалось получить график запрещенной активности. Пользователь неавторизован",
        );
      }

      try {
        const { $api } = useNuxtApp();

        const params: Record<string, string | number> = {
          year,
          ...(hostname ? { hostname } : {}),
        };

        const result = await $api.get<DashboardProhActivityResponse>('/analytics/dashboards/proh-activity', {
          params,
          ...getTokenHeaders(token)
        });


        if (result) {
          this.proh_activity = result;
          return result;
        } else {
          throw new Error("Не удалось получить график запрещенной активности");
        }
      } catch (error: any) {
        throw new Error(error.message || "Ошибка при получении графика запрещенной активности");
      }
    }
  },
});
