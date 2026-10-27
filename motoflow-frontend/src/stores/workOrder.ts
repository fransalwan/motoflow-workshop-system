import { defineStore } from 'pinia'
import { ref } from 'vue'
import api from '../services/api'
import type { WorkOrder, APIResponse, PaymentMethod, WorkOrderStatus } from '../types'

export const useWorkOrderStore = defineStore('workOrder', () => {
  const workOrders = ref<WorkOrder[]>([])
  const currentWorkOrder = ref<WorkOrder | null>(null)
  const loading = ref(false)
  const error = ref<string | null>(null)

  async function fetchWorkOrders(filters?: { status?: string; search?: string; payment_status?: string }) {
    loading.value = true
    error.value = null
    try {
      const res = await api.get<APIResponse<WorkOrder[]>>('/work-orders', { params: filters })
      workOrders.value = res.data.data
    } catch (err: any) {
      error.value = err.response?.data?.message || 'Gagal memuat data Work Order'
    } finally {
      loading.value = false
    }
  }

  async function fetchWorkOrder(id: number) {
    loading.value = true
    error.value = null
    try {
      const res = await api.get<APIResponse<WorkOrder>>(`/work-orders/${id}`)
      currentWorkOrder.value = res.data.data
      return res.data.data
    } catch (err: any) {
      error.value = err.response?.data?.message || 'Gagal memuat detail Work Order'
      return null
    } finally {
      loading.value = false
    }
  }

  async function createWorkOrder(payload: any) {
    loading.value = true
    error.value = null
    try {
      const res = await api.post<APIResponse<WorkOrder>>('/work-orders', payload)
      return res.data.data
    } catch (err: any) {
      error.value = err.response?.data?.message || 'Gagal membuat Work Order'
      throw new Error(error.value || 'Gagal membuat Work Order')
    } finally {
      loading.value = false
    }
  }

  async function updateStatus(id: number, status: WorkOrderStatus, notes?: string) {
    loading.value = true
    try {
      const res = await api.patch<APIResponse<WorkOrder>>(`/work-orders/${id}/status`, { status, notes })
      if (currentWorkOrder.value && currentWorkOrder.value.id === id) {
        currentWorkOrder.value = res.data.data
      }
      return res.data.data
    } catch (err: any) {
      throw new Error(err.response?.data?.message || 'Gagal mengubah status Work Order')
    } finally {
      loading.value = false
    }
  }

  async function processPayment(id: number, paymentMethod: PaymentMethod, paidAmount: number | string, notes?: string) {
    loading.value = true
    try {
      const res = await api.post<APIResponse<WorkOrder>>(`/work-orders/${id}/pay`, {
        payment_method: paymentMethod,
        paid_amount: paidAmount,
        notes,
      })
      if (currentWorkOrder.value && currentWorkOrder.value.id === id) {
        currentWorkOrder.value = res.data.data
      }
      return res.data.data
    } catch (err: any) {
      throw new Error(err.response?.data?.message || 'Gagal memproses pembayaran kasir')
    } finally {
      loading.value = false
    }
  }

  async function getThermalReceipt(id: number, width: number = 58) {
    const res = await api.get<APIResponse<any>>(`/work-orders/${id}/receipt`, { params: { width } })
    return res.data.data
  }

  return {
    workOrders,
    currentWorkOrder,
    loading,
    error,
    fetchWorkOrders,
    fetchWorkOrder,
    createWorkOrder,
    updateStatus,
    processPayment,
    getThermalReceipt,
  }
})

