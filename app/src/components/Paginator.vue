<script setup lang="ts">
import { ref, watch } from "vue"

const props = defineProps<{
  maxPage: number
  page: number
  itemsPerPage: number
}>()

const emit = defineEmits<{
  (e: "update:page", value: number): void
  (e: "update:itemsPerPage", value: number): void
}>()

const localPage = ref(props.page)
const localItemsPerPage = ref(props.itemsPerPage)

// actualizar cuando cambian las props desde afuera
watch(
  () => props.page,
  (val) => (localPage.value = val)
)
watch(
  () => props.itemsPerPage,
  (val) => (localItemsPerPage.value = val)
)


watch(localPage, (val) => emit("update:page", val))
watch(localItemsPerPage, (val) => {
  emit("update:itemsPerPage", val)
  emit("update:page", 1)
})
</script>

<template>
  <div class="d-flex align-center justify-space-between pa-2">

    
    <v-pagination
      v-model="localPage"
      :length="maxPage"
      total-visible="7"
    />


    <v-select
      v-model="localItemsPerPage"
      :items="[10,20,50,100,200]"
      label="Items por página"
      density="compact"
      style="max-width: 150px"
    />

  </div>
</template>
