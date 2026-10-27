<template>
  <div class="space-y-8">
    <div v-if="dashboardStore.loading && !dashboardStore.summary">
      <SkeletonLoader />
    </div>

    <div v-else-if="dashboardStore.summary" class="space-y-8">
      <!-- Top Metrics Cards -->
      <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-6">
        <!-- Metric 1: Omzet Hari Ini -->
        <div class="bg-white p-6 rounded-2xl border border-slate-200/80 shadow-xs relative overflow-hidden group hover:border-red-300 transition">
          <div class="flex items-center justify-between">
            <span class="text-xs font-bold text-slate-500 uppercase tracking-wider">Omzet Kasir Hari Ini</span>
            <div class="w-10 h-10 rounded-xl bg-red-50 text-red-600 flex items-center justify-center">
              <DollarSign class="w-5 h-5" />
            </div>
          </div>
          <div class="mt-4">
            <h3 class="text-2xl font-black text-slate-900 tracking-tight">
              {{ formatRupiah(dashboardStore.summary.today_revenue) }}
            </h3>
            <p class="text-xs text-slate-500 mt-1 flex items-center space-x-1">
              <Bike class="w-3.5 h-3.5 text-slate-400" />
              <span>{{ dashboardStore.summary.today_bikes_serviced }} unit motor masuk hari ini</span>
            </p>
          </div>
        </div>

        <!-- Metric 2: Omzet Bulan Ini (MTD) -->
        <div class="bg-white p-6 rounded-2xl border border-slate-200/80 shadow-xs relative overflow-hidden group hover:border-emerald-300 transition">
          <div class="flex items-center justify-between">
            <span class="text-xs font-bold text-slate-500 uppercase tracking-wider">Omzet Bulan Ini (MTD)</span>
            <div class="w-10 h-10 rounded-xl bg-emerald-50 text-emerald-600 flex items-center justify-center">
              <TrendingUp class="w-5 h-5" />
            </div>
          </div>
          <div class="mt-4">
            <h3 class="text-2xl font-black text-slate-900 tracking-tight">
              {{ formatRupiah(dashboardStore.summary.month_revenue) }}
            </h3>
            <p class="text-xs text-slate-500 mt-1 flex items-center space-x-1">
              <Calendar class="w-3.5 h-3.5 text-slate-400" />
              <span>{{ dashboardStore.summary.month_bikes_serviced }} total unit servis MTD</span>
            </p>
          </div>
        </div>

        <!-- Metric 3: Motor Aktif di Pit -->
        <div class="bg-white p-6 rounded-2xl border border-slate-200/80 shadow-xs relative overflow-hidden group hover:border-blue-300 transition">
          <div class="flex items-center justify-between">
            <span class="text-xs font-bold text-slate-500 uppercase tracking-wider">Motor Dalam Pengerjaan</span>
            <div class="w-10 h-10 rounded-xl bg-blue-50 text-blue-600 flex items-center justify-center">
              <Wrench class="w-5 h-5" />
            </div>
          </div>
          <div class="mt-4">
            <h3 class="text-2xl font-black text-slate-900 tracking-tight">
              {{ dashboardStore.summary.active_jobs_count }} Unit
            </h3>
            <p class="text-xs text-blue-600 mt-1 font-medium">Sedang ditangani di Pit Servis</p>
          </div>
        </div>

        <!-- Metric 4: Peringatan Stok HGP Menipis -->
        <div class="bg-white p-6 rounded-2xl border border-slate-200/80 shadow-xs relative overflow-hidden group hover:border-amber-300 transition">
          <div class="flex items-center justify-between">
            <span class="text-xs font-bold text-slate-500 uppercase tracking-wider">Peringatan Stok Menipis</span>
            <div class="w-10 h-10 rounded-xl bg-amber-50 text-amber-600 flex items-center justify-center">
              <AlertTriangle class="w-5 h-5" />
            </div>
          </div>
          <div class="mt-4">
            <h3 class="text-2xl font-black text-slate-900 tracking-tight" :class="dashboardStore.summary.low_stock_count > 0 ? 'text-amber-600' : 'text-slate-900'">
              {{ dashboardStore.summary.low_stock_count }} Part
            </h3>
            <router-link to="/inventory" class="text-xs text-slate-500 hover:text-red-600 mt-1 inline-block font-medium underline">
              Cek gudang suku cadang &rarr;
            </router-link>
          </div>
        </div>
      </div>

      <!-- Work Order Status Overview Banner -->
      <div class="bg-white p-6 rounded-2xl border border-slate-200/80 shadow-xs">
        <h4 class="text-xs font-bold text-slate-500 uppercase tracking-wider mb-4">Status Alur Kerja Servis AHASS Real-time</h4>
        <div class="grid grid-cols-2 sm:grid-cols-4 gap-4">
          <div class="p-4 rounded-xl bg-amber-50/60 border border-amber-200/60">
            <p class="text-xs font-bold text-amber-800 uppercase">Draft / Pendaftaran</p>
            <p class="text-2xl font-black text-amber-900 mt-1">{{ dashboardStore.summary.status_counts['draft'] || 0 }}</p>
          </div>
          <div class="p-4 rounded-xl bg-blue-50/60 border border-blue-200/60">
            <p class="text-xs font-bold text-blue-800 uppercase">In Progress (Pit)</p>
            <p class="text-2xl font-black text-blue-900 mt-1">{{ dashboardStore.summary.status_counts['in_progress'] || 0 }}</p>
          </div>
          <div class="p-4 rounded-xl bg-purple-50/60 border border-purple-200/60">
            <p class="text-xs font-bold text-purple-800 uppercase">Selesai Inspeksi</p>
            <p class="text-2xl font-black text-purple-900 mt-1">{{ dashboardStore.summary.status_counts['completed'] || 0 }}</p>
          </div>
          <div class="p-4 rounded-xl bg-emerald-50/60 border border-emerald-200/60">
            <p class="text-xs font-bold text-emerald-800 uppercase">Lunas Kasir (Paid)</p>
            <p class="text-2xl font-black text-emerald-900 mt-1">{{ dashboardStore.summary.status_counts['paid'] || 0 }}</p>
          </div>
        </div>
      </div>

      <!-- Two Columns: Top 5 HGP Terlaris & Mechanic Ranking -->
      <div class="grid grid-cols-1 lg:grid-cols-2 gap-8">
        <!-- Top 5 Honda Genuine Parts -->
        <div class="bg-white p-6 rounded-2xl border border-slate-200/80 shadow-xs">
          <div class="flex items-center justify-between mb-4">
            <div>
              <h3 class="font-bold text-slate-800 text-base">Top 5 Honda Genuine Parts Terlaris</h3>
              <p class="text-xs text-slate-500">Berdasarkan total omset penjualan part</p>
            </div>
            <Package class="w-5 h-5 text-red-500" />
          </div>

          <div class="overflow-x-auto">
            <table class="w-full text-xs text-left">
              <thead>
                <tr class="border-b border-slate-100 text-slate-400 font-semibold uppercase text-[10px]">
                  <th class="py-2.5">Part Number</th>
                  <th class="py-2.5">Nama Suku Cadang</th>
                  <th class="py-2.5 text-center">Terjual</th>
                  <th class="py-2.5 text-right">Total Omzet</th>
                </tr>
              </thead>
              <tbody class="divide-y divide-slate-100">
                <tr v-for="part in dashboardStore.summary.top_selling_parts" :key="part.part_number" class="hover:bg-slate-50/60">
                  <td class="py-3 font-mono font-medium text-slate-600">{{ part.part_number }}</td>
                  <td class="py-3 font-semibold text-slate-800">{{ part.part_name }}</td>
                  <td class="py-3 text-center font-bold text-slate-700">{{ part.total_sold_qty }} pcs</td>
                  <td class="py-3 text-right font-mono font-bold text-red-600">{{ formatRupiah(part.total_revenue) }}</td>
                </tr>
                <tr v-if="!dashboardStore.summary.top_selling_parts || dashboardStore.summary.top_selling_parts.length === 0">
                  <td colspan="4" class="py-6 text-center text-slate-400">Belum ada riwayat transaksi suku cadang.</td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>

        <!-- Mechanic Rankings & Productivity -->
        <div class="bg-white p-6 rounded-2xl border border-slate-200/80 shadow-xs">
          <div class="flex items-center justify-between mb-4">
            <div>
              <h3 class="font-bold text-slate-800 text-base">Produktivitas Mekanik AHASS</h3>
              <p class="text-xs text-slate-500">Peringkat pengerjaan motor dan omset jasa</p>
            </div>
            <Award class="w-5 h-5 text-amber-500" />
          </div>

          <div class="space-y-3">
            <div 
              v-for="m in dashboardStore.summary.mechanic_rankings" 
              :key="m.mechanic_id"
              class="flex items-center justify-between p-3.5 rounded-xl border border-slate-100 bg-slate-50/50 hover:bg-slate-50 transition">
              <div class="flex items-center space-x-3">
                <div class="w-8 h-8 rounded-lg bg-red-100 text-red-700 font-bold flex items-center justify-center text-xs">
                  Pit {{ m.pit_number }}
                </div>
                <div>
                  <h4 class="font-bold text-slate-800 text-xs">{{ m.name }}</h4>
                  <p class="text-[11px] text-slate-500 font-mono">{{ m.code }}</p>
                </div>
              </div>
              <div class="text-right">
                <p class="text-xs font-black text-slate-900">{{ m.completed_jobs }} Unit Servis</p>
                <p class="text-[11px] font-mono text-emerald-600 font-semibold">{{ formatRupiah(m.labor_revenue) }} Jasa</p>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- 7-Day Performance Trend -->
      <div class="bg-white p-6 rounded-2xl border border-slate-200/80 shadow-xs">
        <h3 class="font-bold text-slate-800 text-base mb-1">Tren Pendapatan & Kunjungan 7 Hari Terakhir</h3>
        <p class="text-xs text-slate-500 mb-6">Monitoring performa harian operasional AHASS NSS</p>

        <div class="grid grid-cols-7 gap-3">
          <div 
            v-for="d in dashboardStore.summary.daily_trend" 
            :key="d.date"
            class="bg-slate-50 border border-slate-200/60 p-3.5 rounded-xl text-center flex flex-col justify-between hover:border-red-300 transition">
            <p class="text-xs font-semibold text-slate-500">{{ d.date }}</p>
            <div class="my-3">
              <p class="text-base font-black text-slate-900 font-mono">{{ formatRupiah(d.total_revenue) }}</p>
            </div>
            <span class="inline-block bg-white text-slate-700 px-2 py-0.5 rounded text-[10px] font-bold border border-slate-200">
              {{ d.bike_count }} Motor
            </span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import {
  DollarSign,
  TrendingUp,
  Wrench,
  AlertTriangle,
  Package,
  Award,
  Bike,
  Calendar,
} from 'lucide-vue-next'
import { useDashboardStore } from '../stores/dashboard'
import { formatRupiah } from '../utils/format'
import SkeletonLoader from '../components/SkeletonLoader.vue'

const dashboardStore = useDashboardStore()

onMounted(() => {
  dashboardStore.fetchSummary()
})
</script>

