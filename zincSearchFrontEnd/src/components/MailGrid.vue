<template>
  <div class="flex min-h-full flex-col justify-center px-6 lg:px-8">
    <div class="mt-10 sm:mx-auto sm:w-full sm:max-w-sm">
      <div>
        <label for="search" class="block text-sm font-medium leading-6 text-gray-900">Buscar Emails</label>
        <div class="mt-2">
          <input id="search" name="search" type="text" v-model="state.inputValue"  class="px-2 block w-full rounded-md border-0 py-1.5 text-gray-900 shadow-sm ring-1 ring-inset ring-gray-300 placeholder:text-gray-400 focus:ring-2 focus:ring-inset focus:ring-indigo-600 sm:text-sm sm:leading-6">
        </div>
      </div>
    </div>
  </div>
  <ul role="list" id="table-heading" class="divide-y divide-gray-100 rounded-md p-1 my-3" v-if="state.apiResult.length > 0">
    <MailGridHeader />
  </ul>
  <ul role="list" class="divide-y divide-gray-100" v-if="state.apiResult">
    <MailItem v-for="mail in state.apiResult" :key="mail.ID" :mail="mail" :isHeader="false"/>
  </ul>
</template>

<script lang="ts">
import { watch } from 'vue';
import MailItem from './MailItem.vue';
import MailGridHeader from './MailGridHeader.vue';
import { reactive } from 'vue';
import {type Mail} from '@/types/mailTypes';

export default {
    setup() {
      const state:{
        inputValue: string
        apiResult: Mail[]
      } = reactive({
        inputValue: '',
        apiResult: []
      });

      watch(
        () => state.inputValue,
        async (newQuery) => {
          try {
            const response = await fetch(`http://localhost:9000/search?term=${newQuery}&from=0&max=20`);
            const data = await response.json();
            const formatedData = data.hits.hits.map((mail: { [x: string]: Mail; }) => mail['_source'])
            state.apiResult = formatedData;
          } catch (error) {
            console.error('Error al llamar a la API:', error);
          }
        }
      );

      return {
        state
      };
    },
    components: { MailItem, MailGridHeader }
}
</script>