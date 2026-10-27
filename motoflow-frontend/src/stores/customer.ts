import { defineStore } from 'pinia'
import { ref } from 'vue'
import api from '../services/api'
import type { Customer, APIResponse } from '../types'

export const useCustomerStore = defineStore('customer', () => {
  const customers = ref<Customer[]>([])
  const loading = ref(false)
  const error = ref<string | null>(null)

  async function fetchCustomers(search?: string) {
    loading.value = true
    error.value = null
    try {
      const res = await api.get<APIResponse<Customer[]>>('/customers', { params: { search } })
      customers.value = res.data.data
    } catch (err: any) {
      error.value = err.response?.data?.message || 'Gagal memuat pelanggan'
    } finally {
      loading.value = false
    }
  }

  async function createCustomer(payload: Partial<Customer>) {
    const res = await api.post<APIResponse<Customer>>('/customers', payload)
    customers.value.unshift(res.data.data)
    return res.data.data
  }

  return {
    customers,
    loading,
    error,
    fetchCustomers,
    createCustomer,
  }
})

