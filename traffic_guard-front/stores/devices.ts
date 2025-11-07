import { defineStore } from "pinia";
import { useNuxtApp } from "#app";
import type { defaultResponse } from "~/types/api";
import { useUserStore } from "./user";
import { getTokenHeaders } from "~/helpers";
import type { Device, DevicesState } from "~/types/devices";


export const useDevicesStore = defineStore("devices", {
  state: (): DevicesState => {
    return {
      devices: null,
    };
  },

  getters: {
    // devices: (state) => state.devices,
  },

  actions: {
    async fetchDevices(): Promise<Device[]> {

      const userStore = useUserStore();
      try {
        await userStore.ensureValidToken();
      } catch {
        userStore.clearToken();
        throw new Error(
          "Не удалось получить список девайсов. Пользователь неавторизован",
        );
      }

      const token = useCookie('auth_token').value;

      if (!token) {
        userStore.clearToken();
        throw new Error(
          "Не удалось получить список девайсов. Пользователь неавторизован",
        );
      }

      try {
        const { $api } = useNuxtApp();

        const devices = await $api.get<Device[]>('/devices', {
          ...getTokenHeaders(token)
        });


        if (devices) {
          this.devices = devices;
          return devices;
        } else {
          throw new Error("Не удалось получить список девайсов");
        }
      } catch (error: any) {
        throw new Error(error.message || "Ошибка при получении списка девайсов");
      }
    }
  },
});
