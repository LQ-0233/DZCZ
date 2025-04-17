<template>
  <el-container class="h-full min-h-screen">
    <el-aside class="h-full min-h-screen w-1/5 max-w-50 flex flex-col bg-color-#3f3f3f color-white">
      <div class="h-20 flex p-2">
        <img :src="logo" class="mr-1 w-5" />
        <h5>基于联盟链的智慧司法证据存储系统</h5>
      </div>
      <el-menu :default-active="activeIndex" class="h-full" background-color="#3f3f3f" text-color="white" router>
        <template v-for="router in routes" :key="router.path">
          <el-menu-item v-if="!router.children || router.children.length == 0" :index="router.path">
            <span>{{ router.meta.name }}</span>
          </el-menu-item>
          <el-sub-menu v-else :index="router.path">
            <template #title>
              <span>{{ router.meta.name }}</span>
            </template>
            <el-menu-item :index="getAccessPath(router, child)" v-for="child in router.children">
              {{
                child.meta.name
              }}
            </el-menu-item>
          </el-sub-menu>
        </template>
      </el-menu>
    </el-aside>
    <el-container>
      <el-main class="w-full min-h-screen h-screen bg-slate-300" style="padding: 0;">
        <div class="bg-white mb-5 p-3 flex justify-end">
          <el-icon>
            <user-filled />
          </el-icon>
          <el-dropdown @command="handleCommand">
            <span class="el-dropdown-link">
              {{ tokenStore.getNickname }}
              <el-icon class="el-icon--right">
                <arrow-down />
              </el-icon>
            </span>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item command="changePassword">修改密码</el-dropdown-item>
                <el-dropdown-item command="logout">退出登录</el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </div>
        <div class="min-h-9/10 p-5 mb-2">
          <!-- <router-view :key="key" /> -->
          <router-view v-slot="{ Component, route }">
            <transition name="fade-transform" mode="out-in">
              <component :is="Component" :key="route.path" />
            </transition>
          </router-view>
        </div>

      </el-main>
    </el-container>
    <el-dialog v-model="changePasswordVisiable" title="修改密码" width="400px">
      <el-form :model="form" ref="formRef" label-width="80px" label-position="left">
        <el-form-item label="旧密码: " prop="oldPwd">
          <el-input v-model="form.oldPwd" autocomplete="off" type="password" />
        </el-form-item>
        <el-form-item label="新密码: " prop="pwd">
          <el-input v-model="form.pwd" autocomplete="off" type="password" />
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button type="primary" @click="changePassword">确定</el-button>
        </span>
      </template>
    </el-dialog>
  </el-container>
</template>

<script setup>
import _ from 'lodash';
import { computed, onMounted, ref } from 'vue';
import { useRoute } from 'vue-router';
import logo from '/@/assets/logo.svg';
import api from '/@/api/api';
import { useTokenStore } from '/@/store/token';
const tokenStore = useTokenStore()

import { realityRoutes } from '/@/router/index';
const router = useRoute()
const routes = ref([])

import { useRouter } from 'vue-router';
// setup
const { currentRoute, push } = useRouter()
const activeIndex = ref(currentRoute.value.path)

onMounted(async () => {
  console.log("routes: ", realityRoutes);
  let tmp = []
  for (const route of realityRoutes) {
    if (route.meta && route.meta.show === true && (!route.meta.role || route.meta.role == tokenStore.getRole)) {
      if (route.meta.flat) {
        for (const child of route.children) {
          if (!child.meta || !child.meta.show || child.meta.show !== true) {
            continue
          }
          if(!(!child.meta.role || child.meta.role == tokenStore.getRole)){
            continue
          }
          if(child.meta.hiddenOnAdmin && tokenStore.getRole == "1"){
            continue
          }
          const _child = _.cloneDeep(child)
          let middle = "/"
          if (route.path.endsWith("/")) {
            middle = ""
          }
          _child.path = route.path + middle + _child.path
          tmp.push(_child)
        }
      } else {
        tmp.push(route)
      }
    }
  }
  tmp = _.cloneDeep(tmp)
  console.log("tmp: ", tmp);
  tmp.forEach(route => {
    route.children = route.children?.filter(child => {
      const result = child.meta && child.meta.show && child.meta.show === true
        && (!child.meta.role || child.meta.role == tokenStore.getRole) 
      // console.log("child: ", child, result);
      // console.log("child show: ", child.meta.role == null || child.meta.role == info.value.role);

      return result
    })
  })
  routes.value = tmp
  console.log("routes: ", routes.value)
});


const key = computed(() => useRoute().name)

const getAccessPath = (parent, child) => {
  if (parent.path == "/") {
    return child.path
  }
  return parent.path + "/" + child.path
}


const changePasswordVisiable = ref(false)

const form = ref({
  oldPwd: "",
  pwd: ""
})

const changePassword = () => {
  if (form.value.pwd == "") {
    ElMessage.error("密码不能为空")
    return
  }
  // 旧密码不能和新密码相同
  if (form.value.oldPwd == form.value.pwd) {
    ElMessage.error("新密码不能和旧密码相同")
    return
  }
  api.user.changePwd(form.value).then(res => {
    ElMessage.success("修改成功")
    changePasswordVisiable.value = false
  }).catch(err => {
    console.log(err);
  })
}

const handleCommand = (command) => {
  if (command == 'logout') {
    tokenStore.setToken("")
    push("/login")
  }
  if (command == 'changePassword') {
    changePasswordVisiable.value = true
    form.value = {
      oldPwd: "",
      pwd: "",
    }
  }
}

</script>

<style scoped>
.el-menu {
  border-right: 0 !important;
}
</style>