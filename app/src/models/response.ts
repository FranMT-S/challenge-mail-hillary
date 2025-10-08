import { Mail } from "./mails";

export enum ErrorType {
  ABORT_ERROR = "AbortError",
  INVALID_JSON = "InvalidJSON",
}

export interface ApiResponse<T> {
  msg: string;
  data?: T | null;
  error?:string | null;
  status: number;
  errorType?: ErrorType;
}

export interface MailResponse {
  mails: Mail[];
  total: number;
}