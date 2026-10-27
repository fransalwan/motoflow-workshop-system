import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '../stores/auth'

const routes = [
  {
    path: '/login',
    name: 'Login',
    component: () => import('../pages/LoginPage.vue'),
    meta: { public: true },
  },
  {
    path: '/',
    component: () => import('../layouts/MainLayout.vue'),
    redirect: '/dashboard',
    children: [
      {
        path: 'dashboard',
        name: 'Dashboard',
        component: () => import('../pages/DashboardPage.vue'),
      },
      {
        path: 'work-orders',
        name: 'WorkOrders',
        component: () => import('../pages/WorkOrderListPage.vue'),
      },
      {
        path: 'work-orders/create',
        name: 'WorkOrderCreate',
        component: () => import('../pages/WorkOrderCreatePage.vue'),
      },
      {
        path: 'work-orders/:id',
        name: 'WorkOrderDetail',
        component: () => import('../pages/WorkOrderDetailPage.vue'),
      },
      {
        path: 'inventory',
        name: 'Inventory',
        component: () => import('../pages/InventoryPage.vue'),
      },
      {
        path: 'customers',
        name: 'Customers',
        component: () => import('../pages/CustomerListPage.vue'),
      },
      {
        path: 'mechanics',
        name: 'Mechanics',
        component: () => import('../pages/MechanicListPage.vue'),
      },
    ],
  },
  {
    path: '/:pathMatch(.*)*',
    redirect: '/dashboard',
  },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

router.beforeEach((to, _from, next) => {
  const authStore = useAuthStore()
  if (!to.meta.public && !authStore.isAuthenticated) {
    next({ name: 'Login' })
  } else if (to.name === 'Login' && authStore.isAuthenticated) {
    next({ name: 'Dashboard' })
  } else {
    next()
  }
})

export default router

