<template>
  <aside class="w-64 bg-slate-900 text-slate-300 flex flex-col justify-between border-r border-slate-800 select-none">
    <div>
      <!-- Brand Logo / AHASS NSS -->
      <div class="h-20 flex items-center px-6 border-b border-slate-800/80 bg-slate-950/40">
        <div class="flex items-center space-x-3">
          <div class="w-10 h-10 rounded-xl bg-gradient-to-tr from-red-700 to-red-500 flex items-center justify-center shadow-lg shadow-red-900/30 text-white font-black text-xl">
            M
          </div>
          <div>
            <div class="flex items-center space-x-1.5">
              <span class="font-extrabold text-white text-base tracking-tight">MotoFlow</span>
              <span class="bg-red-600/90 text-white text-[9px] font-black px-1.5 py-0.5 rounded tracking-wide">AHASS</span>
            </div>
            <p class="text-[10px] text-slate-400 font-medium tracking-wide">Nusantara Sakti Group</p>
          </div>
        </div>
      </div>

      <!-- Navigation Menu -->
      <div class="p-4 space-y-1">
        <p class="px-3 pb-2 text-[10px] font-bold text-slate-500 uppercase tracking-widest">Menu Bengkel</p>

        <router-link
          v-for="item in menuItems"
          :key="item.to"
          :to="item.to"
          class="flex items-center space-x-3 px-3 py-2.5 rounded-xl text-sm font-medium transition-all group"
          :class="$route.path.startsWith(item.to) ? 'bg-red-600 text-white font-semibold shadow-md shadow-red-900/20' : 'text-slate-400 hover:text-white hover:bg-slate-800/60'">
          <component :is="item.icon" class="w-4 h-4 transition-transform group-hover:scale-110" />
          <span>{{ item.label }}</span>
        </router-link>
      </div>
    </div>

    <!-- User Profile & System Info -->
    <div class="p-4 border-t border-slate-800/80 bg-slate-950/30">
      <div class="flex items-center justify-between">
        <div class="flex items-center space-x-3">
          <div class="w-9 h-9 rounded-full bg-slate-800 border border-slate-700 flex items-center justify-center font-bold text-red-500 text-xs">
            {{ userInitial }}
          </div>
          <div class="overflow-hidden">
            <p class="text-xs font-semibold text-white truncate max-w-[120px]">{{ authStore.user?.full_name || 'Petugas AHASS' }}</p>
            <span class="inline-block text-[10px] font-bold text-red-400 uppercase tracking-wider">
              {{ authStore.user?.role?.replace('_', ' ') || 'Staff' }}
            </span>
          </div>
        </div>
        <button 
          @click="logout" 
          title="Keluar" 
          class="p-2 text-slate-400 hover:text-red-400 hover:bg-slate-800 rounded-lg transition">
          <LogOut class="w-4 h-4" />
        </button>
      </div>
    </div>
  </aside>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRouter } from 'vue-router'
import {
  LayoutDashboard,
  ClipboardList,
  PlusCircle,
  PackageCheck,
  Users,
  Wrench,
  LogOut,
} from 'lucide-vue-next'
import { useAuthStore } from '../stores/auth'

const router = useRouter()
const authStore = useAuthStore()

const menuItems = [
  { label: 'Executive Dashboard', to: '/dashboard', icon: LayoutDashboard },
  { label: 'Work Orders / Kasir', to: '/work-orders', icon: ClipboardList },
  { label: 'Pendaftaran Servis', to: '/work-orders/create', icon: PlusCircle },
  { label: 'Suku Cadang HGP', to: '/inventory', icon: PackageCheck },
  { label: 'Data Pelanggan & Motor', to: '/customers', icon: Users },
  { label: 'Mekanik & Pit Servis', to: '/mechanics', icon: Wrench },
]

const userInitial = computed(() => {
  if (!authStore.user?.full_name) return 'U'
  return authStore.user.full_name.charAt(0).toUpperCase()
})

function logout() {
  authStore.logout()
  router.push('/login')
}
</script>

