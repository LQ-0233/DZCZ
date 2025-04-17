<template>
  <el-upload class="avatar-uploader" :http-request="handleUpload" accept="image/*"   list-type="picture-card"     :on-remove="handleRemove"     :on-preview="handlePictureCardPreview"
  :file-list="fileList"
    :limit="1"
    :on-exceed="onExceed">
    <el-icon class="avatar-uploader-icon">
      <Plus />
    </el-icon>
  </el-upload>
  <el-dialog v-model="dialogVisible">
    <img w-full :src="dialogImageUrl" alt="Preview Image" />
  </el-dialog>
</template>

<script setup>
import { Plus } from '@element-plus/icons-vue';
import { defineEmits, ref, watch } from 'vue';
import { create } from 'ipfs-http-client'
const client = create({
    host: "localhost",
    port: "5001",
    protocol: "http"
});

import { ElMessage } from 'element-plus';

const emit = defineEmits(['update:modelValue'])


const props = defineProps({
  modelValue: "",
})
const fileList = ref(props.modelValue != '' ? [{
  name: "image",
  url: props.modelValue
}] : [])
const imageUrl = ref("")

watch(() => props.modelValue, (val) => {
  if(val == ""){
    fileList.value = []
  }else{
    fileList.value = [{
      name: "image",
      url: val
    }]
  }
})

const handleUpload = (data) => {
  return new Promise(async (resolve, reject) => {
    const file = data.file;
    console.log(file);
    try {
      const result = await client.add({
        path: file.filepath || file.webkitRelativePath || file.name,
        content: file,
        size: file.size
      });
      console.log("result: ", result.cid.toString());
      imageUrl.value = formatImg(result.cid.toString());
      console.log("imageUrl: ", imageUrl.value);
      emit('update:modelValue', imageUrl.value)
      resolve(result.cid.toString());
    } catch (error) {
      reject(error);
    }
  })
}

const handleRemove = () => {
  imageUrl.value = "";
  emit('update:modelValue', "")
}


const dialogImageUrl = ref('')
const dialogVisible = ref(false)


const handlePictureCardPreview = (uploadFile) => {
  dialogImageUrl.value = uploadFile.url
  dialogVisible.value = true
}

const formatImg = (url) => {
  return `http://127.0.0.1:8080/ipfs/${url}`
}

</script>

<style>
.avatar-uploader .el-upload {
  border: 1px dashed var(--el-border-color);
  border-radius: 6px;
  cursor: pointer;
  position: relative;
  overflow: hidden;
  transition: var(--el-transition-duration-fast);
}

.avatar-uploader .el-upload:hover {
  border-color: var(--el-color-primary);
}

.el-icon.avatar-uploader-icon {
  font-size: 28px;
  color: #8c939d;
  width: 178px;
  height: 178px;
  text-align: center;
}
</style>