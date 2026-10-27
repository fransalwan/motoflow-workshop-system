export function formatRupiah(amount: number | string | undefined | null): string {
  if (amount === undefined || amount === null || amount === '') return 'Rp 0'
  const num = typeof amount === 'string' ? parseFloat(amount) : amount
  if (isNaN(num)) return 'Rp 0'
  return new Intl.NumberFormat('id-ID', {
    style: 'currency',
    currency: 'IDR',
    minimumFractionDigits: 0,
    maximumFractionDigits: 0,
  }).format(num)
}

export function formatDate(dateStr: string | undefined | null): string {
  if (!dateStr) return '-'
  const d = new Date(dateStr)
  return new Intl.DateTimeFormat('id-ID', {
    day: '2-digit',
    month: 'short',
    year: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
  }).format(d) + ' WIB'
}

export function getStatusBadge(status: string): { label: string; bgClass: string; textClass: string; borderClass: string } {
  switch (status) {
    case 'draft':
      return { label: 'Draft / Masuk', bgClass: 'bg-amber-50', textClass: 'text-amber-700', borderClass: 'border-amber-200' }
    case 'in_progress':
      return { label: 'Dikerjakan (Pit)', bgClass: 'bg-blue-50', textClass: 'text-blue-700', borderClass: 'border-blue-200' }
    case 'completed':
      return { label: 'Selesai Servis', bgClass: 'bg-purple-50', textClass: 'text-purple-700', borderClass: 'border-purple-200' }
    case 'paid':
      return { label: 'Lunas Kasir', bgClass: 'bg-emerald-50', textClass: 'text-emerald-700', borderClass: 'border-emerald-200' }
    case 'cancelled':
      return { label: 'Dibatalkan', bgClass: 'bg-slate-100', textClass: 'text-slate-600', borderClass: 'border-slate-200' }
    default:
      return { label: status, bgClass: 'bg-slate-50', textClass: 'text-slate-600', borderClass: 'border-slate-200' }
  }
}

