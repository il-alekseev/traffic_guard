<template>
  <div class="users-page">
    <div class="users__header">
      <h1 class="users__title page-title">Управление пользователями</h1>
      <div class="users__toolbar">
        <BaseSearch
          v-model="searchQuery"
          class="users__search"
          placeholder="Поиск"
          @search="applyFilters"
        />
        <FilterButton/>
        <BaseButton
          @click.stop="openCreateUserModal"
          type="button"
          variant="primary"
          class="users__create-button">
          <template #icon>
            <PlusIcon />
          </template>
          Создать пользователя
        </BaseButton>
      </div>
    </div>

    <ErrorBlock v-if="fetchError !== ''" :fetch-error="fetchError" />

    <BaseTable
      v-if="fetchError === ''"
      :items="users"
      :columns="columns"
      :loading="loading"
      :total-items="totalUsers"
      :total-pages="totalPages"
      :current-page="currentPage"
      :items-per-page="itemsPerPage"
      :allowedItemsCount="[11, 33, 66]"
      prefix="users"
      :show-actions="true"
      empty-message="Пользователи не найдены"
      loading-message="Загрузка пользователей..."
      item-key="user_id"
      @page-change="handleChangePage"
      @set-items-per-page="handleChangeItemsPerPage"
      @action-click="openEditUserModal"
    >
      <template #cell-login="{ value }">
        <span class="users__login-cell">{{ value }}</span>
      </template>

      <template #cell-role="{ value }">
        <span class="users__role-cell">{{ getRoleDisplayName(value) }}</span>
      </template>

      <template #cell-context="{ item }">
        <span
          class="users__table-cell__badge users__context-cell"
          v-bind="getOnlyRole(item.role) !== 'SA' 
                  ? { style: generateColor(getOnlyDeviceName(item.role)) } 
                  : { style: {background: 'transparent'} }"
        >
          {{ getOnlyDeviceName(item.role) || '' }}
        </span>
      </template>

      <template #cell-full_name="{ item }">
        {{ item.last_name }} {{ item.first_name }} {{ item.patronymic || '' }}
      </template>

      <template #cell-email="{ value }">
        <span class="users__email-cell">{{ value }}</span>
      </template>

      <template #action-button="{ item }">
        <EditDataIcon class="table__action-icon users__action-icon" />
      </template>
    </BaseTable>

    <SideModal
      v-model="showCreateUserModal"
      title="Создать пользователя"
      @close="closeCreateUserModal"
    >
      <ManageUserForm
        :manageType="'create'"
        :user-data="null"
        @close="closeCreateUserModal"
        @success="manageUser"
      />
    </SideModal>

    <SideModal
      v-model="showEditUserModal"
      title="Редактировать пользователя"
      @close="closeEditUserModal"
    >
      <ManageUserForm
        :manageType="'update'"
        :userData="selectedUser"
        @close="closeEditUserModal"
        @success="manageUser"
        @resetPassword="handlePasswordReset"
        @deleteUser="handleUserDeleted"
      />
    </SideModal>

    <TempPasswordUserModal
      v-if="showTempPasswordUserModal"
      :user="selectedUser"
      :password="tempPassword"
      :type="tempPasswordType"
      @close="closeTempPasswordUserModal"
    />
  </div>
</template>

<script setup lang="ts">
import {definePageMeta} from '#imports';
import { useUsersControlStore } from '~/stores/usersControl';
import { useUserStore } from '~/stores/user';
import type { User, UsersTable } from '~/types/user';
import BaseSearch from '~/components/ui/BaseSearch.vue';
import BaseButton from '~/components/ui/BaseButton.vue';
import SideModal from '~/components/ui/SideModal.vue';
import ManageUserForm from '~/components/users/ManageUserForm.vue';
import TempPasswordUserModal from '~/components/users/TempPasswordUserModal.vue';
import FilterButton from '~/components/ui/FilterButton.vue';
import ErrorBlock from '~/components/ui/ErrorBlock.vue';
import BaseTable from '~/components/ui/BaseTable.vue';
import PlusIcon from "~/assets/img/plus.svg";
import EditDataIcon from '~/assets/img/edit.svg';
import { getRoleDisplayName, getOnlyRole, getOnlyDeviceName } from "~/helpers";
import { useDeviceColors } from '~/composables/useDeviceColors';


definePageMeta({
  layout: 'dashboard',
  middleware: ['auth']
});

const { generateColor } = useDeviceColors();

const route = useRoute();
const router = useRouter();

const usersControlStore = useUsersControlStore();
const userStore = useUserStore();

const loading = ref(true);
const fetchError = ref('');

const selectedUser = ref<User | null>(null);
const tempPassword = ref('');
const showCreateUserModal = ref(false);
const showEditUserModal = ref(false);
const showTempPasswordUserModal = ref(false);
const tempPasswordType = ref<'create' | 'reset'>('create')

const users = ref<User[]>([]);
const columns = [
  { key: 'login', label: 'Логин' },
  { key: 'role', label: 'Роль' },
  { key: 'context', label: 'Контекст' },
  { key: 'full_name', label: 'ФИО' },
  { key: 'email', label: 'Email' }
]
const currentPage = ref(1);
const itemsPerPage = ref(11);
const tableMetaData = ref<{ total: number; pages: number } | null>(null);
const totalUsers = computed(() =>
  tableMetaData.value ? tableMetaData.value.total : users.value.length
);
const totalPages = computed(() =>
  tableMetaData.value ? tableMetaData.value.pages : 1
);

const searchQuery = ref('');

const fetchUsers = async () => {
  loading.value = true;
  fetchError.value = '';

  try {
    const result: UsersTable = await usersControlStore.fetchUsers(
      currentPage.value,
      itemsPerPage.value,
      searchQuery.value,
    );

    if (result) {
      users.value = result.data;
      tableMetaData.value = result.meta;
    } else {
      users.value = [];
    }
  } catch (error) {
    console.error('Ошибка при загрузке пользователей:', error);
    fetchError.value = 'Произошла ошибка при загрузке пользователей';
  } finally {
    loading.value = false;
  }
};

const toggleLockBodyScroll = (isLock: boolean) => {
  window.scrollTo({
    top: 0,
    behavior: "smooth"
  })
  isLock 
    ? document.body.classList.add('no-scroll')
    : document.body.classList.remove('no-scroll')
}

const openCreateUserModal = () => {
  toggleLockBodyScroll(true);
  selectedUser.value = null;
  showCreateUserModal.value = true;
};

const closeCreateUserModal = () => {
  toggleLockBodyScroll(false);
  showCreateUserModal.value = false;
};

const handleUserCreated = (newUser: User, password: string) => {
  usersControlStore.addUser(newUser);
  closeCreateUserModal();
  openTempPasswordUserModal(newUser, password, 'create');
};

const openEditUserModal = (user: User) => {
  toggleLockBodyScroll(true);
  selectedUser.value = { ...user };
  showEditUserModal.value = true;
};

const closeEditUserModal = () => {
  toggleLockBodyScroll(false);
  selectedUser.value = null;
  showEditUserModal.value = false;
};

const handleUserUpdated = async (updatedUser: User) => {
  const userIndex = users.value.findIndex((u) => u.user_id === updatedUser.user_id);
  if (userIndex !== -1) {
    users.value[userIndex] = updatedUser;
  } else {
    await fetchUsers();
  }
  closeEditUserModal();
};

const manageUser = (opType: 'create' | 'update', user: User, password: string | null) => {
  if (opType === 'create' && user && password !== null) {
    handleUserCreated(user, password)
  } else if (opType === 'update' && user) {
    handleUserUpdated(user)
  } else return;
}

const handlePasswordReset = (userId: string, password: string) => {
  closeEditUserModal();
  const user: User | undefined = users.value.find(u => u.user_id === userId);
  if (!user) return;
  openTempPasswordUserModal(user, password, 'reset');
};

const openTempPasswordUserModal = (user: User, password: string, tempPassType: 'create' | 'reset') => {
  toggleLockBodyScroll(true);
  selectedUser.value = { ...user };
  tempPassword.value = password;
  showTempPasswordUserModal.value = true;
  tempPasswordType.value = tempPassType;
};

const closeTempPasswordUserModal = () => {
  toggleLockBodyScroll(false);
  selectedUser.value = null;
  tempPassword.value = '';
  showTempPasswordUserModal.value = false;
};

const handleUserDeleted = (userId: string) => {
  closeEditUserModal();
  users.value = users.value.filter(u => u.user_id !== userId);
  usersControlStore.removeUser(userId);
};

const handleChangePage = (page: number) => {
  if (page < 1 || page > totalPages.value) return;
  currentPage.value = page;
  updateUrlParams();
};

const handleChangeItemsPerPage = (value: number) => {
  itemsPerPage.value = value;
  currentPage.value = 1;
  updateUrlParams();
}

const initFiltersFromUrl = () => {
  const query = route.query;

  currentPage.value = Number(query.page) || 1;
  itemsPerPage.value = Number(query.per_page) || 11;
  searchQuery.value = typeof query.search === 'string' ? query.search : '';
};

const updateUrlParams = () => {
  const query: Record<string, string | number> = {};

  if (currentPage.value > 1) query.page = currentPage.value;
  if (itemsPerPage.value !== 11) query.per_page = itemsPerPage.value;
  if (searchQuery.value.trim()) query.search = searchQuery.value.trim();

  router.replace({ query });
};

const applyFilters = () => {
  currentPage.value = 1;
  updateUrlParams();
};

onMounted(async () => {
  initFiltersFromUrl();
  await fetchUsers();
});
watch(
  () => route.query,
  async (newQuery, oldQuery) => {
    const withoutCreate = (q: typeof newQuery) => {
      const { create, ...rest } = q;
      return JSON.stringify(rest);
    };

    if (withoutCreate(newQuery) === withoutCreate(oldQuery)) {
      return;
    }

    initFiltersFromUrl();
    await fetchUsers();
  },
  { deep: true }
);

</script>

<style scoped>
.users__header {
  display: flex;
  gap: 1.5rem;
  margin-bottom: 1.5rem;
}

.users__title {
  min-width: 27rem;
}

.users__toolbar {
  width: 100%;
  display: flex;
  gap: 0.5rem;
}

.users__search {
  max-width: 47.25rem;
  min-width: 16rem;
  width: 100%;
}

.users__create-button {
  width: 14.5rem;
}

:deep(.users__table) {
  min-width: 77.5rem;
}

:deep(td.users__table-cell-login) {
  font-weight: 500;
  font-size: 1rem;
  line-height: 1.25rem;
  color: #3F3F46;
}

:deep(td.users__table-cell-email) {
  color: #2563EB;
}

:deep(.table__action-icon) {
  width: 1.25rem;
  height: 1.25rem;
  color: #2563EB;
}

:deep(.users__table-cell__badge) {
  font-weight: 500;
  font-size: 0.75rem;
  line-height: 1rem;
  padding: 0.25rem 0.5rem;
  text-transform: uppercase;
  border-radius: 6px;
}

:deep(.users__table-column-login) {
  width: 10%
}
:deep(.users__table-column-role) {
  width: 17.5%
}
:deep(.users__table-column-context) {
  width: 20.5%
}
:deep(.users__table-column-full_name) {
  width: 25%;
}
:deep(.users__table-column-email) {
  width: 25%;
}
:deep(.users__table-column-button) {
  width: 3%;
}

</style>