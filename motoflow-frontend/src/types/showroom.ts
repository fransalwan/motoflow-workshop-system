export interface MotorUnit {
  id: number
  model_name: string
  category: string
  frame_number: string
  engine_number: string
  color: string
  year: number
  otr_price: string | number
  min_dp: string | number
  status: 'ready' | 'booked' | 'delivered'
  image_url: string
  specs_summary: string
  warehouse_location: string
  created_at?: string
}

export interface CreditSimulationResult {
  otr_price: string | number
  dp_percent: string | number
  dp_amount: string | number
  principal_loan: string | number
  tenor_months: number
  annual_interest_rate: string | number
  total_interest: string | number
  total_loan_payable: string | number
  monthly_installment: string | number
  admin_fee: string | number
  insurance_tlo: string | number
  total_first_payment: string | number
}

export interface SPKUnit {
  id: number
  spk_number: string
  customer_id: number
  customer?: {
    id: number
    name: string
    phone: string
    address?: string
  }
  sales_id: number
  sales?: {
    id: number
    full_name: string
    username: string
  }
  motor_unit_id: number
  motor_unit?: MotorUnit
  purchase_type: 'cash' | 'credit_nsc'
  booking_fee: string | number
  down_payment: string | number
  tenor_months: number
  monthly_installment: string | number
  leasing_status: string
  status: string
  notes?: string
  created_at: string
}

