<script setup lang="ts">

import { ref, watch } from 'vue';
import DOMPurify from 'dompurify';
import { Mail } from '@/models/mails';

const emit = defineEmits(['onClose'])
const props = defineProps<{mail: Mail}>()
const content = ref(DOMPurify.sanitize(props.mail.content));

watch(() => props.mail, () => {
  content.value = DOMPurify.sanitize(props.mail.content);
})



</script>

<template>
  <v-dialog max-width="800" :model-value="props.mail != null" @update:model-value="$emit('onClose', $event)">
    <template v-slot:default="{ isActive }">
      <v-card> 
        <v-toolbar color="primary">
          <v-toolbar-title :title="props.mail.subject">{{ props.mail.subject }}</v-toolbar-title>
          <v-btn
            icon="mdi-close"
            @click="$emit('onClose', false)"
          ></v-btn>
        </v-toolbar>

        <v-list>
          <v-list-item>
            <v-card-text class="tw-whitespace-pre-line">
              <div v-html="content"></div>
            </v-card-text>
          </v-list-item>
        </v-list>
      </v-card>
    </template>
  </v-dialog>
</template>

<style scoped>
  :deep(.v-card-text h2) {
    font-size: 1.2rem;
    margin-bottom: 1rem;
    font-weight: bold !important;
  }
</style>