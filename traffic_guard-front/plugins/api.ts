import { defineNuxtPlugin, useRuntimeConfig } from '#app'
import type { ApiClient, ApiRequestOptions, ApiResponse } from '~/types/api'


export default defineNuxtPlugin((nuxtApp) => {
  const config = useRuntimeConfig();
  const apiBaseUrl = `${config.public.apiGateway.host}:${config.public.apiGateway.port}/api/v1`

  const getFullUrl = (endpoint: string, params?: Record<string, string | number | boolean>) => {
    const url = `${apiBaseUrl}${endpoint.startsWith('/') ? endpoint : `/${endpoint}`}`
    
    if (!params) return url

    const queryParams = new URLSearchParams()
    Object.entries(params).forEach(([key, value]) => {
      queryParams.append(key, String(value))
    })

    return `${url}${url.includes('?') ? '&' : '?'}${queryParams.toString()}`
  }

  const handleResponse = async <T>(response: Response): Promise<ApiResponse<T>> => {
    if (!response.ok) {
      const errorData = await response.json().catch(() => ({}))
      throw new Error(`${response.status}` || errorData.error)
    }

    const contentType = response.headers.get('content-type')
    let data: T

    if (contentType && contentType.includes('application/json')) {
      data = await response.json()
    } else {
      data = await response.text() as unknown as T
    }

    return {
      data,
      status: response.status,
      headers: response.headers
    }
  }

  const api: ApiClient = {
    get: async <T>(endpoint: string, options: ApiRequestOptions = {}) => {
      const { params, ...fetchOptions } = options
      const response = await fetch(getFullUrl(endpoint, params), {
        ...fetchOptions,
        method: 'GET',
        // credentials: 'include',
        headers: {
          'Accept': 'application/json',
          ...fetchOptions.headers
        }
      })
      return handleResponse<T>(response).then(res => res.data)
    },

    post: async <T>(endpoint: string, data?: any, options: ApiRequestOptions = {}) => {
      const { params, ...fetchOptions } = options
      const response = await fetch(getFullUrl(endpoint, params), {
        ...fetchOptions,
        method: 'POST',
        // credentials: 'include',
        headers: {
          'Accept': 'application/json',
          'Content-Type': 'application/json',
          ...fetchOptions.headers
        },
        body: data ? JSON.stringify(data) : undefined
      })
      const result = await handleResponse<T>(response).then(res => res.data)
      return result
    },

    put: async <T>(endpoint: string, data?: any, options: ApiRequestOptions = {}) => {
      const { params, ...fetchOptions } = options
      const response = await fetch(getFullUrl(endpoint, params), {
        ...fetchOptions,
        method: 'PUT',
        // credentials: 'include',
        headers: {
          'Accept': 'application/json',
          'Content-Type': 'application/json',
          ...fetchOptions.headers
        },
        body: data ? JSON.stringify(data) : undefined
      })
      return handleResponse<T>(response).then(res => res.data)
    },

    patch: async <T>(endpoint: string, data?: any, options: ApiRequestOptions = {}) => {
      const { params, ...fetchOptions } = options
      const response = await fetch(getFullUrl(endpoint, params), {
        ...fetchOptions,
        method: 'PATCH',
        // credentials: 'include',
        headers: {
          'Accept': 'application/json',
          'Content-Type': 'application/json',
          ...fetchOptions.headers
        },
        body: data ? JSON.stringify(data) : undefined
      })
      return handleResponse<T>(response).then(res => res.data)
    },


    delete: async <T>(endpoint: string, options: ApiRequestOptions = {}) => {
      const { params, ...fetchOptions } = options
      const response = await fetch(getFullUrl(endpoint, params), {
        ...fetchOptions,
        method: 'DELETE',
        // credentials: 'include',
        headers: {
          'Accept': 'application/json',
          ...fetchOptions.headers
        }
      })
      return handleResponse<T>(response).then(res => res.data)
    },
  }

  nuxtApp.provide("api", api)
})
