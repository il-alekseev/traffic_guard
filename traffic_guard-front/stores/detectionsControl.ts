import { defineStore } from "pinia";
import { useNuxtApp } from "#app";
import { useUserStore } from "./user";
import { getTokenHeaders } from "~/helpers";
import type { DetectionControlState, DetectionStats, DetectionTable } from "~/types/detectionsControl";
import type { Categories } from "~/types/categories";


export const useDetectionsControlStore = defineStore("detectionControl", {
  state: (): DetectionControlState => {
    return {
      detections: null,
    };
  },

  getters: {
    // detections: (state) => state.detections,
  },

  actions: {
    async fetchDetections(from: string = 'now-10m', to: string = 'now', page: number = 1, limit: number = 6, _status?: string, category?: Categories, hostname?: string): Promise<DetectionTable> {
      const userStore = useUserStore();
      try {
        await userStore.ensureValidToken();
      } catch {
        userStore.clearToken();
        throw new Error(
          "Не удалось получить список выявлений. Пользователь неавторизован",
        );
      }

      const token = useCookie('auth_token').value;

      if (!token) {
        userStore.clearToken();
        throw new Error(
          "Не удалось получить список выявлений. Пользователь неавторизован",
        );
      }

      try {
        const { $api } = useNuxtApp();

        const params: Record<string, string | number> = {
          from,
          to,
          page,
          limit,
          ...(_status ? { _status } : {}),
          ...(hostname ? { hostname } : {}),
          ...(category ? { category } : {}),
        };

        const detections = await $api.get<DetectionTable>('/detections/', {
          params,
          ...getTokenHeaders(token)
        });


        if (detections) {
          this.detections = detections.data;
          return detections;
        } else {
          throw new Error("Не удалось получить список выявлений");
        }
      } catch (error: any) {
        throw new Error(error.message || "Ошибка при получении списка выявлений");
      }
    },

    async fetchDetectionStats(from: string = 'now-10m', to: string = 'now', hostname?: string, category?: Categories): Promise<DetectionStats> {
      // const userStore = useUserStore();
      // try {
      //   await userStore.ensureValidToken();
      // } catch {
      //   userStore.clearToken();
      //   throw new Error(
      //     "Не удалось получить статистику по выявлениям. Пользователь неавторизован",
      //   );
      // }

      // const token = useCookie('auth_token').value;

      // if (!token) {
      //   userStore.clearToken();
      //   throw new Error(
      //     "Не удалось получить статистику по выявлениям. Пользователь неавторизован",
      //   );
      // }

      try {
        const { $api } = useNuxtApp();

        const params: Record<string, string | number> = {
          from,
          to,
          ...(hostname ? { hostname } : {}),
          ...(category ? { category } : {}),
        };

        const detectionStats = await $api.get<DetectionStats>('/detections/stat', {
          params,
          // ...getTokenHeaders(token)
        });


        if (detectionStats) {
          return detectionStats;
        } else {
          throw new Error("Не удалось получить статистику по выявлениям");
        }
      } catch (error: any) {
        throw new Error(error.message || "Не удалось получить статистику по выявлениям");
      }
    }
  },
});
