// serach enums
export enum Operator{
  GREATER_THAN = ">",
  GREATER_THAN_OR_EQUAL = ">=",
  LESS_THAN = "<",
  LESS_THAN_OR_EQUAL = "<=",
}

export enum OrderBy{
  ASC = "asc",
  DESC = "desc",
}

export enum TypeSearch{
  AND = "AND",
  OR = "OR",
}


export interface Mail {
  id: string;
  from: string;
  to: string;
  subject: string;
  content: string;
  date: string;
}

export interface DateSearch{
  date?: Date | null;
  operator: Operator;
}

export interface MailSearch{
  query: string;
  type: TypeSearch;
  page: number;
  limit: number;
  orderBy: OrderBy;
  dateSearch?: DateSearch | null;
}
