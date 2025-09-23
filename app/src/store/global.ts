import { defineStore } from "pinia";
import { computed, ref, watch } from "vue";
import { useTheme } from "vuetify";

export const useGlobalStore = defineStore('global', () => {

    const Theme = useTheme();
    Theme.change(localStorage.getItem('theme') || 'light');
    
    const updateTheme = () =>{
      try {
        if (Theme.global.name.value === 'dark') {
          document.body.classList.add('tw-dark')
        } else {
          document.body.classList.remove('tw-dark')
        }

        localStorage.setItem('theme', Theme.global.name.value);
      } catch (error) {
        console.error('error saving theme:', error);
      }
    }

    updateTheme();

    watch(() => Theme.global.name.value, (value) => {
      updateTheme();
    });

    const toggleTheme = () => {
      Theme.change(Theme.global.name.value === 'light' ? 'dark' : 'light');
      updateTheme();
    };

  
    const isDark = computed(() => Theme.global.name.value === 'dark');
    return {
      toggleTheme,
      isDark,
    };
  });