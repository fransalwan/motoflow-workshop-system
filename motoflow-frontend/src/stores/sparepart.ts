import { defineStore } from 'pinia'
import { ref } from 'vue'
import api from '../services/api'
import type { Sparepart, Category, APIResponse } from '../types'

export const useSparepartStore = defineStore('sparepart', () => {
  const spareparts = ref<Sparepart[]>([])
  const categories = ref<Category[]>([])
  const loading = ref(false)
  const error = ref<string | null>(null)

  async function fetchSpareparts(filters?: { search?: string; category_id?: string; type?: string; low_stock?: boolean }) {
    loading.value = true
    error.value = null
    try {
      const res = await api.get<APIResponse<Sparepart[]>>('/spareparts', { params: filters })
      spareparts.value = res.data.data
    } catch (err: any) {
      error.value = err.response?.data?.message || 'Gagal memuat suku cadang'
    } finally {
      loading.value = false
    }
  }

  async function fetchCategories() {
    try {
      const res = await api.get<APIResponse<Category[]>>('/spareparts/categories')
      categories.value = res.data.data
    } catch (err) {
      console.error('Failed to load categories', err)
    }
  }

  async function createSparepart(payload: Partial<Sparepart>) {
    const res = await api.post<APIResponse<Sparepart>>('/spareparts', payload)
    return res.data.data
  }

  async function updateStock(id: number, quantity: number, reason: string) {
    const res = await api.patch<APIResponse<Sparepart>>(`/spareparts/${id}/stock`, { quantity, reason })
    const idx = spareparts.value.findIndex(p => p.id === id)
    if (idx !== -1) {
      spareparts.value[idx] = res.data.data
    }
    return res.data.data
  }

  return {
    spareparts,
    categories,
    loading,
    error,
    fetchSpareparts,
    fetchCategories,
    createSparepart,
    updateStock,
  }
})

