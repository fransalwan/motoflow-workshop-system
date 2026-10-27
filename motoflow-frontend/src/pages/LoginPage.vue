<template>
  <div class="min-h-screen bg-slate-950 flex flex-col justify-center items-center p-4 relative overflow-hidden">
    <!-- Red glowing ambient accents -->
    <div class="absolute -top-40 -left-40 w-96 h-96 bg-red-600/20 rounded-full blur-3xl pointer-events-none"></div>
    <div class="absolute -bottom-40 -right-40 w-96 h-96 bg-red-800/20 rounded-full blur-3xl pointer-events-none"></div>

    <div class="max-w-md w-full relative z-10">
      <!-- Brand Card Header -->
      <div class="text-center mb-8">
        <div class="inline-flex items-center justify-center w-16 h-16 rounded-2xl bg-gradient-to-tr from-red-700 to-red-500 shadow-xl shadow-red-900/40 text-white font-black text-3xl mb-4 border border-red-400/30">
          M
        </div>
        <h1 class="text-2xl font-black text-white tracking-tight flex items-center justify-center space-x-2">
          <span>MotoFlow</span>
          <span class="bg-red-600 text-white text-xs font-black px-2 py-0.5 rounded uppercase tracking-wider">AHASS</span>
        </h1>
        <p class="text-xs text-slate-400 font-medium mt-1">Sistem Bengkel Resmi Honda & POS Billing PT Nusantara Sakti Group</p>
      </div>

      <!-- Main Login Form -->
      <div class="bg-slate-900/90 backdrop-blur-md p-8 rounded-3xl border border-slate-800 shadow-2xl">
        <div v-if="authStore.error" class="mb-5 p-3.5 bg-red-950/60 border border-red-800/80 rounded-xl text-red-200 text-xs flex items-center space-x-2">
          <AlertCircle class="w-4 h-4 shrink-0 text-red-400" />
          <span>{{ authStore.error }}</span>
        </div>

        <form @submit.prevent="handleLogin" class="space-y-4">
          <div>
            <label class="block text-xs font-semibold text-slate-300 uppercase tracking-wider mb-1.5">Username / NIK</label>
            <div class="relative">
              <input
                v-model="username"
                type="text"
                required
                placeholder="Masukkan username"
                class="w-full bg-slate-950 border border-slate-700 rounded-xl px-4 py-2.5 text-sm text-white placeholder-slate-500 focus:outline-none focus:border-red-500 transition" />
              <UserIcon class="w-4 h-4 text-slate-500 absolute right-3.5 top-3" />
            </div>
          </div>

          <div>
            <label class="block text-xs font-semibold text-slate-300 uppercase tracking-wider mb-1.5">Kata Sandi</label>
            <div class="relative">
              <input
                v-model="password"
                type="password"
                required
                placeholder="Masukkan kata sandi"
                class="w-full bg-slate-950 border border-slate-700 rounded-xl px-4 py-2.5 text-sm text-white placeholder-slate-500 focus:outline-none focus:border-red-500 transition" />
              <Lock class="w-4 h-4 text-slate-500 absolute right-3.5 top-3" />
            </div>
          </div>

          <button
            type="submit"
            :disabled="authStore.loading"
            class="w-full mt-2 bg-gradient-to-r from-red-600 to-red-700 hover:from-red-500 hover:to-red-600 text-white font-bold py-3 px-4 rounded-xl shadow-lg shadow-red-900/30 text-sm transition active:scale-[0.98] disabled:opacity-50">
            <span v-if="!authStore.loading">Masuk ke Sistem AHASS</span>
            <span v-else>Memverifikasi Kredensial...</span>
          </button>
        </form>

        <!-- Quick Demo Role Switcher for HRD / Technical Interviewers -->
        <div class="mt-6 pt-6 border-t border-slate-800">
          <p class="text-[11px] font-bold text-slate-400 uppercase tracking-wider text-center mb-3">
            Akses Cepat Akun Demo (1-Click)
          </p>
          <div class="grid grid-cols-2 gap-2 text-xs">
            <button
              @click="setCredentials('admin', 'password123')"
              type="button"
              class="bg-slate-800/80 hover:bg-slate-800 text-slate-200 border border-slate-700 p-2 rounded-lg text-left transition hover:border-red-500/50">
              <p class="font-bold text-white text-[11px]">Kepala Bengkel</p>
              <p class="text-[10px] text-slate-400">admin / pass123</p>
            </button>
            <button
              @click="setCredentials('advisor', 'password123')"
              type="button"
              class="bg-slate-800/80 hover:bg-slate-800 text-slate-200 border border-slate-700 p-2 rounded-lg text-left transition hover:border-red-500/50">
              <p class="font-bold text-white text-[11px]">Service Advisor</p>
              <p class="text-[10px] text-slate-400">advisor / pass123</p>
            </button>
            <button
              @click="setCredentials('kasir', 'password123')"
              type="button"
              class="bg-slate-800/80 hover:bg-slate-800 text-slate-200 border border-slate-700 p-2 rounded-lg text-left transition hover:border-red-500/50">
              <p class="font-bold text-white text-[11px]">Kasir POS AHASS</p>
              <p class="text-[10px] text-slate-400">kasir / pass123</p>
            </button>
            <button
              @click="setCredentials('mekanik', 'password123')"
              type="button"
              class="bg-slate-800/80 hover:bg-slate-800 text-slate-200 border border-slate-700 p-2 rounded-lg text-left transition hover:border-red-500/50">
              <p class="font-bold text-white text-[11px]">Mekanik Senior</p>
              <p class="text-[10px] text-slate-400">mekanik / pass123</p>
            </button>
          </div>
        </div>
      </div>

      <!-- Footer Info -->
      <p class="text-center text-[11px] text-slate-600 mt-6">
        Portofolio Programmer Junior NSS — Dibuat dengan Go & Vue 3
      </p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { User as UserIcon, Lock, AlertCircle } from 'lucide-vue-next'
import { useAuthStore } from '../stores/auth'

const router = useRouter()
const authStore = useAuthStore()

const username = ref('admin')
const password = ref('password123')

function setCredentials(u: string, p: string) {
  username.value = u
  password.value = p
  handleLogin()
}

async function handleLogin() {
  const success = await authStore.login(username.value, password.value)
  if (success) {
    router.push('/dashboard')
  }
}
</script>

