<template>
  <div class="flex h-screen bg-slate-50 overflow-hidden font-sans">
    <!-- Sidebar (Desktop & Mobile Drawer) -->
    <Sidebar 
      :is-open="isMobileMenuOpen" 
      @close="isMobileMenuOpen = false" />

    <!-- Main Content Area -->
    <div class="flex-1 flex flex-col min-w-0 overflow-hidden">
      <Navbar @toggle-mobile-menu="isMobileMenuOpen = !isMobileMenuOpen" />

      <main class="flex-1 overflow-y-auto p-4 sm:p-6 lg:p-8 pb-24 lg:pb-8">
        <router-view v-slot="{ Component }">
          <transition name="fade" mode="out-in">
            <component :is="Component" />
          </transition>
        </router-view>
      </main>

      <!-- Mobile Bottom Navigation Bar (Thumb-friendly for smartphone users) -->
      <nav class="lg:hidden fixed bottom-0 left-0 right-0 h-16 bg-white/95 backdrop-blur-md border-t border-slate-200 flex items-center justify-around px-2 z-40 shadow-lg">
        <router-link
          to="/dashboard"
          class="flex flex-col items-center justify-center w-14 py-1 text-[10px] font-semibold transition"
          :class="$route.path === '/dashboard' ? 'text-red-600' : 'text-slate-400'">
          <LayoutDashboard class="w-5 h-5 mb-0.5" />
          <span>Omzet</span>
        </router-link>

        <router-link
          to="/work-orders"
          class="flex flex-col items-center justify-center w-14 py-1 text-[10px] font-semibold transition"
          :class="$route.path.startsWith('/work-orders') && $route.path !== '/work-orders/create' ? 'text-red-600' : 'text-slate-400'">
          <ClipboardList class="w-5 h-5 mb-0.5" />
          <span>Servis</span>
        </router-link>

        <!-- Center Raised Action Button (Buat Servis Baru) -->
        <router-link
          to="/work-orders/create"
          class="flex flex-col items-center justify-center -mt-5">
          <div class="w-12 h-12 rounded-full bg-gradient-to-tr from-red-700 to-red-500 text-white flex items-center justify-center shadow-lg shadow-red-900/30 border-4 border-white active:scale-90 transition">
            <Plus class="w-6 h-6 stroke-[2.5]" />
          </div>
          <span class="text-[9px] font-bold text-slate-700 mt-1">Daftar</span>
        </router-link>

        <router-link
          to="/inventory"
          class="flex flex-col items-center justify-center w-14 py-1 text-[10px] font-semibold transition"
          :class="$route.path === '/inventory' ? 'text-red-600' : 'text-slate-400'">
          <PackageCheck class="w-5 h-5 mb-0.5" />
          <span>HGP</span>
        </router-link>

        <router-link
          to="/customers"
          class="flex flex-col items-center justify-center w-14 py-1 text-[10px] font-semibold transition"
          :class="$route.path === '/customers' ? 'text-red-600' : 'text-slate-400'">
          <Users class="w-5 h-5 mb-0.5" />
          <span>Motor</span>
        </router-link>
      </nav>
    </div>

    <!-- Global Toast Notification Container -->
    <ToastContainer />
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import {
  LayoutDashboard,
  ClipboardList,
  Plus,
  PackageCheck,
  Users,
} from 'lucide-vue-next'
import Sidebar from '../components/Sidebar.vue'
import Navbar from '../components/Navbar.vue'
import ToastContainer from '../components/ToastContainer.vue'

const isMobileMenuOpen = ref(false)
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
