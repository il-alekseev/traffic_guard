import { defineStore } from "pinia";
import { useNuxtApp } from "#app";
import type { defaultResponse } from "~/types/api";
import { useUserStore } from "./user";
import { getTokenHeaders } from "~/helpers";
import type { Categories } from "~/types/categories";
import type { CategoriesControlState } from "~/types/categoriesControl";


export const useCategoriesControlStore = defineStore("categoriesControl", {
  state: (): CategoriesControlState => {
    return {
      categories: null,
    };
  },

  getters: {
    // categories: (state) => state.categories,
  },

  actions: {
    async fetchCategories(): Promise<Categories[]> {
      const userStore = useUserStore();
      try {
        await userStore.ensureValidToken();
      } catch {
        userStore.clearToken();
        throw new Error(
          "Не удалось получить список категорий. Пользователь неавторизован",
        );
      }

      const token = useCookie('auth_token').value;

      if (!token) {
        userStore.clearToken();
        throw new Error(
          "Не удалось получить список категорий. Пользователь неавторизован",
        );
      }

      try {
        const { $api } = useNuxtApp();

        const categories = await $api.get<Categories[]>('/categories', {
          ...getTokenHeaders(token)
        });


        if (categories) {
          this.categories = categories;
          return categories;
        } else {
          throw new Error("Не удалось получить список категорий");
        }
      } catch (error: any) {
        throw new Error(error.message || "Ошибка при получении списка категорий");
      }
    }
  },
});
