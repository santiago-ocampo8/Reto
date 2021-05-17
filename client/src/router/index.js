import Vue from 'vue'
import VueRouter from 'vue-router'
import Home from '../views/Home.vue'

Vue.use(VueRouter)

const routes = [
  {
    path: '/',
    name: 'Home',
    component: Home
  },
  {
    path: '/about',
    name: 'About',

    component: function () {
      return import(/* webpackChunkName: "about" */ '../views/About.vue')
    }
  },
  {
    path: '/buyers',
    name: 'Buyers',
    
    component: function () {
      return import( '../views/ListBuyers.vue')
    }
  },
  {
    path: '/detailsbuyer/:id',
    name: 'DetailsBuyer',
    
    component: function () {
      return import( '../views/DetailsBuyer.vue')
    }
  },
  {
    path: '/listbuyersip/:ip',
    name: 'ListBuyersIp',
    
    component: function () {
      return import( '../views/ListBuyersIp.vue')
    }
  },
  {
    path: '/load',
    name: 'LoadData',
    
    component: function () {
      return import( '../views/LoadData.vue')
    }
  }
]

const router = new VueRouter({
  mode: 'history',
  base: process.env.BASE_URL,
  routes
})

export default router
