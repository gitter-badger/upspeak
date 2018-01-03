import Vue from 'vue'
import Router from 'vue-router'

Vue.use(Router)

export default new Router({
  mode: 'history',

  routes: [

    // Home
    {
      path: '/',
      name: 'Home',
      component: () => import('@/routes/Home.vue')
    },

    // Node detail
    {
      path: '/node/:nid',
      name: 'NodeDetail',
      component: () => import('@/routes/NodeDetail.vue')
    }
  ]
})
