<script setup lang="ts">
import { Mail, Operator } from '@/models/mails';
import { useMailStore } from '@/store/mailStore';
import { storeToRefs } from 'pinia';
import { watch } from 'vue';
import { VDateInput } from 'vuetify/labs/VDateInput';
import MailList from './components/MailList.vue';
import MailModal from './components/MailModal.vue';
import ToastNotification from './components/toast/components/ToastNotification.vue';
import { useToast } from './components/toast/composable/useToast';
import ToggleTheme from './components/ToggleTheme.vue';

const { open } = useToast()

const mailStore = useMailStore();
const {currentMail,error} = storeToRefs(mailStore);


const selectMail = (mail: Mail) => {
  currentMail.value = mail;
};



watch(error, () => {
  if(error.value != null && error.value != ''){
    open(error.value, 'red-darken-2')
  }
})


</script>

<template>
  <ToastNotification />
  <div class="tw-container-fluid tw-mx-auto tw-p-0 tw-h-screen ">
    <v-app class="app-container">
      <v-header class="tw-px-4 tw-py-2 tw-items-center tw-flex tw-justify-between app-header tw-border-b tw-border-gray-300">
        <h1 class="tw-text-xl tw-font-bold  ">Hillary Clinton Emails</h1>
        
        <v-text-field density="compact"  
        class="mx-8"
        variant="underlined"
        v-model="mailStore.query" 
        placeholder="Search"
        @keyup.enter="mailStore.fetchMails()"
        hide-details clearable prepend-inner-icon="mdi-magnify"/>
        <ToggleTheme />
      </v-header>

      <v-main class="app-main">
        <v-container>
          <v-row>
            <v-col cols="auto" md="auto">
              <v-select
                density="compact"
                variant="underlined"
                class="max-w-con capitalize"
                v-model="mailStore.dateSearch.operator"
                :items="[
                  { title: 'Less Than', value: Operator.LESS_THAN_OR_EQUAL },
                  { title: 'Greater Than', value: Operator.GREATER_THAN_OR_EQUAL }
                ]"
                @change="mailStore.fetchMails()"
              ></v-select>
            </v-col>
            <v-col cols="auto" md="auto">
              <VDateInput
                density="compact"
                v-model="mailStore.dateSearch.date"
                label="Select a date"
                min-width="200"
              ></VDateInput>
            </v-col>
          </v-row>
        </v-container>

        <v-container class="tw-px-4 py-0">
          <template v-if="currentMail">
            <MailModal :mail="currentMail" @onClose="currentMail = null" />
          </template>
          <MailList @selectMail="selectMail" class="tw-overflow-y-auto py-0" />
        </v-container>
      </v-main>
    </v-app>
  </div>
</template>

<style scoped>
/* Custom styles can be added here */
.app-container{
  display: grid;
  grid-template-areas: "header "
                       "main ";
  grid-template-rows: 10vh 90vh  ;
  grid-template-columns:  1fr;

}

.app-header{
  grid-area: header;
}

.app-main{
  grid-area: main;
  display: grid;
  grid-template-rows: 10vh 75vh;
}

</style>
