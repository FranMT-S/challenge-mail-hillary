<script setup lang="ts">
import { useToast } from '../composable/useToast'
const { state, close } = useToast()

</script>

<template>
  <v-snackbar
    v-model="state.show"
    :color="state.color"
    :timeout="state.timeout"
    location="top right"

  >
    {{ state.message }}
    <template #actions>
      <v-btn color="white" text @click="close">Cerrar</v-btn>
    </template>
  </v-snackbar>
</template>
