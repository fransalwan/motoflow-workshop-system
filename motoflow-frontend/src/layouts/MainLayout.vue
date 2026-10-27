<template>
  <div class="flex h-screen bg-slate-50 overflow-hidden font-sans">
    <!-- Sidebar (Desktop & Mobile Drawer) -->
    <Sidebar 
      :is-open="isMobileMenuOpen" 
      @close="isMobileMenuOpen = false" />

    <!-- Main Content Area -->
    <div class="flex-1 flex flex-col min-w-0 overflow-hidden">
      <Navbar @toggle-mobile-menu="isMobileMenuOpen = !isMobileMenuOpen" />

      <main class="flex-1 overflow-y-auto p-3.5 sm:p-6 lg:p-8 pb-28 lg:pb-8">
        <router-view v-slot="{ Component }">
          <transition name="fade" mode="out-in">
            <component :is="Component" />
          </transition>
        </router-view>
      </main>

      <!-- Mobile Bottom Navigation Bar (Dynamic Role-Based & Safe-Area Aware) -->
      <nav class="lg:hidden fixed bottom-0 left-0 right-0 bg-white/95 backdrop-blur-md border-t border-slate-200/90 flex items-center justify-around px-2 z-40 shadow-lg pb-[max(0.5rem,env(safe-area-inset-bottom))] pt-1">
        <router-link
          v-for="item in visibleBottomItems"
          :key="item.to"
          :to="item.to"
          class="flex flex-col items-center justify-center transition active:scale-90"
          :class="item.isCenter ? '-mt-5' : 'w-14 py-1 text-[10px] font-semibold ' + ($route.path === item.to ? 'text-red-600 font-bold' : 'text-slate-400')">
          <div v-if="item.isCenter" class="w-12 h-12 rounded-full bg-gradient-to-tr from-red-700 to-red-500 text-white flex items-center justify-center shadow-lg shadow-red-900/30 border-4 border-white transition-transform">
            <component :is="item.icon" class="w-5 h-5 stroke-[2.2]" />
          </div>
          <component v-else :is="item.icon" class="w-5 h-5 mb-0.5" />
          <span :class="item.isCenter ? 'text-[9px] font-bold text-slate-700 mt-1' : ''">{{ item.label }}</span>
        </router-link>
      </nav>
    </div>

    <!-- Global Toast Notification Container -->
    <ToastContainer />
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import {
  LayoutDashboard,
  ClipboardList,
  Plus,
  PackageCheck,
  Users,
  ShoppingBag,
  Bike,
  Wrench,
} from 'lucide-vue-next'
import Sidebar from '../components/Sidebar.vue'
import Navbar from '../components/Navbar.vue'
import ToastContainer from '../components/ToastContainer.vue'
import { useAuthStore } from '../stores/auth'
import type { Role } from '../types'

const authStore = useAuthStore()
const isMobileMenuOpen = ref(false)

interface BottomItem {
  label: string
  to: string
  icon: any
  isCenter?: boolean
  roles: Role[]
}

const allBottomItems: BottomItem[] = [
  { label: 'Omzet', to: '/dashboard', icon: LayoutDashboard, roles: ['admin'] },
  { label: 'Servis', to: '/work-orders', icon: ClipboardList, roles: ['admin', 'service_advisor', 'cashier', 'mechanic'] },
  { label: 'Motor', to: '/customers', icon: Users, roles: ['service_advisor'] },
  { label: 'Daftar', to: '/work-orders/create', icon: Plus, isCenter: true, roles: ['service_advisor'] },
  { label: 'Retail', to: '/pos-retail', icon: ShoppingBag, isCenter: true, roles: ['admin', 'cashier'] },
  { label: 'Showroom', to: '/showroom', icon: Bike, roles: ['admin', 'service_advisor', 'cashier'] },
  { label: 'Pit Servis', to: '/mechanics', icon: Wrench, roles: ['mechanic'] },
  { label: 'HGP', to: '/inventory', icon: PackageCheck, roles: ['admin', 'cashier', 'service_advisor'] },
]

const visibleBottomItems = computed(() => {
  const role = authStore.user?.role
  if (!role) return []
  return allBottomItems.filter(item => item.roles.includes(role))
})
</script>

<style scoped>
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.15s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>
