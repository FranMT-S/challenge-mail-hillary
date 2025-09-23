import { API_CONFIG, DEFAULT_PAGINATION } from '@/constants/api';

import { PaginationParams } from '@/models/pagination';
import { customFetch } from './customFetch';
import { ApiResponse, MailResponse } from '@/models/response';
import { MailSearch, Operator, OrderBy, TypeSearch } from '@/models/mails';
import DOMPurify from 'dompurify';

export const fetchMails = async (search: MailSearch, cancelationToken: AbortSignal): Promise<ApiResponse<MailResponse>> => {
  const url = new URL(`${API_CONFIG.ENDPOINTS.MAILS}`, API_CONFIG.BASE_URL);

  search.query = DOMPurify.sanitize(search.query.trim());

  let defaultBody: MailSearch = {
    query: "",
    type: TypeSearch.OR,
    page: DEFAULT_PAGINATION.PAGE,
    limit: DEFAULT_PAGINATION.LIMIT,
    dateSearch: {
      date: null ,
      operator: Operator.LESS_THAN_OR_EQUAL
    },
    orderBy: OrderBy.ASC,  
  }

  // merge default body with search
  let body = {...defaultBody,...search, }

  const response = await customFetch<MailResponse>(url.toString(),
  {
    method: 'POST',
    body: JSON.stringify(body),
    signal: cancelationToken
  });
  
  return response;
};
