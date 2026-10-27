export type Role = 'admin' | 'service_advisor' | 'mechanic' | 'cashier'

export interface User {
  id: number
  username: string
  email: string
  full_name: string
  role: Role
  phone: string
  is_active: boolean
  created_at: string
}

export interface Customer {
  id: number
  name: string
  phone: string
  email?: string
  address?: string
  plate_number: string
  motor_model: string
  year?: number
  frame_number?: string
  engine_number?: string
  current_km: number
  work_orders?: WorkOrder[]
  created_at: string
}

export interface Category {
  id: number
  name: string
  code: string
  description?: string
}

export type ItemType = 'part' | 'service'

export interface Sparepart {
  id: number
  part_number: string
  name: string
  category_id: number
  category?: Category
  type: ItemType
  cost_price: string | number
  selling_price: string | number
  stock: number
  min_stock: number
  unit: string
  location?: string
  is_active: boolean
  created_at: string
}

export interface Mechanic {
  id: number
  code: string
  name: string
  phone?: string
  specialization: string
  pit_number: number
  is_available: boolean
}

export type WorkOrderStatus = 'draft' | 'in_progress' | 'completed' | 'paid' | 'cancelled'
export type PaymentMethod = 'cash' | 'transfer' | 'qris' | 'debit'
export type PaymentStatus = 'unpaid' | 'paid'

export interface WorkOrderItem {
  id?: number
  work_order_id?: number
  sparepart_id?: number
  item_type: ItemType
  part_number_snapshot: string
  item_name_snapshot: string
  quantity: number
  price_snapshot: string | number
  discount_percent: string | number
  subtotal: string | number
  mechanic_notes?: string
}

export interface WorkOrder {
  id: number
  wo_number: string
  customer_id: number
  customer?: Customer
  mechanic_id?: number
  mechanic?: Mechanic
  service_advisor_id?: number
  service_advisor?: User
  entry_km: number
  complaints: string
  service_diagnosis: string
  status: WorkOrderStatus
  payment_method: PaymentMethod
  payment_status: PaymentStatus
  parts_subtotal: string | number
  services_subtotal: string | number
  total_discount: string | number
  tax_rate: string | number
  tax_amount: string | number
  grand_total: string | number
  paid_amount: string | number
  change_amount: string | number
  completed_at?: string
  paid_at?: string
  notes?: string
  items: WorkOrderItem[]
  created_at: string
  updated_at: string
}

export interface TopPartStat {
  part_number: string
  part_name: string
  total_sold_qty: number
  total_revenue: string | number
}

export interface MechanicStat {
  mechanic_id: number
  code: string
  name: string
  pit_number: number
  completed_jobs: number
  labor_revenue: string | number
}

export interface DailyRevenueStat {
  date: string
  total_revenue: string | number
  bike_count: number
}

export interface DashboardSummary {
  today_revenue: string | number
  today_bikes_serviced: number
  month_revenue: string | number
  month_bikes_serviced: number
  active_jobs_count: number
  low_stock_count: number
  status_counts: Record<string, number>
  top_selling_parts: TopPartStat[]
  mechanic_rankings: MechanicStat[]
  daily_trend: DailyRevenueStat[]
}

export interface APIResponse<T = any> {
  success: boolean
  message: string
  data: T
  error?: string
}

