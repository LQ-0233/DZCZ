<script setup>
import { ref } from 'vue';
import List from '/@/components/List.vue';
import api from '/@/api/api';


const tableColProps = [
    {
        "prop": "caseNumber",
        "label": "案件编号",
    },
    {
        "prop": "viewUser",
        "label": "查看人",
    
    },
    {
        "prop": "viewTime",
        "label": "查看时间",
        "type": "datetime",
    }
]

// 查询参数
const query = ref({
    page: 1,
    pageSize: 10,
    total: 0
})

// 表格数据
const tableData = ref([
])

// 获取查看记录列表
const find = () => {


    api.evidence.viewRecordList().then((res) => {
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


</script>

<template>
    <div class="min-h-full h-full">
        <List 
            :page="query" 
            :tableColProps="tableColProps" 
            :tableDate="tableData"
            @reset="reset" 
            @find="find"
            :showOperate="false"
        >
        </List>
    </div>
</template>

<style scoped>
.el-form-item {
    margin-right: 16px;
}
</style>