<script setup>
import { ref } from 'vue';
import { ElMessageBox } from 'element-plus';
import List from '/@/components/List.vue';
import api from '/@/api/api';

const tableColProps = [
    {
        "prop": "caseNumber",
        "label": "案件编号",
        "width": "150"
    },
    {
        "prop": "authTime",
        "label": "授权时间",
        "type": "date",
    },
    {
        "prop": "authorizedUser",
        "label": "被授权人",
    },
    {
        "prop": "status",
        "label": "状态",
        "type": "map",
        "map": {
            "authorized": '已授权',
            "canceled": '已取消'
        }
    },
    {
        "prop": "viewCount",
        "label": "查看次数",
    }
]

// 取消授权
const cancelAuth = (row) => {
    ElMessageBox.confirm(
        '确定要取消该用户的授权吗？',
        '取消授权确认',
        {
            confirmButtonText: '确定',
            cancelButtonText: '取消',
            type: 'warning',
        }
    )
    .then(() => {
        api.evidence.cancelAuthorize({
            id: row.id
        }).then(() => {
            ElMessage.success('取消授权成功')
            find()
        })
    })

}

// 查询参数
const query = ref({
    page: 1,
    pageSize: 10,
    caseNumber: '',
    status: '',
    authorizedUser: ''
})

// 表格数据
const tableData = ref([

])

// 获取授权记录列表
const find = () => {
    api.evidence.authorizedList().then((res) => {
        tableData.value = res
    })
}

// 重置查询
const reset = () => {
    query.value = {
        page: 1,
        pageSize: 10,

    }
    find()
}

// 初始化加载
find()
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
                    type="danger" 
                    @click="cancelAuth(scope.row)"
                    v-if="scope.row.status === 'authorized'"
                >取消授权</el-button>
            </template>
        </List>
    </div>
</template>

<style scoped>
.el-form-item {
    margin-right: 16px;
}
</style> 