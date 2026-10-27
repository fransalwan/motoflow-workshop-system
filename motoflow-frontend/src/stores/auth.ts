import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import api from '../services/api'
import type { User, APIResponse } from '../types'

export const useAuthStore = defineStore('auth', () => {
  const token = ref<string | null>(localStorage.getItem('motoflow_token'))
  const user = ref<User | null>(
    localStorage.getItem('motoflow_user') ? JSON.parse(localStorage.getItem('motoflow_user')!) : null
  )
  const loading = ref(false)
  const error = ref<string | null>(null)

  const isAuthenticated = computed(() => !!token.value)
  const userRole = computed(() => user.value?.role || '')

  async function login(username: string, password: string): Promise<boolean> {
    loading.value = true
    error.value = null
    try {
      const res = await api.post<APIResponse<{ token: string; user: User }>>('/auth/login', {
        username,
        password,
      })
      if (res.data.success) {
        token.value = res.data.data.token
        user.value = res.data.data.user
        localStorage.setItem('motoflow_token', res.data.data.token)
        localStorage.setItem('motoflow_user', JSON.stringify(res.data.data.user))
        return true
      }
      return false
    } catch (err: any) {
      error.value = err.response?.data?.message || 'Login gagal. Periksa username dan kata sandi.'
      return false
    } finally {
      loading.value = false
    }
  }

  function logout() {
    token.value = null
    user.value = null
    localStorage.removeItem('motoflow_token')
    localStorage.removeItem('motoflow_user')
  }

  return {
    token,
    user,
    loading,
    error,
    isAuthenticated,
    userRole,
    login,
    logout,
  }
})

