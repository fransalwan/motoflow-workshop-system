import { defineStore } from 'pinia'
import { ref } from 'vue'
import api from '../services/api'
import type { MotorUnit, CreditSimulationResult, SPKUnit } from '../types'

export const useShowroomStore = defineStore('showroom', () => {
  const units = ref<MotorUnit[]>([])
  const spkList = ref<SPKUnit[]>([])
  const loading = ref(false)
  const simulationResult = ref<CreditSimulationResult | null>(null)

  async function fetchUnits(category = '', status = '', search = '') {
    loading.value = true
    try {
      const res = await api.get('/showroom/units', {
        params: { category, status, search },
      })
      units.value = res.data.data
    } finally {
      loading.value = false
    }
  }

  async function simulateCredit(otrPrice: number, dpPercent: number, tenorMonths: number) {
    try {
      const res = await api.post('/showroom/simulate-credit', {
        otr_price: otrPrice,
        dp_percent: dpPercent,
        tenor_months: tenorMonths,
      })
      simulationResult.value = res.data.data
      return res.data.data as CreditSimulationResult
    } catch (err) {
      console.error('Failed to simulate credit', err)
      throw err
    }
  }

  async function createSPK(payload: {
    customer_id: number
    motor_unit_id: number
    purchase_type: 'cash' | 'credit_nsc'
    booking_fee: number
    down_payment: number
    tenor_months: number
    monthly_installment: number
    notes?: string
  }) {
    loading.value = true
    try {
      const res = await api.post('/showroom/spk', payload)
      await fetchUnits()
      await fetchSPKList()
      return res.data.data as SPKUnit
    } finally {
      loading.value = false
    }
  }

  async function fetchSPKList(search = '') {
    loading.value = true
    try {
      const res = await api.get('/showroom/spk', { params: { search } })
      spkList.value = res.data.data
    } finally {
      loading.value = false
    }
  }

  async function fetchSPKReceipt(id: number) {
    const res = await api.get(`/showroom/spk/${id}/receipt`)
    return res.data
  }

  return {
    units,
    spkList,
    loading,
    simulationResult,
    fetchUnits,
    simulateCredit,
    createSPK,
    fetchSPKList,
    fetchSPKReceipt,
  }
})
