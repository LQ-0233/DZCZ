import router from './router';
import { useTokenStore } from '/@/store/token';

const whiteList = ['/login'];

router.beforeEach(async (to, from, next) => {
  const tokenStore = useTokenStore()
  if (tokenStore.getToken && tokenStore.getToken != '') {
    if (to.path == '/login') {
      next({ path: '/' })
    } else {
      next()
    }
  } else {
    if (whiteList.indexOf(to.path) !== -1) {
      // 在免登录白名单，直接进入
      next()
    }else{
      console.log('token不存在')
      next({ path: '/login' })
    }
  }
})

router.afterEach(() => { })
