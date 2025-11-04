<template>
  <button class="filter-button">
    <div class="filter-button__icon">
      <FilterIcon />
    </div>
    <span class="filter-button__text">Фильтр</span>
  </button>
</template>

<script setup lang="ts">
import FilterIcon from "~/assets/img/filter-icon.svg"


</script>

<style lang="scss">
.filter-button {
  background-color: #FFFFFF;
  box-shadow: 0px 1px 2px 0px #0000000D;
  padding: 0.5rem 0.75rem;
  display: flex;
  align-items: center;
  gap: 0.375rem;
  border: none;
  border-radius: 0.375rem;
  cursor: pointer;
}

.filter-button__icon {
  display: flex;
  justify-content: center;
  align-items: center;
  width: 1.25rem;
  height: 1.25rem;
  color: #A1A1AA;
}

.filter-button__icon svg {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.filter-button__text {
  font-weight: 500;
  font-size: 1rem;
  line-height: 1.25rem;
  color: #3F3F46;
}  
</style>