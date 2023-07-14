<script setup lang="ts">
import MailGrid from '../components/MailGrid.vue'
import HeaderComponent from '@/components/HeaderComponent.vue';
</script>

<template>
  <HeaderComponent :msg="'Desafío Zincsearch-Vue'"/>
  <main>
    <MailGrid />
  </main>
</template>
