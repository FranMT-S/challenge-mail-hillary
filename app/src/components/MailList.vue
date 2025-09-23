<script setup lang="ts">
import { ref, onMounted } from 'vue';
import Paginator from './Paginator.vue';
import { useMailStore } from '@/store/mailStore';
import { storeToRefs } from 'pinia';
import { useGlobalStore } from '@/store/global';
import { Mail } from '@/models/mails';


const globalStore = useGlobalStore();
const { isDark } = globalStore;
const { mails, total, page, limit, error, isLoading,maxPage} = storeToRefs(useMailStore());
const {fetchMails} = useMailStore();

const emit = defineEmits(['selectMail'])

const onSelectMail = (mail: any) => {
  emit('selectMail', mail)
}

onMounted(async () => {
  await fetchMails();
});

</script>

<template>
<v-data-table-server
  density="comfortable"
  height="100%"
  
  class="tw-h-full table-container"
  :loading="isLoading"
  @update:items-per-page="limit = $event"
  @update:page="page = $event"
  :fixed-header="true"
  :items="mails"
  @click:row="(item: any) => onSelectMail(item)"
  :headers="[
    { title: 'N°', value: 'id' },
    { title: 'To', value: 'to' },
    { title: 'Subject', value: 'subject' },
    { title: 'From', value: 'from' },
    { title: 'Date', value: 'date' }
  ]"
  :items-length="total"

>
  <template v-slot:no-data>
      Can't find any mail
  </template>

  <template v-slot:headers="{ columns, isSorted, getSortIcon, toggleSort }">
        <tr   class="tw-bg-[#f3f4f6] dark:tw-bg-[#181818]  darK:tw-text-[#deffff] ">
          <template v-for="column in columns" :key="column.key">
            <td 
              class="tw-px-2 tw-py-1 tw-font-medium tw-cursor-pointer"
              @click="toggleSort(column)"
              :class="{ 'tw-text-primary': isSorted(column) }"
            >
              {{ column.title }}
              <v-icon 
                v-if="isSorted(column)" 
                :icon="getSortIcon(column)"
                size="small"
                class="ml-1"
              />
            </td>
          </template>
        </tr>
      </template>

  <!-- Filas personalizadas -->
  <template #item="{ item }">
    <tr @click="onSelectMail(item)"  class="row tw-border-b tw-border-gray-700 dark:hover:tw-bg-[#181818] hover:tw-bg-gray-100 ">
      <td class="tw-px-2 tw-py-1">{{ item.id }}</td>
      <td class="tw-px-2 tw-py-1">{{ item.to }}</td>
      <td class="tw-px-2 tw-py-1 font-medium">{{ item.subject || '(No subject)' }}</td>
      <td class="tw-px-2 tw-py-1">{{ item.from }}</td>
      <td class="tw-px-2 tw-py-1">{{ new Date(item.date).toLocaleDateString() }}</td>
    </tr>
  </template>
</v-data-table-server>



</template>

<style>

  .table-container  table {
    height: 85vh;

  }

</style>