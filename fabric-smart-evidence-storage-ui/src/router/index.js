import {
  createRouter,
  createWebHistory
} from "vue-router"

export const realityRoutes = [{
    path: "/login",
    name: "Login",
    meta: {
      name: "登录",
      show: false,
    },
    component: () => import("/@/views/login.vue")
  },
  {
    path: "/",
    redirect: "/user/list",
  },
  {
    path: "/user",
    name: "User",
    redirect: "/user/list",
    component: () => import("/@/layout/index.vue"),
    meta: {
      name: "用户管理",
      show: true,
      role: "1",
    },
    children: [{
      path: "list",
      name: "UserList",
      meta: {
        name: "用户列表",
        show: true,
      },
      component: () => import("/@/views/user/index.vue"),
    }, ]
  },
  {
    path: "/storage",
    name: "Storage",
    redirect: "/storage/index",
    component: () => import("/@/layout/index.vue"),
    meta: {
      name: "电子证据管理",
      show: true,
      flat: true,
      role: "2",
    },
    children: [
      {
        path: "add",
        name: "StorageAdd",
        meta: {
          name: "上传电子证据",
          show: true,
        },
        component: () => import("/@/views/storage/add.vue"),
      },
      {
        path: "index",
        name: "StorageIndex",
        meta: {
          name: "电子证据列表",
          show: true,
        },
        component: () => import("/@/views/storage/list.vue"),
      },
    ]
  },
  {
    path: "/authorized",
    name: "Authorized",
    redirect: "/authorized/index",
    component: () => import("/@/layout/index.vue"),
    meta: {
      name: "授权管理",
      show: true,
      flat: true,
    },
    children: [
      {
        path: "index",
        name: "AuthorizedIndex",
        meta: {
          name: "授权管理",
          show: true,
          role: "2",
        },
        component: () => import("/@/views/authorized/list.vue"),
      },
      {
        path: "view",
        name: "AuthorizedView",
        meta: {
          name: "可查阅案件",
          show: true,
          hiddenOnAdmin: true
        },
        component: () => import("../views/authorized/view.vue"),
      },
      {
        path: "record",
        name: "AuthorizedRecord",
        meta: {
          name: "查看纪录",
          show: true,
          role: "2",
        },
        component: () => import("../views/authorized/record.vue"),
      },
    ]
  },

  {
    path: "/:pathMatch(.*)*",
    name: "NotFound",
    meta: {
      name: "404",
      show: false,
    },
    redirect: "/",
  }
]



const router = createRouter({
  history: createWebHistory(),
  routes: realityRoutes
})


export default router