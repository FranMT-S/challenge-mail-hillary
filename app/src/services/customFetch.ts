import { ApiResponse } from "@/models/response";

/**
 * Custom fetch function to handle API responses
 * transform the response in a json format
 * if exist a error is saved inthe error property
 */
export const customFetch = async <T>(input: RequestInfo | URL, init?: RequestInit): Promise<ApiResponse<T>> => {
  let data: ApiResponse<T>;
  let response: Response;
  
  try {
    response = await fetch(input, init);
    data = await response.json();
    data.status = response.status;
  } catch (error) {
    
    console.error(`Invalid JSON, error: ${error}`);
    data = {
      msg: "Error",
      error: "there was an error, try again later",
      data: undefined,
      status: 500,
    }
  }

  return data;
}