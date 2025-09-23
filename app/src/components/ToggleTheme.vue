<script setup lang="ts">
import { useGlobalStore } from "../store/global";

const theme = useGlobalStore();

function toggleTheme() {
  theme.toggleTheme();
}
</script>

<template>
  <v-btn :icon="theme.isDark ? 'mdi-weather-sunny' : 'mdi-weather-night'" @click="toggleTheme" variant="text" size="small" >
    
  </v-btn>
</template>
  