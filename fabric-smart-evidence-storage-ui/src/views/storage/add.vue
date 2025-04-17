<template>
  <div class="upload-container">
    <div class="title text-center">
      <h1>上传电子证据</h1>
    </div>
    <el-divider />
    <el-form :model="uploadForm" label-width="120px" :rules="rules">
      <el-form-item label="案件编号" prop="caseNumber">
        <el-input v-model="uploadForm.caseNumber" placeholder="请输入15位案件编号"></el-input>
      </el-form-item>
      
      <el-form-item label="案件基本信息" prop="caseInfo">
        <el-input
          type="textarea"
          v-model="uploadForm.caseInfo"
          placeholder="请输入案件基本信息"
          :maxlength="400"
          show-word-limit
        ></el-input>
      </el-form-item>

      <el-form-item label="负责人1" prop="manager1">
        <el-input v-model="uploadForm.manager1" placeholder="请输入负责人1姓名"></el-input>
      </el-form-item>

      <el-form-item label="负责人2" prop="manager2">
        <el-input v-model="uploadForm.manager2" placeholder="请输入负责人2姓名"></el-input>
      </el-form-item>

      <el-form-item label="电子证据形式" prop="evidenceType">
        <el-select v-model="uploadForm.evidenceType" placeholder="请选择电子证据形式">
          <el-option label="图片" value="image"></el-option>
          <el-option label="视频" value="video"></el-option>
          <el-option label="源代码" value="code"></el-option>
          <el-option label="文字形式" value="text"></el-option>
        </el-select>
      </el-form-item>

      <el-form-item label="文件" prop="fileUrl">
        <el-upload
          class="upload-demo"
           action="/api/v1/common/upload"
          :limit="1"
          :on-exceed="handleExceed"
          :on-success="handleSuccess"
          :before-upload="beforeUpload"
        >
          <el-button type="primary">点击上传</el-button>
          <template #tip>
            <div class="el-upload__tip">只能上传单个文件</div>
          </template>
        </el-upload>
      </el-form-item>

    </el-form>
    <el-divider />
    <el-button class="w-full" type="primary" @click="submitForm">提交</el-button>

  </div>
</template>

<script setup>
import { ref } from 'vue';
import { useRouter } from 'vue-router';
import api from '/@/api/api';
const router = useRouter()
const uploadForm = ref({
  caseNumber: '',
  caseInfo: '',
  manager1: '',
  manager2: '',
  evidenceType: '',
  ipfsLink: ''
})

const rules = {
  caseNumber: [
    { required: true, message: '请输入案件编号', trigger: 'blur' },
    { min: 15, max: 15, message: '案件编号必须为15位', trigger: 'blur' }
  ],
  caseInfo: [
    { required: true, message: '请输入案件基本信息', trigger: 'blur' },
    { max: 400, message: '案件基本信息不能超过400字', trigger: 'blur' }
  ],
  manager1: [
    { required: true, message: '请输入负责人1', trigger: 'blur' }
  ],
  manager2: [
    { required: true, message: '请输入负责人2', trigger: 'blur' }
  ],
  evidenceType: [
    { required: true, message: '请选择电子证据形式', trigger: 'change' }
  ],
  fileUrl: [
    { required: true, message: '请上传电子证据', trigger: 'blur' }
  ]
}

const handleExceed = () => {
  ElMessage.warning('只能上传一个文件')
}

const handleSuccess = (response) => {
  uploadForm.value.filePath = response.filePath
  uploadForm.value.fileName = response.fileName
  ElMessage.success('上传成功')
}

const beforeUpload = (file) => {
  return true
}

const submitForm = () => {
  console.log('表单数据：', uploadForm.value)
  api.evidence.create(uploadForm.value).then((res) => {
    ElMessage.success('提交成功')
    router.push('/storage/index')
  })
}

</script>

<style scoped>
.upload-container {
  max-width: 600px;
  margin: 20px auto;
  padding: 20px;
  background-color: #f5f5f5;
  border-radius: 10px;
}

.upload-demo {
  text-align: center;
}

.el-upload__tip {
  margin-top: 8px;
}
</style>
