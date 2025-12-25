<template>
  <div class="pagination">
    <div class="pagination__left">
      <div class="pagination__select-wrapper" ref="selectWrapperRef">
        <button
          class="pagination__select-button"
          @click="toggleDropdown"
        >
          <span class="pagination__select-label">
            {{ selectedItemsCount }} на странице
          </span>
          <ArrowDownIcon :class="['pagination__select-icon', { 'pagination__select-icon--open': isDropdownOpen }]" alt="dropdown icon"
          />
        </button>

        <div
          v-if="isDropdownOpen"
          class="pagination__dropdown"
          :class="{ 'pagination__dropdown--top': dropdownDirection === 'top' }"
        >
          <button
            v-for="count in allowedItemsCount"
            :key="count"
            class="pagination__dropdown-item"
            :class="{ 'pagination__dropdown-item--active': count === selectedItemsCount }"
            @click="selectItemsCount(count)"
          >
            {{ count }} на странице
          </button>
        </div>
      </div>
    </div>

    <div class="pagination__right">
      <button
        class="pagination__nav-button"
        :disabled="currentPage === 1"
        @click="goToPreviousPage"
        aria-label="Previous page"
      >
        <ArrowLeftIcon alt="previous page" />
      </button>

      <span class="pagination__info">
        Страница {{ currentPage }} - {{ getEndItem() }} из {{ total }}
      </span>

      <button
        class="pagination__nav-button"
        :disabled="currentPage === totalPages"
        @click="goToNextPage"
        aria-label="Next page"
      >
        <ArrowLeftIcon alt="next page" class="pagination__nav-icon--right"/>
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import ArrowLeftIcon from '~/assets/img/arrow-left.svg'
import ArrowDownIcon from '~/assets/img/arrow-down.svg'

interface Props {
  total: number
  allowedItemsCount: number[]
  currentPage: number
  selectedItemsCount: number
  totalPages: number
}

interface Emits {
  'update:currentPage': [page: number]
  'update:selectedItemsCount': [count: number]
}

const emit = defineEmits<Emits>()

const isDropdownOpen = ref(false)

const selectItemsCount = (count: number) => {
  emit('update:selectedItemsCount', count)
  isDropdownOpen.value = false
}

const goToPreviousPage = () => {
  if (props.currentPage > 1) {
    emit('update:currentPage', props.currentPage - 1)
  }
}

const goToNextPage = () => {
  if (props.currentPage < props.totalPages) {
    emit('update:currentPage', props.currentPage + 1)
  }
}

const props = defineProps<Props>()

const getEndItem = () => {
  return Math.min(props.currentPage * props.selectedItemsCount, props.total)
}

const dropdownDirection = ref<'bottom' | 'top'>('bottom')
const selectWrapperRef = ref<HTMLElement | null>(null)

const checkDropdownDirection = () => {
  if (!selectWrapperRef.value) return

  const wrapper = selectWrapperRef.value
  const rect = wrapper.getBoundingClientRect()
  const spaceBelow = window.innerHeight - rect.bottom
  const dropdownHeight = Math.min(props.allowedItemsCount.length * 44 + 8, 200)

  if (spaceBelow < dropdownHeight) {
    dropdownDirection.value = 'top'
  } else {
    dropdownDirection.value = 'bottom'
  }
}

const toggleDropdown = () => {
  isDropdownOpen.value = !isDropdownOpen.value
  if (isDropdownOpen.value) {
    checkDropdownDirection()
  }
}
</script>

<style scoped lang="scss">
.pagination {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0.75rem;
  background-color: #FFFFFF;
  border-radius: 6px;

  &__left {
    flex: 0 0 auto;
  }

  &__right {
    display: flex;
    align-items: center;
    gap: 1rem;
  }

  &__select-wrapper {
    position: relative;
  }

  &__select-button {
    display: flex;
    align-items: center;
    gap: 0.25rem;
    padding: 0.5rem 0.75rem;
    background-color: #ffffff;
    border: 1px solid #D4D4D8;
    box-shadow: 0px 1px 2px 0px #0000000D;
    border-radius: 6px;
    cursor: pointer;
    font-size: 0.875rem;
    line-height: 1.25rem;
    font-weight: 500;
    color: #3F3F46;
    transition: all 0.2s ease;
  }

  &__select-label {
    white-space: nowrap;
  }

  &__select-icon {
    width: 1.25rem;
    height: 1.25rem;
    color: #A1A1AA;
    transition: transform 0.2s ease;

    &--open {
      transform: rotate(180deg);
    }
  }

  &__dropdown {
    position: absolute;
    top: calc(100% + 6px);
    left: 0;
    background-color: #ffffff;
    border: 1px solid #D4D4D8;
    border-radius: 6px;
    z-index: 10;
    min-width: 100%;
    overflow: hidden;

    &--top {
      top: auto;
      bottom: calc(100% + 4px);
    }
  }

  &__dropdown-item {
    display: block;
    width: 100%;
    padding: 0.75rem;
    border: none;
    background-color: #ffffff;
    text-align: left;
    font-size: 0.75rem;
    line-height: 1.25rem;
    font-weight: 500;
    color: #3F3F46;
    cursor: pointer;
    transition: background-color 0.2s ease;

    &--active {
      opacity: 0.5;
    }
  }

  &__nav-button {
    display: flex;
    align-items: center;
    justify-content: center;
    width: 2.25rem;
    height: 2.25rem;
    padding: 0;
    background-color: #ffffff;
    border: 1px solid #D4D4D8;
    border-radius: 6px;
    cursor: pointer;
    transition: all 0.2s ease;

    &:disabled {
      opacity: 0.5;
      cursor: not-allowed;
    }

    svg {
      width: 1.25rem;
      height: 1.25rem;
      color: #A1A1AA;
    }
  }

  &__nav-icon--right {
    transform: rotate(180deg);
  }

  &__info {
    font-size: 0.875rem;
    line-height: 1.25rem;
    font-weight: 500;
    color: 3F3F46;
    white-space: nowrap;
  }
}
</style>