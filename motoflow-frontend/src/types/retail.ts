import type { Sparepart } from './index'

export interface RetailSaleItem {
  id?: number
  retail_sale_id?: number
  sparepart_id: number
  sparepart?: Sparepart
  part_number_snapshot: string
  part_name_snapshot: string
  quantity: number
  price_snapshot: string | number
  discount_percent: string | number
  subtotal: string | number
}

export interface RetailSale {
  id: number
  invoice_number: string
  customer_id?: number | null
  customer_name: string
  customer_phone?: string
  cashier_id: number
  cashier?: {
    id: number
    full_name: string
    username: string
  }
  subtotal: string | number
  discount_amount: string | number
  tax_ppn: string | number
  grand_total: string | number
  paid_amount: string | number
  change_amount: string | number
  payment_method: 'cash' | 'qris' | 'debit' | 'transfer'
  payment_status: 'paid' | 'cancelled'
  promo_code?: string
  notes?: string
  items: RetailSaleItem[]
  created_at: string
}

export interface Promotion {
  id: number
  promo_code: string
  title: string
  description: string
  promo_type: 'percentage' | 'fixed_nominal' | 'bundle'
  discount_value: string | number
  min_spend: string | number
  is_active: boolean
  valid_until: string
}
