import { defineStore } from 'pinia'
import { ref } from 'vue'
import api from '../services/api'
import type { DashboardSummary, APIResponse } from '../types'

export const useDashboardStore = defineStore('dashboard', () => {
  const summary = ref<DashboardSummary | null>(null)
  const loading = ref(false)
  const error = ref<string | null>(null)

  async function fetchSummary() {
    loading.value = true
    error.value = null
    try {
      const res = await api.get<APIResponse<DashboardSummary>>('/dashboard/summary')
      summary.value = res.data.data
    } catch (err: any) {
      error.value = err.response?.data?.message || 'Gagal memuat ringkasan dashboard'
    } finally {
      loading.value = false
    }
  }

  return {
    summary,
    loading,
    error,
    fetchSummary,
  }
})

