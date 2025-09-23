import { Mail } from "./mails";

export interface ApiResponse<T> {
  msg: string;
  data?: T | null;
  error?:string | null;
  status: number;
}

export interface MailResponse {
  mails: Mail[];
  total: number;
}