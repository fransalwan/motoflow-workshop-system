import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '../stores/auth'
import { useToastStore } from '../stores/toast'
import type { Role } from '../types'

export function getDefaultRouteForRole(role?: Role): string {
  switch (role) {
    case 'admin':
      return '/dashboard'
    case 'service_advisor':
      return '/work-orders'
    case 'cashier':
      return '/pos-retail'
    case 'mechanic':
      return '/work-orders'
    default:
      return '/login'
  }
}

const routes = [
  {
    path: '/login',
    name: 'Login',
    component: () => import('../pages/LoginPage.vue'),
    meta: { public: true },
  },
  {
    path: '/waiting-display',
    name: 'WaitingDisplay',
    component: () => import('../pages/WaitingDisplayPage.vue'),
    meta: { public: true },
  },
  {
    path: '/',
    component: () => import('../layouts/MainLayout.vue'),
    children: [
      {
        path: 'dashboard',
        name: 'Dashboard',
        component: () => import('../pages/DashboardPage.vue'),
        meta: { roles: ['admin'] as Role[] },
      },
      {
        path: 'work-orders',
        name: 'WorkOrders',
        component: () => import('../pages/WorkOrderListPage.vue'),
        meta: { roles: ['admin', 'service_advisor', 'cashier', 'mechanic'] as Role[] },
      },
      {
        path: 'work-orders/create',
        name: 'WorkOrderCreate',
        component: () => import('../pages/WorkOrderCreatePage.vue'),
        meta: { roles: ['admin', 'service_advisor'] as Role[] },
      },
      {
        path: 'work-orders/:id',
        name: 'WorkOrderDetail',
        component: () => import('../pages/WorkOrderDetailPage.vue'),
        meta: { roles: ['admin', 'service_advisor', 'cashier', 'mechanic'] as Role[] },
      },
      {
        path: 'pos-retail',
        name: 'PosRetail',
        component: () => import('../pages/PosRetailPage.vue'),
        meta: { roles: ['admin', 'cashier', 'service_advisor'] as Role[] },
      },
      {
        path: 'showroom',
        name: 'Showroom',
        component: () => import('../pages/ShowroomPage.vue'),
        meta: { roles: ['admin', 'service_advisor', 'cashier'] as Role[] },
      },
      {
        path: 'inventory',
        name: 'Inventory',
        component: () => import('../pages/InventoryPage.vue'),
        meta: { roles: ['admin', 'cashier', 'service_advisor'] as Role[] },
      },
      {
        path: 'customers',
        name: 'Customers',
        component: () => import('../pages/CustomerListPage.vue'),
        meta: { roles: ['admin', 'service_advisor'] as Role[] },
      },
      {
        path: 'mechanics',
        name: 'Mechanics',
        component: () => import('../pages/MechanicListPage.vue'),
        meta: { roles: ['admin', 'service_advisor', 'mechanic'] as Role[] },
      },
    ],
  },
  {
    path: '/:pathMatch(.*)*',
    redirect: (to: any) => {
      const authStore = useAuthStore()
      return getDefaultRouteForRole(authStore.user?.role)
    },
  },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

router.beforeEach((to, _from, next) => {
  const authStore = useAuthStore()

  // 1. Public route handling
  if (to.meta.public) {
    if (to.name === 'Login' && authStore.isAuthenticated) {
      return next(getDefaultRouteForRole(authStore.user?.role))
    }
    return next()
  }

  // 2. Unauthenticated user handling
  if (!authStore.isAuthenticated) {
    return next({ name: 'Login' })
  }

  // 3. Root path '/' handling
  if (to.path === '/') {
    return next(getDefaultRouteForRole(authStore.user?.role))
  }

  // 4. Role Authorization Check
  const allowedRoles = to.meta.roles as Role[] | undefined
  if (allowedRoles && allowedRoles.length > 0) {
    const userRole = authStore.user?.role
    if (!userRole || !allowedRoles.includes(userRole)) {
      const toast = useToastStore()
      const roleLabel = userRole ? userRole.replace('_', ' ').toUpperCase() : 'ANDA'
      toast.error(`Akses Terbatas: Peran ${roleLabel} tidak memiliki wewenang untuk halaman ini.`)
      return next(getDefaultRouteForRole(userRole))
    }
  }

  next()
})

export default router
