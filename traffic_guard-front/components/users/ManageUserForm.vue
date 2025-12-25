<template>
  <div class="create-user">
    <BaseAlert v-if="error">{{ error }}</BaseAlert>
    <form @submit.prevent="submitForm" class="create-user__form">
      <div class="create-user-form__group">

        <div
          v-if="props.manageType === 'update' && !props.isProfileForm"
          class="info-block"
        >
          <p class="info-block__title">
            ID
          </p>
          <p class="info-block__text">
            {{ props.userData?.user_id }}
          </p>
        </div>

        <div
          v-if="props.manageType === 'update'"
          class="info-block"
        >
          <p class="info-block__title">
            Логин
          </p>
          <p class="info-block__text">
            {{ props.userData?.login }}
          </p>
        </div>

        <BaseInput
          v-if="props.manageType === 'create'"
          id="login"
          name="login"
          label="Логин"
          :showRequiredLabel="true"
          v-model="form.login"
          :error="errors.login"
        />

        <BaseInput
          id="email"
          name="email"
          label="Email"
          :showRequiredLabel="true"
          type="email"
          v-model="form.email"
          :error="errors.email"
        />

        <BaseInput
          id="last_name"
          name="last_name"
          label="Фамилия"
          :showRequiredLabel="true"
          v-model="form.last_name"
          :error="errors.last_name"
        />

        <BaseInput
          id="first_name"
          name="first_name"
          label="Имя"
          :showRequiredLabel="true"
          v-model="form.first_name"
          :error="errors.first_name"
        />

        <BaseInput
          id="patronymic"
          :required="false"
          name="patronymic"
          label="Отчество"
          :showRequiredLabel="true"
          v-model="form.patronymic"
          :error="errors.patronymic"
        />

        <div v-if="!isProfileForm" class="base-form-field">
          <label for="role" class="base-label">Роль</label>
          <Dropdown
              :search="true"
              :localSearch="true"
              id="role"
              :items="allowedRoles"
              searchPlaceholder="Поиск..."
              v-model="form.role"
          />
          <p v-if="errors.role" class="user-form__error-text">
            {{ errors.role }}
          </p>
        </div>

        <div v-if="form.role.id === 'CA' && !isProfileForm" class="base-form-field">
          <label for="device" class="base-label">NGFW</label>
          <Dropdown
            :search="true"
            :localSearch="true"
            id="device"
            :items="allowedDevices"
            searchPlaceholder="Поиск..."
            v-model="form.device"
            :selectedTextClass="'ngfw-bagde'"
            :selectedTextStyle="{
              backgroundColor: generateColor(form.device.id).background,
              color: generateColor(form.device.id).color
            }"
          />
        </div>

        <div class="extras" v-if="props.showExtraButtons && props.userData !== undefined && props.userData?.user_id">
          <p class="extras__title">Дополнительно</p>
          <div class="extras__btns">
            <BaseButton
              class="extras__btn"
              type="button"
              variant="extra"
              :disabled="loadingResetPass"
              loadingText="Сброс пароля..."
              @click="resetPassword"
            >
              Сбросить пароль
            </BaseButton>
            <BaseButton
              class="extras__btn"
              type="button"
              variant="danger"
              :disabled="loadingDelete"
              loadingText="Удаление..."
              @click="deleteUser"
            >
              Удалить пользователя
            </BaseButton>
          </div>
        </div>
      </div>

      <div class="create-user-form__actions">
        <BaseButton
          class="cancel-btn"
          type="button"
          variant="secondary"
          :disabled="loading"
          loadingText="Cохранение..."
          @click="emit('close')"
        >
          Отменить
        </BaseButton>

        <BaseButton
          class="submit-btn"
          type="submit"
          variant="primary"
          :loading="loading"
          loadingText="Загрузка..."
        >
          {{ props.manageType === 'create' ? "Создать пользователя" : "Сохранить изменения" }}
        </BaseButton>
      </div>
    </form>
  </div>
</template>

<script setup lang="ts">
import { useDevicesStore } from '~/stores/devices';
import { useDeviceColors } from '~/composables/useDeviceColors';
import { useUsersControlStore } from '~/stores/usersControl';
import { useNotificationStore } from '~/stores/notification';
import { useUserStore } from "~/stores/user";
import type { User, UserForm } from '~/types/user';
import type { DropdownItem } from '~/types/dropdown';
import { getOnlyDeviceName, getOnlyRole, getRoleDisplayName, getTempPassword } from '~/helpers';
import BaseInput from '~/components/ui/BaseInput.vue';
import BaseButton from '~/components/ui/BaseButton.vue';
import BaseAlert from '~/components/ui/BaseAlert.vue';
import Dropdown from '~/components/ui/Dropdown.vue';


const { generateColor } = useDeviceColors();

const emit = defineEmits<{
  close: [];
  success: [opType: 'create' | 'update', user: User, password: string | null ];
  resetPassword: [userId: string, password: string];
  deleteUser: [userId: string];
}>();

const props = withDefaults(
  defineProps<{
    manageType: 'create' | 'update';
    userData: User | null;
    showExtraButtons?: boolean;
    isProfileForm?: boolean
  }>(),
  {
    manageType: 'create',
    userData: null,
    showExtraButtons: true,
    isProfileForm: false
  }
);

interface createUserDTO {
  login: string;
  first_name: string;
  last_name: string;
  patronymic: string;
  role: string;
  email: string;
  password: string;
}
const allowedRoles = ref<{id: string, name: string}[]>([{id: 'SA', name: 'Системный администратор'}, {id: 'CA', name: 'Администратор узла NGFW'}]);

const usersControlStore = useUsersControlStore();
const userStore = useUserStore();
const loading = ref(false);
const loadingResetPass = ref(false);
const loadingDelete = ref(false);
const error = ref('');

const deviceStore = useDevicesStore();
const allowedDevices = computed<DropdownItem[]>(() => {
  const devices = deviceStore.devices;
  if (!devices) return [];
  const result = devices.map((device) => {
    return {
      id: device,
      name: device
    }
  })
  return result;
})


const form = reactive<UserForm>({
  login: '',
  first_name: '',
  last_name: '',
  patronymic: '',
  role: {
    id: '',
    name: '',
  },
  device: {
    id: '',
    name: '',
  },
  email: '',
});

const errors = reactive({
  login: '',
  first_name: '',
  last_name: '',
  patronymic: '',
  role: '',
  email: '',
});

const resetPassword = async () => {
  if (loading.value || loadingDelete.value || loadingResetPass.value) return;
  if (!props.userData || !props.userData.user_id) return;

  loadingResetPass.value = true;

  try {
    const password = getTempPassword();
    const result = await usersControlStore.resetPassword(props.userData.user_id, password);
    if (result) {
      emit('resetPassword', props.userData.user_id, password);
    } else {
      error.value = 'Не удалось удалить пользователя';
    }

  } catch (e) {
    error.value = 'Произошла ошибка при удалении пользователя';
  }
}

const deleteUser = async () => {
  if (loading.value || loadingDelete.value || loadingResetPass.value) return;
  if (!props.userData || !props.userData.user_id) return;

  loadingDelete.value = true;

  try {
    const result = await usersControlStore.deleteUser(props.userData.user_id);
    if (result) {
      emit('deleteUser', props.userData.user_id);
    } else {
      error.value = 'Не удалось удалить пользователя';
    }

  } catch (e) {
    error.value = 'Произошла ошибка при удалении пользователя';
  }
}

const validateForm = () => {
  let isValid = true;

  (Object.keys(errors) as (keyof typeof errors)[]).forEach(key => {
    errors[key] = '';
  });

if (!form.login.trim()) {
    errors.login = 'Логин обязателен для заполнения';
    isValid = false;
  } else if (form.login.length < 4) {
    errors.login = 'Логин должен содержать не менее 4 символов';
    isValid = false;
  } else if (!/^[a-zA-Z0-9_-]+$/.test(form.login)) {
    errors.login = 'Логин может содержать только латинские буквы, цифры и символ подчеркивания';
    isValid = false;
  }

  if (!form.first_name.trim()) {
    errors.first_name = 'Имя обязательно для заполнения';
    isValid = false;
  } else if (form.first_name.length < 2) {
    errors.first_name = 'Имя должно содержать не менее 2 символов';
    isValid = false;
  } else if (form.first_name.length > 30) {
    errors.first_name = 'Имя должно содержать не более 30 символов';
    isValid = false;
  }

  if (!form.last_name.trim()) {
    errors.last_name = 'Фамилия обязательна для заполнения';
    isValid = false;
  } else if (form.last_name.length < 2) {
    errors.last_name = 'Фамилия должна содержать не менее 2 символов';
    isValid = false;
  } else if (form.last_name.length > 30) {
    errors.last_name = 'Фамилия должна содержать не более 30 символов';
    isValid = false;
  }

  if (form.patronymic.trim().length > 30) {
    errors.last_name = 'Фамилия должно содержать не более 30 символов';
    isValid = false;
  }

  if (form.email.trim().length > 0 && !/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(form.email)) {
    errors.email = 'Введите корректный email';
    isValid = false;
  }

  if (!form.role.id) {
    errors.role = 'Выберите роль';
    isValid = false;
  }

  return isValid;
};

const createUser = async () => {
  if (loading.value) return;

  error.value = '';

  if (!validateForm()) {
    return;
  }

  loading.value = true;

  try {
    const password = getTempPassword();
    const fullRole = form.role.id === 'SA' ? form.role.id : form.role.id + '-' + form.device.id

    const data: createUserDTO = {
      login: form.login.trim(),
      first_name: form.first_name.trim(),
      last_name: form.last_name.trim(),
      patronymic: form.patronymic.trim(),
      role: fullRole,
      email: form.email.trim(),
      password: password,
    }

    const result = await usersControlStore.createUser(data);
    if (result) {
      emit('success', 'create', result, password);
    } else {
      error.value = 'Не удалось создать пользователя';
    }
  } catch (err: any) {
    if (err?.message === '409') {
      error.value = 'Логин и/или Email заняты';
    } else {
      error.value = 'Произошла ошибка при создании пользователя';
    }
  } finally {
    loading.value = false;
  }
};

const editProfile = async () => {
  if (loading.value) return;

  error.value = '';

  if (props.userData === null || !props.userData.user_id) {
    error.value = 'Произошла ошибка при обновлении пользователя';
    return;
  }

  if (!validateForm()) {
    return;
  }

  loading.value = true;

  try {
    const data: Partial<User> = {
      first_name: form.first_name,
      last_name: form.last_name,
      patronymic: form.patronymic,
      role: form.role.id,
      email: form.email,
    }

    const result = await userStore.updateUserInfo(data);

    if (result) {
      const userObject: User = {
        ...data,
        login: props.userData.login,
        user_id: props.userData.user_id,
        is_super_admin: props.userData.is_super_admin,
        created_at: props.userData.created_at,
        is_need_to_change_password: props.userData.is_need_to_change_password,
      } as User;

      emit('success', 'update', userObject, null);
      const notify = useNotificationStore();
      notify.notify('Профиль изменён', 'success');
    } else {
      error.value = 'Не удалось обновить профиль';
    }
  } catch (err) {
    console.error('Ошибка при обновлении профиля:', err);
    error.value = 'Произошла ошибка при обновлении профиля';
  } finally {
    loading.value = false;
  }
}

const editUser = async () => {
  if (loading.value) return;

  error.value = '';

  if (props.userData === null || !props.userData.user_id) {
    error.value = 'Произошла ошибка при обновлении пользователя';
    return;
  }

  if (!validateForm()) {
    return;
  }

  loading.value = true;

  try {
    const data: Partial<User> = {
      first_name: form.first_name,
      last_name: form.last_name,
      patronymic: form.patronymic,
      role: form.role.id,
      email: form.email,
    }

    const result = await usersControlStore.updateUserInfo(props.userData.user_id, data);

    if (result) {
      const userObject: User = {
        ...data,
        login: props.userData.login,
        user_id: props.userData.user_id,
        is_super_admin: props.userData.is_super_admin,
        created_at: props.userData.created_at,
        is_need_to_change_password: props.userData.is_need_to_change_password,
      } as User;

      emit('success', 'update', userObject, null);
      const notify = useNotificationStore();
      notify.notify('Пользователь изменён', 'success');
    } else {
      error.value = 'Не удалось обновить пользователя';
    }
  } catch (err) {
    console.error('Ошибка при обновлении пользователя:', err);
    error.value = 'Произошла ошибка при обновлении пользователя';
  } finally {
    loading.value = false;
  }
}

const submitForm = () => {
  if (props.manageType === 'create') {
    createUser();
  } else if (props.manageType === 'update') {
    if (props.isProfileForm) editProfile();
    else editUser();
  } else return;
}

const initForm = async () => {
  if (props.userData === null) return;

  const onlyRole = getOnlyRole(props.userData.role);
  const onlyDevice = getOnlyDeviceName(props.userData.role)

  form.login = props.userData.login || '';
  form.email = props.userData.email || '';
  form.first_name = props.userData.first_name || '';
  form.last_name = props.userData.last_name || '';
  form.patronymic = props.userData.patronymic || '';
  form.role.id = onlyRole;
  form.role.name = getRoleDisplayName(onlyRole) || '';
  form.device.id = onlyDevice || '';
  form.device.name = onlyDevice || '';
}

onMounted(() => {
  props.manageType === 'update'
    ? initForm()
    : null
})
</script>

<style scoped>
.create-user {
  height: 100%;
}

.create-user .base-alert {
  margin: 1.5rem 1.5rem;
}

.info-block {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.info-block__title {
  font-weight: 500;
  font-size: 1rem;
  line-height: 1.5rem;
  color: #3F3F46;
}

.info-block__text {
  font-weight: 400;
  font-size: 1rem;
  line-height: 1.25rem;
  color: #71717A;
}

.create-user__form {
  display: flex;
  flex-direction: column;
  height: 100%;
}

.create-user-form__group {
  padding: 0 1.5rem;
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
}

.create-user-form__actions {
  display: flex;
  gap: 0.5rem;
  margin-top: auto;
  padding: 1rem;
  border-top: 1px solid #E4E4E7;
  justify-content: flex-end;
}

.extras {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.extras__title {
  font-weight: 500;
  font-size: 1rem;
  line-height: 1.5rem;
  color: #3F3F46;
}

.extras__btns {
  display: flex;
  gap: 0.5rem;
}

.cancel-btn {
  max-width: 5.8125rem;
}

.submit-btn {
  max-width: 13rem;
}

.user-form__error-text {
  font-size: var(--font-size-xs);
  color: #EF4444;
}
</style>
