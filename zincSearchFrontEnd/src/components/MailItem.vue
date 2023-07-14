<template>
  <li class="flex justify-between gap-x-6 py-5 px-3 grow">
    <div class="flex gap-x-4">
      <div class="min-w-0 flex-auto">
        <p class="text-sm font-semibold leading-6 text-gray-900 content">{{ mail['X-From'] }}</p>
        <p class="mt-1 truncate text-xs leading-5 text-gray-500 content">{{ mail.From }}</p>
      </div>
    </div>
    <div class="hidden sm:flex sm:flex-col sm:items-end">
      <p class="text-sm leading-6 text-gray-900 content">{{ mail.Subject }}</p>
      <p class="mt-1 text-xs leading-5 text-gray-500 content"><time datetime="2023-01-23T13:23Z">{{ mail.Date }}</time></p>
    </div>
    <div>
      <button class="rounded bg-indigo-500 px-2 text-slate-50 me-2">
        Ver detalles
      </button>
    </div>
  </li>
</template>

<script setup lang="ts">
import { type Mail } from '@/types/mailTypes';

defineProps<{
  mail: Mail
}>()
</script>

<style>
.content {
  width: 200px;
  white-space: nowrap; 
  overflow: hidden; 
  text-overflow: ellipsis;
}
</style>