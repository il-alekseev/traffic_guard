<template>
    <button class="download-btn">
      <DownloadIcon />
    </button>
</template>

<script setup lang="ts">
import DownloadIcon from '~/assets/img/download.svg';


</script>



<style>
.download-btn {
  padding: 0.5rem;
  background-color: #FFFFFF;
  box-shadow: 0px 1px 2px 0px #0000000D;
  border-radius: 6px;
  border: none;
  display: flex;
  align-items: center;
  justify-content: center;
}

.download-btn svg {
  width: 20px;
  height: 20px;
  object-fit: contain;
}
</style>