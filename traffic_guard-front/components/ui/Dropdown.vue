<template>
  <div class="dropdown">
    <div class="dropdown__container" ref="dropdownRef">
      <button
        type="button"
        class="dropdown__button"
        :class="{ 'dropdown__button--error': validationError, 'dropdown__button--disabled': disabled }"
        @click="!disabled ? open = !open : null"
      >
        <span v-if="modelValue.name !== ''" class="dropdown__selected-text" :class="[{ 'dropdown__selected-text--error': validationError }]">
          <p :class="props.selectedTextClass" :style="props.selectedTextStyle">{{ modelValue.name }} <span @click.stop="selectItem({id: '', name: ''})"><CrossIcon class="dropdown__selected-cross-icon"/></span></p>
          
        </span>
        <span v-else class="dropdown__selected-text--placeholder">{{ mainPlaceholder }}</span>
        <span class="dropdown__icon-container">
          <CaretDownIcon class="dropdown__icon" :class="{ 'dropdown__icon--rotated': open, 'dropdown__icon--disabled': disabled }" aria-hidden="true" />
        </span>
      </button>

      <div v-if="validationError" class="validation-error-text">
        {{ validationError }}
      </div>

      <div v-if="open" class="dropdown__options" ref="dropdownList" @scroll="onScroll">
        <div v-if="search" class="dropdown__search-container">
          <input
            type="text"
            class="dropdown__search-input"
            :placeholder="searchPlaceholder"
            v-model="searchQuery"
            @input="debouncedSearch"
          />
        </div>
        <ul class="dropdown__list">
          <template v-if="props.loading === true">
            <li class="dropdown__option">
              <span class="dropdown__option-text">
                Загрузка...
              </span>
            </li>
          </template>
          <template v-else-if="props.error !== '' && props.error !== undefined">
            <li class="dropdown__option">
              <span class="dropdown__option-text">
                {{ props.error }}
              </span>
            </li>
          </template>
          <template v-else-if="filteredItems.length > 0">
            <li
              v-for="item in filteredItems"
              :key="item.id"
              class="dropdown__option"
              :class="{ 'dropdown__option--selected dropdown__option--active': isSelected(item) }"
              @click="selectItem(item)"
            >
              <span class="dropdown__option-text" :class="{ 'dropdown__option-text--selected': isSelected(item) }">
                {{ item.name }}
              </span>

              <span v-if="isSelected(item)" class="dropdown__check-icon-container">
                <CheckIcon class="dropdown__check-icon" aria-hidden="true" />
              </span>
            </li>
          </template>
          <template v-else>
            <li class="dropdown__option">
              <span class="dropdown__option-text">
                Элементы отсутствуют
              </span>
            </li>
          </template>

        </ul>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { DropdownItem } from '~/types/dropdown';
import CaretDownIcon from '~/assets/img/caret-down.svg'
import CheckIcon from '~/assets/img/check.svg'
import CrossIcon from '~/assets/img/cross.svg'
import type { StyleValue } from 'vue';


const props = defineProps<{
  items: DropdownItem[];
  modelValue: DropdownItem;
  search?: boolean;
  localSearch: boolean;
  searchPlaceholder: string;
  mainPlaceholder?: string;
  loading?: boolean;
  error?: string;
  disabled?: boolean;
  validationError?: string
  selectedTextClass?: string
  selectedTextStyle?: StyleValue
}>()

const emit = defineEmits<{
  (e: 'update:modelValue', value: DropdownItem): void
  (e: 'searchOut', searchQuery: string): void;
  (e: 'scrollEnd'): void;
}>()

const open = ref(false)
const searchQuery = ref('')
let searchTimeout: number | null = null

const filteredItems = computed(() => {
  const items = props.items ?? [];

  const query = searchQuery.value?.trim().toLowerCase()

  const shouldSearch =
    props.search &&
    query &&
    query.length > 0 &&
    props.localSearch

  if (!shouldSearch) return items

  return items.filter(({ name }) =>
    name?.trim().toLowerCase().includes(query)
  )
})


const debouncedSearch = () => {
  if (searchTimeout) {
    clearTimeout(searchTimeout)
  }

  searchTimeout = setTimeout(() => {
    searchTimeout = null
  }, 300) as unknown as number
}

const isActive = (item: DropdownItem) => {
  return props.modelValue === item
}

const isSelected = (item: DropdownItem) => {
  return props.modelValue.id === item.id
}

const selectItem = (item: DropdownItem) => {
  emit('update:modelValue', item)
  open.value = false
  searchQuery.value = ''
}

const dropdownList = ref<HTMLElement | null>(null);

const onScroll = () => {
  if (!dropdownList.value || props.loading || props.localSearch) return;

  const scrollBottom = dropdownList.value.scrollTop + dropdownList.value.clientHeight;
  const scrollHeight = dropdownList.value.scrollHeight;

  const distanceToBottom = scrollHeight - scrollBottom;

  if (distanceToBottom < 80) {
    emit('scrollEnd');
  }
};

const dropdownRef = ref<HTMLElement | null>(null)

const onClickOutside = (event: MouseEvent) => {
  if (open.value && dropdownRef.value && !dropdownRef.value.contains(event.target as Node)) {
    open.value = false
  }
}

onMounted(() => {
  document.addEventListener('click', onClickOutside)
})

onBeforeUnmount(() => {
  document.removeEventListener('click', onClickOutside)
})


watch(searchQuery, () => {
  if (!props.localSearch) {
    emit('searchOut', searchQuery.value);
  }
})
</script>

<style lang="scss">
.dropdown {
  position: relative;
  width: 100%;
}

.dropdown__container {
  position: relative;
  width: 100%;
}

.dropdown__button {
  position: relative;
  width: 100%;
  height: 2.25rem;
  padding: 0.375rem 0.75rem;
  text-align: left;
  cursor: pointer;
  color: var(--color-typo-primary);
  background-color: var(--color-white);
  outline: 1px solid var(--color-border-primary);
  border-radius: var(--radius-sm);
  border: 0;
}

.dropdown__button--disabled {
  cursor: not-allowed;
  background-color: #EDF1F2;
  color: #B8C1C5;
}

.dropdown__button--error {
  outline-color: #EB5757;
}

.dropdown__selected-text_badge {
  padding: 0.125rem 0.375rem;
  background-color: var(--color-gray-disabled);
  display: flex;
  gap: 0.25rem;
  align-items: center;
  font-size: 0.75rem;
  font-weight: 500;
  line-height: 1rem;
}

.dropdown__selected-text {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  width: fit-content;
  color: #3F3F46;
  font-size: 0.875rem;
  font-weight: 400;
  line-height: 1.5rem;
}

.dropdown__selected-text p {
  display: flex;
  gap: var(--size-1);
  align-items: center;
}

.dropdown__selected-text--error {
  background-color: #F5E6E6;
}

.dropdown__selected-cross-icon {
  display: flex;
  justify-content: center;
  align-items: center;
  cursor: pointer;
  width: 0.75rem;
  height: 0.75rem;
  padding-top: 1px;
}

.dropdown__selected-text--placeholder {
  color: var(--color-typo-placeholder);
  font-size: var(--font-size-sm);
}

.dropdown__icon-container {
  position: absolute;
  top: 0;
  bottom: 0;
  right: 0.75rem;
  display: flex;
  align-items: center;
  justify-content: center;
  pointer-events: none;
}

.dropdown__icon {
  width: 1rem;
  height: 1rem;
  color: var(--color-typo-primary);
  transition: transform 0.3s ease;
}

.dropdown__icon--rotated {
  transform: rotate(180deg);
}

.dropdown__icon--disabled {
  color: #B4BDC1;
}

.dropdown__chevron {
  height: 1.25rem;
  width: 1.25rem;
  color: #9ca3af;
}

.dropdown__options {
  position: absolute;
  z-index: 100;
  margin-top: 0.25rem;
  max-height: 18.75rem;
  width: 100%;
  overflow: auto;
  border-radius: 0.375rem;
  background-color: white;
  box-shadow: 0 10px 15px -3px rgba(0, 0, 0, 0.1), 0 4px 6px -2px rgba(0, 0, 0, 0.05);
  border: 1px solid rgba(0, 0, 0, 0.1);
}

.dropdown__search-container {
  padding: 0.5rem;
  border-bottom: 1px solid #e5e7eb;
}

.dropdown__search-input {
  width: 100%;
  padding: 0.375rem 0.5rem;
  border-radius: 0.25rem;
  border: 1px solid #d1d5db;
  font-size: 0.875rem;
}

.dropdown__search-input:focus {
  outline: none;
  border-color: color-mix(in srgb, var(--color-blue-main), white 45%);
  box-shadow: 0 0 0 2px rgba(79, 70, 229, 0.2);
}

.dropdown__list {
  padding: 4px 0;
  list-style: none;
  margin: 0;
}

.dropdown__option {
  position: relative;
  cursor: pointer;
  user-select: none;
  padding: 0.5rem 0.75rem;
  padding-right: 2.25rem;
}

.dropdown__option--active {
  background-color: var(--color-blue-main);
  color: white;
}

.dropdown__option-text {
  display: block;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  color: var(--color-typo-primary);
  font-size: var(--font-size-sm);
}

.dropdown__option-text--selected {
  font-weight: 600;
  color: var(--color-white);
}

.dropdown__check-icon-container {
  position: absolute;
  top: 0;
  bottom: 0;
  right: 0.75rem;
  display: flex;
  align-items: center;
  justify-content: center;
  pointer-events: none;
}

.dropdown__check-icon {
  width: 1rem;
  height: 1rem;
  color: white;
}

.dropdown__option--active .dropdown__check {
  color: white;
}

.dropdown__empty {
  padding: 0.5rem 0.75rem;
}

.validation-error-text {
  padding-top: 0.5rem;
  font-weight: 400;
  font-size: 0.75rem;
  line-height: 0.75rem;
  color: #E5381A;
}

.category-badge {
  max-width: fit-content;
  padding: 0.25rem 0.5rem;
  border-radius: 0.375rem;
  font-size: 0.75rem;
  line-height: 1rem;
  font-weight: 500;

  &--aggression { background: #FEE2E2; color: #DC2626; }
  &--racism { background: #FECACA; color: #B91C1C; }
  &--terrorism { background: #FEE2E2; color: #DC2626; }
  &--drugs { background: #FECACA; color: #B91C1C; }
  &--malware { background: #FEE2E2; color: #EF4444; }
  &--malicious-software { background: #FECACA; color: #DC2626; }

  &--botnets { background: #FFEDD5; color: #EA580C; }
  &--cryptomining { background: #FDBA74; color: #C2410C; }
  &--cryptojacking { background: #FFEDD5; color: #D97706; }
  &--phishing { background: #FEF3C7; color: #D97706; }

  &--pornography { background: #F3E8FF; color: #9333EA; }
  &--adult { background: #E9D5FF; color: #7E22CE; }
  &--gambling { background: #F3E8FF; color: #7C3AED; }

  &--proxy { background: #DBEAFE; color: #2563EB; }
  &--torrents { background: #BFDBFE; color: #1D4ED8; }

  &--leisure { background: #DCFCE7; color: #16A34A; }
  &--games { background: #BBF7D0; color: #15803D; }
  &--positive { background: #BBF7D0; color: #15803D; }
  &--online-games { background: #DCFCE7; color: #15803D; }
  &--gaming-platforms { background: #BBF7D0; color: #16A34A; }
  &--movies { background: #DCFCE7; color: #059669; }

  &--social { background: #E0F2FE; color: #0284C7; }
  &--chats { background: #BAE6FD; color: #0369A1; }
  &--webmail { background: #E0F2FE; color: #0EA5E9; }

  &--shopping { background: #E0E7FF; color: #4F46E5; }
  &--advertising { background: #C7D2FE; color: #4338CA; }

  &--files { background: #F3F4F6; color: #6B7280; }
  &--blocked { background: #E5E7EB; color: #4B5563; }
  &--additional { background: #F3F4F6; color: #6B7280; }

  &--depressive { background: #FCE7F3; color: #DB2777; }
  
  &--alcohol { background: #FEF3C7; color: #D97706; }

  &--unknown { background: #f3f4f6; color: #6b7280; }

  &--reverse {
    &--aggression { background: #DC2626; color: #DC2626; }
    &--racism { background: #B91C1C; color: #B91C1C; }
    &--terrorism { background: #DC2626; color: #DC2626; }
    &--drugs { background: #B91C1C; color: #B91C1C; }
    &--malware { background: #EF4444; color: #EF4444; }
    &--malicious-software { background: #DC2626; color: #DC2626; }

    &--botnets { background: #EA580C; color: #EA580C; }
    &--cryptomining { background: #C2410C; color: #C2410C; }
    &--cryptojacking { background: #D97706; color: #D97706; }
    &--phishing { background: #D97706; color: #D97706; }

    &--pornography { background: #9333EA; color: #9333EA; }
    &--adult { background: #7E22CE; color: #7E22CE; }
    &--gambling { background: #7C3AED; color: #7C3AED; }

    &--proxy { background: #2563EB; color: #2563EB; }
    &--torrents { background: #1D4ED8; color: #1D4ED8; }

    &--leisure { background: #16A34A; color: #16A34A; }
    &--games { background: #15803D; color: #15803D; }
    &--positive { background: #15803D; color: #15803D; }
    &--online-games { background: #15803D; color: #15803D; }
    &--gaming-platforms { background: #16A34A; color: #16A34A; }
    &--movies { background: #059669; color: #059669; }

    &--social { background: #0284C7; color: #0284C7; }
    &--chats { background: #0369A1; color: #0369A1; }
    &--webmail { background: #0EA5E9; color: #0EA5E9; }

    &--shopping { background: #4F46E5; color: #4F46E5; }
    &--advertising { background: #4338CA; color: #4338CA; }

    &--files { background: #6B7280; color: #6B7280; }
    &--blocked { background: #4B5563; color: #4B5563; }
    &--additional { background: #6B7280; color: #6B7280; }

    &--depressive { background: #DB2777; color: #DB2777; }

    &--alcohol { background: #D97706; color: #D97706; }

    &--unknown { background: #6b7280; color: #6b7280; }
  }
}

.session-status-badge_greeen {
  max-width: fit-content;
  padding: 0.25rem 0.5rem;
  border-radius: 0.375rem;
  font-size: 0.75rem;
  line-height: 1rem;
  font-weight: 500;
  color: #008236;
  background-color: #DCFCE7;
}

.session-status-badge_red {
  max-width: fit-content;
  padding: 0.25rem 0.5rem;
  border-radius: 0.375rem;
  font-size: 0.75rem;
  line-height: 1rem;
  font-weight: 500;
  color: #C10007;
  background-color: #FFE2E2;
}

.session-status-badge_yellow {
  max-width: fit-content;
  padding: 0.25rem 0.5rem;
  border-radius: 0.375rem;
  font-size: 0.75rem;
  line-height: 1rem;
  font-weight: 500;
  color: #894B00;
  background-color: #FEF9C2;
}
</style>
