<script setup>
import { ElMessage, ElMessageBox } from 'element-plus';
import _ from 'lodash';
import moment from 'moment';
import List from '/@/components/List.vue';
import api from '/@/api/api';
import { getCurrentInstance, ref } from 'vue';

const { proxy } = getCurrentInstance()

const defaultQuery = {
    username: "",
    nickname: "",
    role: "",
    status: "",
    page: 1,
    size: 10,
    total: 0

}

const query = ref(_.cloneDeep(defaultQuery))

const reset = () => {
    query.value = _.cloneDeep(defaultQuery)
}

const tableColProps = [
    // {
    //   "prop": "id",
    //   "label": "ID"
    // },
    {
        "prop": "username",
        "label": "用户名"
    },
    {
        "prop": "nickname",
        "label": "名称"
    },
    {
        "prop": "role",
        "label": "角色",
        "type": "map",
        "map": {
            "3": "查阅人员",
            "2": "录入人员"
        }
    },
    {
        "prop": "status",
        "label": "状态",
        "type": "map",
        "map": {
            "enable": "启用",
            "disable": "禁用"
        }
    },
    {
        "prop": "registerTime",
        "label": "注册时间",
        "type": "datetime"
    },
]

const formProps = [
    {
        "prop": "role",
        "label": "角色",
        "type": "select",
        "options": [
            {
                "label": "查阅人员",
                "value": "3"
            },
            {
                "label": "录入人员",
                "value": "2"
            }
        ]
    },
    {
        "prop": "status",
        "label": "状态",
        "type": "switch",
        "default": "enable",
        "activeText": "启用",
        "inactiveText": "禁用",
        "activeValue": "enable",
        "inactiveValue": "disable"
    },
]

const formRules = {

}

const tableData = ref([
])

const find = () => {
    console.log("find query: ", query.value);
    api.user.list({
    }).then((res) => {
        console.log("find res: ", res);
        tableData.value = res
    })
}

find()


const update = (item) => {
    console.log(item);
    api.user.update(item).then((res) => {

        ElMessage.success("更新成功")
        find()

    })
}

</script>
<template>
    <div class="min-h-full h-full">
        <List :page="query" :tableColProps="tableColProps" :tableDate="tableData" :formProps="formProps"
            :showHeader="false" :formRules="formRules" @reset="reset" @find="find" :showOperate="true" :showUpdate="true"
            @update="update">

        </List>
    </div>
</template>