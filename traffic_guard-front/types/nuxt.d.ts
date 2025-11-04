import type { UserActions } from "~/types/user";
import type { NotificationManager } from "~/types/notifications";
import { ApiClient } from './api';

declare module "#app" {
  interface NuxtApp {
    $api: ApiClient;
    $notifications: NotificationManager;
    $user: UserActions;
  }
}

declare module "vue" {
  interface ComponentCustomProperties {
    $api: ApiClient;
    $notifications: NotificationManager;
    $user: UserActions;
  }
}
