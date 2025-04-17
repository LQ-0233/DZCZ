<script setup>
import { ref } from 'vue';
import List from '/@/components/List.vue';
import api from '/@/api/api';
import { convertRes2Blob } from '/@/utils/download';

const tableColProps = [
    {
        "prop": "caseNumber",
        "label": "案件编号",
    },
    {
        "prop": "evidenceCreator",
        "label": "授权人",
    },
    {
        "prop": "authTime",
        "label": "授权时间",
        "type": "datetime",
      
    }
]

// 查询参数
const query = ref({
    page: 1,
    pageSize: 10,
    caseNumber: '',
    authorizerName: ''
})

// 表格数据
const tableData = ref([
])

// 详情弹窗控制
const detailDialogVisible = ref(false)
const currentDetail = ref(null)

// 获取授权文件列表
const find = () => {
    api.evidence.receivedAuthorizedList().then((res) => {
        tableData.value = res
    })
}

// 初始化加载
find()

// 重置查询
const reset = () => {
    query.value = {
        page: 1,
        pageSize: 10,
    }
    find()
}

// 查看详情
const viewDetail = async (row) => {

    const res = await api.evidence.view({
            caseNumber: row.caseNumber
        })
        currentDetail.value = res
        detailDialogVisible.value = true
}

// 下载文件
const downloadFile = async () => {
    api.common.download({
        ipfsLink: currentDetail.value.ipfsLink,
        fileName: currentDetail.value.fileName,
    }).then((res) => {
        convertRes2Blob(res)
    })
}


</script>

<template>
    <div class="min-h-full h-full">
        <List 
            :page="query" 
            :tableColProps="tableColProps" 
            :tableDate="tableData"
            @reset="reset" 
            @find="find"
            :operateWidth="120"
        >
 
            <!-- 操作按钮 -->
            <template #handler="{ scope }">
                <el-button 
                    type="primary" 
                    @click="viewDetail(scope.row)"
                >查看</el-button>
            </template>
        </List>

        <!-- 详情弹窗 -->
        <el-dialog
            v-model="detailDialogVisible"
            title="文件详情"
            width="600px"
        >
            <el-descriptions :column="1" border>
                <el-descriptions-item label="案件编号">
                    {{ currentDetail?.caseNumber }}
                </el-descriptions-item>
                <el-descriptions-item label="案件基本信息">
                    {{ currentDetail?.caseInfo }}
                </el-descriptions-item>
                <el-descriptions-item label="负责人1">
                    {{ currentDetail?.manager1 }}
                </el-descriptions-item>
                <el-descriptions-item label="负责人2">
                    {{ currentDetail?.manager2 }}
                </el-descriptions-item>
                <el-descriptions-item label="电子证据形式">
                    {{ currentDetail?.evidenceType }}
                </el-descriptions-item>
            </el-descriptions>
            
            <template #footer>
                <span class="dialog-footer">
                    <el-button 
                    type="primary"
                    @click="downloadFile(currentDetail?.id, currentDetail?.fileName)"
                >
                    下载文件
                </el-button>
                    <el-button @click="detailDialogVisible = false">关闭</el-button>
                </span>
            </template>
        </el-dialog>
    </div>
</template>

<style scoped>
.el-form-item {
    margin-right: 16px;
}
.mt-4 {
    margin-top: 1rem;
}
</style>