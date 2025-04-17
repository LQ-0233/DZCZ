<script setup>
import { ElMessage } from 'element-plus';
import _ from 'lodash';
import moment from 'moment';
import List from '/@/components/List.vue';
import api from '/@/api/api';
import { getCurrentInstance, ref } from 'vue';
import { convertRes2Blob } from '/@/utils/download';

const query = ref({
    page: 1,
    pageSize: 10,
    total: 20
})

const tableColProps = [
    {
        "prop": "caseNumber",
        "label": "案件编号",
        "width": "150"
    },
    {
        "prop": "caseInfo",
        "label": "案件基本信息",
        "width": "200",
        "type": "overflowTooltip",
    },
    {
        "prop": "manager1",
        "label": "负责人1",
   
    },
    {
        "prop": "manager2",
        "label": "负责人2",
    },
    {
        "prop": "evidenceType",
        "label": "电子证据形式",
    },
    {
        "prop": "fileHash",
        "label": "文件Hash",

    },
    {
        "prop": "viewCount",
        "label": "查看次数",
    },
    {
        "prop": "createTime",
        "label": "创建时间",
        "type": "datetime",
        "width": "160"
    }
]

const tableData = ref([

])

const find = () => {

    api.evidence.userEvidence({
    }).then((res) => {
        tableData.value = res
    })
}
find()

// 添加授权对话框的状态控制
const authDialogVisible = ref(false)
const selectedUsers = ref([])
const currentRow = ref(null)

// 添加用户列表数据
const userList = ref([])

// 获取用户列表
const getUserList = () => {
    api.evidence.authorizeUserList().then((res) => {
        userList.value = res
    })
}

// 打开授权对话框
const openAuthDialog = (row) => {
    currentRow.value = row
    authDialogVisible.value = true
    getUserList()
}

// 提交授权
const submitAuth = () => {
    if (!selectedUsers.value.length) {
        ElMessage.warning('请选择要授权的用户')
        return
    }
    
    api.evidence.authorize({
        caseNumber: currentRow.value.caseNumber,
        authorizedUser: selectedUsers.value
    }).then((res) => {
        ElMessage.success('授权成功')
        authDialogVisible.value = false
        selectedUsers.value = []
        currentRow.value = null
    })
}

const downloadFile = (row) => {
    console.log(row)
    api.common.download({
        ipfsLink: row.ipfsLink,
        fileName: row.fileName,
    }).then((res) => {
        convertRes2Blob(res)
    })
}


</script>

<template>
    <div class="min-h-full h-full">
        <List :page="query" 
              :tableColProps="tableColProps" 
              :tableDate="tableData" 
              :showHeader="false" 
              @reset="reset" 
              @find="find"
              :operateWidth="300"
              >
              <template #handler="{ scope }">
                <el-button type="primary" @click="openAuthDialog(scope.row)">授权</el-button>
                <el-button type="warning" @click="downloadFile(scope.row)">下载</el-button>
           
              </template>
        </List>

        <!-- 授权对话框 -->
        <el-dialog
            v-model="authDialogVisible"
            title="文件授权"
            width="500px"
        >
            <el-form>
                <el-form-item label="选择用户">
                    <el-select
                        v-model="selectedUsers"
                        placeholder="请选择要授权的用户"
                        style="width: 100%"
                    >
                        <el-option
                            v-for="user in userList"
                            :key="user"
                            :label="user"
                            :value="user"
                        />
                    </el-select>
                </el-form-item>
            </el-form>
            <template #footer>
                <span class="dialog-footer">
                    <el-button @click="authDialogVisible = false">取消</el-button>
                    <el-button type="primary" @click="submitAuth">确认授权</el-button>
                </span>
            </template>
        </el-dialog>
    </div>
</template>

<style scoped>
.mx-2 {
    margin: 0 8px;
}
</style>
