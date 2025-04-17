<script setup>
import { reactive, ref } from 'vue';
import { ElMessage, ElMessageBox } from 'element-plus';
import { useRouter } from "vue-router";
import api from "/@/api/api";
const router = useRouter()
import { useTokenStore } from '../store/token';
const tokenStore = useTokenStore()

const formData = ref({
    username: "",
    pwd: "",
    nickname: "",
    role: "",
})
const registFlag = ref(false)
const formRef = ref(null)
const formRules = ref({
    username: [
        { required: true, message: '请输入用户名', trigger: 'blur' },
    ],
    pwd: [
        { required: true, message: '请输入密码', trigger: 'blur' },

    ],
    nickname: [
        { required: true, message: '请输入名称', trigger: 'blur' },

    ],
    role: [
        { required: true, message: '请选择角色', trigger: 'blur' },
    ],

})



const login = async () => {
    registFlag.value = false
    formRef.value.validate((valid) => {
        if (valid) {
            api.common.login(formData.value).then((res) => {

                if (res) {
                    ElMessage.success('登录成功');
                    tokenStore.setToken(res.token)
                    tokenStore.setNickname(res.nickname)
                    tokenStore.setRole(res.role)
                    if (res.role == '1') {
                        router.push({ path: '/user/list' })
                    } else if (res.role == '2') {
                        router.push({ path: '/storage/add' })
                    } else {
                        router.push({ path: '/authorized/view' })
                    }
                }
            })
        } else {
            console.log('error submit!!');
            return false;
        }
    });
}

const register = async () => {
    if (registFlag.value) {
        formRef.value.validate((valid) => {
            if (valid) {
                api.common.register(formData.value).then((res) => {
                    
                    ElMessage.success('注册成功');
                        registFlag.value = false
                        formData.value = {
                            username: "",
                            pwd: "",
                            role: "",
                            nickname: ""
                        }
                    
                })
            } else {
                console.log('error submit!!');
                return false;
            }
        });
    } else {
        registFlag.value = true

    }

}

</script>
<template>
    <div class="wrapper">
        <div class="content" @keyup.enter="login">
            <div class="title">
                <h1>基于联盟链的智慧司法证据存储系统</h1>
            </div>
            <el-form ref="formRef" :model="formData" :rules="formRules" label-width="120px">
                <el-form-item label="用户名: " prop="username" class="item">
                    <el-input v-model="formData.username" style="width:400px" />
                </el-form-item>
                <el-form-item label="密码: " prop="pwd" class="item">
                    <el-input v-model="formData.pwd" type="password" style="width:400px" />
                </el-form-item>
                <el-form-item label="名称: " prop="nickname" class="item" v-if="registFlag">
                    <el-input v-model="formData.nickname" style="width:400px" />
                </el-form-item>
                <el-form-item label="角色: " prop="role" class="item" v-if="registFlag">
                    <el-select v-model="formData.role" style="width:400px">
                        <el-option label="查阅人员" value="3"></el-option>
                        <el-option label="录入人员" value="2"></el-option>
                    </el-select>
                </el-form-item>

            </el-form>
            <div class="flex justify-center mt-2">
                <el-button class="w-40" type="primary" @click="register">注册</el-button>
                <el-button class="w-40" type="primary" @click="login">登录</el-button>
            </div>
        </div>

    </div>
</template>
<style lang="scss">
.wrapper {
    width: 100%;
    height: 100%;
    display: flex;
    flex-direction: column;
    justify-content: center;
    align-items: center;
    background-image: url('/@/assets/bg.png');
    background-size: cover;

    .content {
        background-color: rgba(200,
                200,
                200,
                0.3);
        border-color: rgba(200, 200, 200, 0.5);
        border-radius: 18px;
        padding: 40px;
    }

    .title {
        text-align: center;
        color: black;
        margin-bottom: 1em;
    }


}

.item .el-form-item__label {
    color: black;
}
</style>