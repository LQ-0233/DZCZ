import { defineStore } from 'pinia'


export const useTokenStore= defineStore('token',{
  state: () => ({
    token: localStorage.getItem('token') || '',
    role: localStorage.getItem('role') || '',
    nickname: localStorage.getItem('nickname') || '',
  }),
  getters: {
    getToken: (state) => state.token,
    getRole: (state) => state.role,
    getNickname: (state) => state.nickname,
  },
  actions: {
    setToken(data) {
      this.token = data
      localStorage.setItem('token', data)
    },
    setRole(data) {
      this.role = data
      localStorage.setItem('role', data)
    },
    setNickname(data) {
      this.nickname = data
      localStorage.setItem('nickname', data)
    },
  }
})