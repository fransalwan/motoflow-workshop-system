<template>
  <div>
    <!-- Desktop Static Sidebar -->
    <aside class="hidden lg:flex w-64 bg-slate-900 text-slate-300 flex-col justify-between border-r border-slate-800 select-none h-screen shrink-0">
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

          <!-- TV Monitor Ruang Tunggu Link -->
          <div class="pt-4 mt-2 border-t border-slate-800">
            <router-link
              to="/waiting-display"
              class="flex items-center space-x-3 px-3 py-2.5 rounded-xl text-xs font-semibold bg-slate-800/50 hover:bg-slate-800 text-red-400 border border-red-500/20 transition-all">
              <Tv class="w-4 h-4 text-red-500" />
              <span>TV Display Ruang Tunggu</span>
            </router-link>
          </div>
        </div>
      </div>

      <!-- User Profile & Logout -->
      <div class="p-4 border-t border-slate-800/80 bg-slate-950/30">
        <div class="flex items-center justify-between">
          <div class="flex items-center space-x-3 overflow-hidden">
            <div class="w-9 h-9 rounded-full bg-slate-800 border border-slate-700 flex items-center justify-center font-bold text-red-500 text-xs shrink-0">
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
            class="p-2 text-slate-400 hover:text-red-400 hover:bg-slate-800 rounded-lg transition shrink-0">
            <LogOut class="w-4 h-4" />
          </button>
        </div>
      </div>
    </aside>

    <!-- Mobile Drawer Overlay -->
    <div 
      v-if="isOpen" 
      class="fixed inset-0 z-50 lg:hidden bg-slate-950/70 backdrop-blur-sm flex">
      <div class="w-72 bg-slate-900 text-slate-300 flex flex-col justify-between h-full shadow-2xl p-4 border-r border-slate-800 animate-in slide-in-from-left duration-200">
        <div>
          <!-- Mobile Brand & Close Button -->
          <div class="flex items-center justify-between pb-4 border-b border-slate-800">
            <div class="flex items-center space-x-2">
              <div class="w-8 h-8 rounded-lg bg-red-600 flex items-center justify-center text-white font-black text-sm">
                M
              </div>
              <span class="font-black text-white text-base">MotoFlow AHASS</span>
            </div>
            <button @click="$emit('close')" class="p-1.5 text-slate-400 hover:text-white rounded-lg">
              <X class="w-5 h-5" />
            </button>
          </div>

          <!-- Navigation Links -->
          <div class="py-4 space-y-1">
            <router-link
              v-for="item in menuItems"
              :key="item.to"
              :to="item.to"
              @click="$emit('close')"
              class="flex items-center space-x-3 px-3 py-3 rounded-xl text-sm font-medium transition"
              :class="$route.path.startsWith(item.to) ? 'bg-red-600 text-white font-bold' : 'text-slate-400 hover:bg-slate-800'">
              <component :is="item.icon" class="w-4 h-4" />
              <span>{{ item.label }}</span>
            </router-link>

            <router-link
              to="/waiting-display"
              @click="$emit('close')"
              class="flex items-center space-x-3 px-3 py-3 rounded-xl text-xs font-bold text-red-400 bg-red-950/40 border border-red-800/50 mt-3">
              <Tv class="w-4 h-4" />
              <span>TV Display Ruang Tunggu</span>
            </router-link>
          </div>
        </div>

        <!-- Mobile User & Logout -->
        <div class="pt-4 border-t border-slate-800 flex items-center justify-between">
          <div class="flex items-center space-x-2">
            <div class="w-8 h-8 rounded-full bg-slate-800 border border-slate-700 flex items-center justify-center text-red-500 font-bold text-xs">
              {{ userInitial }}
            </div>
            <div>
              <p class="text-xs font-bold text-white">{{ authStore.user?.full_name || 'Petugas' }}</p>
              <p class="text-[10px] text-red-400 uppercase font-mono">{{ authStore.user?.role }}</p>
            </div>
          </div>
          <button @click="logout" class="p-2 text-slate-400 hover:text-red-400">
            <LogOut class="w-4 h-4" />
          </button>
        </div>
      </div>

      <!-- Backdrop Click to Close -->
      <div class="flex-1" @click="$emit('close')"></div>
    </div>
  </div>
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
  Tv,
  LogOut,
  X,
  ShoppingBag,
  Bike,
} from 'lucide-vue-next'
import { useAuthStore } from '../stores/auth'

defineProps<{
  isOpen?: boolean
}>()

defineEmits(['close'])

const router = useRouter()
const authStore = useAuthStore()

const menuItems = [
  { label: 'Executive Dashboard', to: '/dashboard', icon: LayoutDashboard },
  { label: 'Work Orders Servis', to: '/work-orders', icon: ClipboardList },
  { label: 'POS Retail Part & Oli', to: '/pos-retail', icon: ShoppingBag },
  { label: 'Showroom Motor & NSC', to: '/showroom', icon: Bike },
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
