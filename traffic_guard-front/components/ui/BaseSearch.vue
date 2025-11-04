<template>
  <div class="search">
    <span class="search__icon-container">
      <SearchIcon class="search__icon" />
    </span>
    <input
      class="search__input"
      type="text"
      v-model="searchQuery"
      :placeholder="placeholder"
      @keydown.enter="handleEnter"
    />
  </div>
</template>

<script setup lang="ts">
import SearchIcon from "~/assets/img/search.svg";

const props = defineProps<{
  modelValue: string;
  placeholder: string;
}>()

const emit = defineEmits<{
  (e: 'update:modelValue', value: string): void;
  (e: 'search'): void;
}>()

const searchQuery = ref(props.modelValue);

watch(searchQuery, (val) => {
  emit('update:modelValue', val);
});

watch(() => props.modelValue, (val) => {
  searchQuery.value = val;
});

const handleEnter = () => {
  emit('search');
}
</script>

<style>

.search {
  width: 100%;
  position: relative;
  flex: 1;
}

.search__icon-container {
  position: absolute;
  top: 50%;
  left: 0.5rem;
  transform: translateY(-50%);
  display: flex;
  align-items: center;
  justify-content: center;
}

.search__icon {
  width: 1.25rem;
  height: 1.25rem;
}

.search__input {
  width: 100%;
  padding: 0.5rem 0.5rem 0.5rem 2.125rem;
  font-size: 1rem;
  line-height: 1.25rem;
  color: #3F3F46;
  background-color: #FFFFFF;
  border-radius: 0.375rem;
  border: 0;
  transition: all var(--transition);
}

.search__input::placeholder {
  color: #D4D4D8;
}
  
</style>