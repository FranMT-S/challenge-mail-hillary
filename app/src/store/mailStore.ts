import { DateSearch, Mail } from "@/models/mails";
import { defineStore } from "pinia";
import { computed, nextTick, reactive, ref, watch } from "vue";
import { fetchMails as fetchMailsService } from "@/services/mailService";
import {  DEFAULT_PAGINATION } from "@/constants/api";
import { Operator, OrderBy, TypeSearch } from "@/models/mails";
import { ErrorType } from "@/models/response";

export const useMailStore = defineStore('mail', () => {
    
  const currentMail = ref<Mail | null>(null);
  const query = ref('');
  const mails = ref<Mail[]>([]);
  const total = ref(0);
  const page = ref(DEFAULT_PAGINATION.PAGE);
  const limit = ref(DEFAULT_PAGINATION.LIMIT);
  const type = ref(TypeSearch.AND);
  const orderBy = ref(OrderBy.DESC);
  const dateSearch = reactive<DateSearch>({
    date: undefined,
    operator: Operator.LESS_THAN_OR_EQUAL
  });

  const error = ref<string | null>(null);
  const isLoading = ref(false);
  const cancelationToken = ref<AbortController | null>(null);
  const maxPage = computed(() => Math.ceil(total.value / limit.value));
  

  watch([page, limit],async ([newPage, newLimit], [_, oldLimit]) => {
    if(newPage < 1){
      page.value = DEFAULT_PAGINATION.PAGE;
    } 

    if(newLimit < 1){
      limit.value = DEFAULT_PAGINATION.LIMIT;
    }

    if(newLimit != oldLimit){
      page.value = 1;
    }
    
    await fetchMails();
  });

  watch([type,orderBy,dateSearch],async () => {
    page.value = 1;
    await fetchMails();
  });

  
  const fetchMails = async () => {
    isLoading.value = true;

    if(cancelationToken.value){
      try {
        cancelationToken.value.abort();
      } catch (error) {
      }
    }

    cancelationToken.value = new AbortController();
    const response = await fetchMailsService(
      { 
        query: query.value,
        type: type.value,
        orderBy: orderBy.value,
        page: page.value, 
        limit: limit.value,
        dateSearch: dateSearch
      },
      cancelationToken.value?.signal
    );


    if(response.errorType === ErrorType.ABORT_ERROR){
      return;
    }

    if (response.status >= 400) {
        error.value = ''
        await nextTick()
        error.value = response.error || 'There was an error loading the emails';
        isLoading.value = false;
        cancelationToken.value = null;
        return;
    }

    mails.value = response.data?.mails || [];
    total.value = response.data?.total || 0;
    isLoading.value = false;
    cancelationToken.value = null;
  };

  return {
      mails,
      total,
      page,
      limit,
      error,
      isLoading,
      currentMail,
      maxPage,
      query,
      type,
      orderBy,
      dateSearch,

      fetchMails,
  };
});