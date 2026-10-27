import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import api from '../services/api'
import type { RetailSale, Promotion, Sparepart } from '../types'

export interface CartItem {
  sparepart: Sparepart
  quantity: number
  discount_percent: number
}

export const useRetailStore = defineStore('retail', () => {
  const sales = ref<RetailSale[]>([])
  const promotions = ref<Promotion[]>([])
  const loading = ref(false)
  const cart = ref<CartItem[]>([])
  const appliedPromo = ref<Promotion | null>(null)

  // Subtotal produk dalam keranjang
  const cartSubtotal = computed(() => {
    return cart.value.reduce((acc, item) => {
      const price = Number(item.sparepart.selling_price)
      const qty = item.quantity
      const discount = item.discount_percent || 0
      const net = (price * qty) * (1 - discount / 100)
      return acc + net
    }, 0)
  })

  // Potongan promo
  const promoDiscountAmount = computed(() => {
    if (!appliedPromo.value) return 0
    const promo = appliedPromo.value
    if (promo.promo_type === 'percentage') {
      return (cartSubtotal.value * Number(promo.discount_value)) / 100
    }
    return Math.min(Number(promo.discount_value), cartSubtotal.value)
  })

  // PPN 11%
  const cartTaxPPN = computed(() => {
    const taxable = Math.max(0, cartSubtotal.value - promoDiscountAmount.value)
    return taxable * 0.11
  })

  // Grand Total
  const cartGrandTotal = computed(() => {
    const taxable = Math.max(0, cartSubtotal.value - promoDiscountAmount.value)
    return taxable + cartTaxPPN.value
  })

  const cartTotalItems = computed(() => {
    return cart.value.reduce((acc, item) => acc + item.quantity, 0)
  })

  function addToCart(part: Sparepart) {
    const existing = cart.value.find(i => i.sparepart.id === part.id)
    if (existing) {
      if (existing.quantity < part.stock) {
        existing.quantity++
      }
    } else {
      if (part.stock > 0) {
        cart.value.push({
          sparepart: part,
          quantity: 1,
          discount_percent: 0,
        })
      }
    }
  }

  function removeFromCart(partId: number) {
    cart.value = cart.value.filter(i => i.sparepart.id !== partId)
  }

  function updateQty(partId: number, qty: number) {
    const item = cart.value.find(i => i.sparepart.id === partId)
    if (item) {
      if (qty <= 0) {
        removeFromCart(partId)
      } else if (qty <= item.sparepart.stock) {
        item.quantity = qty
      }
    }
  }

  function clearCart() {
    cart.value = []
    appliedPromo.value = null
  }

  async function fetchSales(search = '') {
    loading.value = true
    try {
      const res = await api.get('/retail/sales', { params: { search } })
      sales.value = res.data.data
    } finally {
      loading.value = false
    }
  }

  async function fetchPromotions() {
    try {
      const res = await api.get('/retail/promotions')
      promotions.value = res.data.data
    } catch (err) {
      console.error('Failed to load promotions', err)
    }
  }

  function applyPromoCode(code: string): { success: boolean; message: string } {
    const cleanCode = code.trim().toUpperCase()
    const found = promotions.value.find(p => p.promo_code.toUpperCase() === cleanCode)
    if (!found) {
      return { success: false, message: 'Kode promo tidak ditemukan atau sudah kadaluarsa' }
    }
    if (cartSubtotal.value < Number(found.min_spend)) {
      return { success: false, message: `Minimal belanja Rp ${Number(found.min_spend).toLocaleString('id-ID')} untuk promo ini` }
    }
    appliedPromo.value = found
    return { success: true, message: `Promo "${found.title}" berhasil dipasang!` }
  }

  async function checkout(payload: {
    customer_name: string
    customer_phone: string
    payment_method: string
    paid_amount: number
    notes?: string
  }) {
    loading.value = true
    try {
      const itemsPayload = cart.value.map(i => ({
        sparepart_id: i.sparepart.id,
        quantity: i.quantity,
        discount_percent: i.discount_percent,
      }))

      const res = await api.post('/retail/checkout', {
        ...payload,
        promo_code: appliedPromo.value ? appliedPromo.value.promo_code : '',
        items: itemsPayload,
      })

      clearCart()
      await fetchSales()
      return res.data.data as RetailSale
    } finally {
      loading.value = false
    }
  }

  async function fetchReceipt(id: number) {
    const res = await api.get(`/retail/sales/${id}/receipt`)
    return res.data
  }

  return {
    sales,
    promotions,
    loading,
    cart,
    appliedPromo,
    cartSubtotal,
    promoDiscountAmount,
    cartTaxPPN,
    cartGrandTotal,
    cartTotalItems,
    addToCart,
    removeFromCart,
    updateQty,
    clearCart,
    fetchSales,
    fetchPromotions,
    applyPromoCode,
    checkout,
    fetchReceipt,
  }
})
