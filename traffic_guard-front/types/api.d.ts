export interface ApiRequestOptions extends RequestInit {
  params?: Record<string, string | number | boolean>;
  headers?: Record<string, string>;
}

export interface ApiResponse<T = any> {
  data: T;
  payload?: any;
  status: number;
  headers: Headers;
}

export interface defaultResponse {
  message: string;
}

export interface ApiClient {
  get: <T = any>(endpoint: string, options?: ApiRequestOptions) => Promise<T>;
  post: <T = any>(
    endpoint: string,
    data?: any,
    options?: ApiRequestOptions,
  ) => Promise<T>;
  put: <T = any>(
    endpoint: string,
    data?: any,
    options?: ApiRequestOptions,
  ) => Promise<T>;
  patch: <T = any>(
    endpoint: string,
    data?: any,
    options?: ApiRequestOptions,
  ) => Promise<T>;
  delete: <T = any>(
    endpoint: string,
    options?: ApiRequestOptions,
  ) => Promise<T>;
}
