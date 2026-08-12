<template>
    <div style="padding: 24px">
        <a-card style="margin-bottom: 16px">
            <div style="display: flex; justify-content: space-between; align-items: center">
                <h2 style="margin: 0">AI托管规则</h2>
                <a-button
                    type="primary"
                    size="large"
                    @click="handleAdd">
                    ＋ 新建规则
                </a-button>
            </div>
        </a-card>

        <a-card>
            <a-table
                :columns="columns"
                :data-source="tableData"
                :pagination="pagination"
                :row-key="(record) => record.id"
                :loading="loading">
                <template #bodyCell="{ column, record }">
                    <template v-if="'date_range' === column.key">
                        <span>{{ record.trigger_start_date }} ~ {{ record.trigger_end_date }}</span>
                    </template>
                    <template v-if="'target_accounts' === column.key">
                        <span>{{ formatJsonArray(record.target_accounts) }}</span>
                    </template>
                    <template v-if="'target_obj' === column.key">
                        <span>{{ getTargetObj(record) }}</span>
                    </template>
                    <template v-if="'status' === column.key">
                        <a-switch
                            :checked="record.status === 1"
                            :loading="statusLoading[record.id]"
                            checked-children="启用"
                            un-checked-children="停用"
                            @change="(val) => handleToggle(record, val)" />
                    </template>
                    <template v-if="'action' === column.key">
                        <a-button
                            type="link"
                            size="small"
                            @click="handleLog(record)">
                            查看
                        </a-button>
                    </template>
                </template>
            </a-table>
        </a-card>
    </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { message } from 'ant-design-vue'
import { getHostRuleList, updateHostRuleStatus } from '@/apis/modules/hostRule'

const router = useRouter()

const columns = [
    { title: '规则ID', dataIndex: 'id', key: 'id', width: 80 },
    { title: '规则名称', dataIndex: 'rule_name', key: 'rule_name', width: 150 },
    { title: '托管目标', dataIndex: 'target', key: 'target', width: 100 },
    { title: '生效日期', key: 'date_range', width: 200 },
    { title: '创建人', dataIndex: 'user_name', key: 'user_name', width: 100 },
    { title: '创建时间', dataIndex: 'created_at', key: 'created_at', width: 170 },
    { title: '托管账户', dataIndex: 'target_accounts', key: 'target_accounts', width: 150 },
    { title: '托管对象', key: 'target_obj', width: 150 },
    { title: '状态', key: 'status', width: 80 },
    { title: '运行记录', key: 'action', width: 100, fixed: 'right' },
]

const loading = ref(false)
const tableData = ref([])
const statusLoading = reactive({})

const pagination = reactive({
    current: 1,
    pageSize: 10,
    total: 0,
    showTotal: (total) => `共 ${total} 条`,
    showSizeChanger: true,
    onChange: (page, pageSize) => {
        pagination.current = page
        pagination.pageSize = pageSize
        loadData()
    },
})

const loadData = async () => {
    loading.value = true
    try {
        const res = await getHostRuleList({
            current: pagination.current,
            pageSize: pagination.pageSize,
        })
        tableData.value = res.data || []
        pagination.total = res.total || 0
    } catch (e) {
        message.error('加载失败')
    } finally {
        loading.value = false
    }
}

const handleAdd = () => {
    router.push({ name: 'aiHostRuleForm' })
}

const handleToggle = async (row, checked) => {
    const newStatus = checked ? 1 : 2
    statusLoading[row.id] = true
    try {
        await updateHostRuleStatus(row.id, { status: newStatus })
        row.status = newStatus
        message.success(checked ? '已启用' : '已停用')
    } catch (e) {
        message.error('操作失败')
    } finally {
        statusLoading[row.id] = false
    }
}

const handleLog = (row) => {
    router.push({ name: 'aiHostTriggerRecord', query: { ruleId: row.id } })
}

const formatJsonArray = (str) => {
    if (!str) return '-'
    try {
        const arr = JSON.parse(str)
        return Array.isArray(arr) ? arr.join(', ') : str
    } catch {
        return str
    }
}

const getTargetObj = (row) => {
    switch (row.target) {
        case 'promotion':
            return formatJsonArray(row.target_promotion)
        case 'creative':
            return formatJsonArray(row.target_material)
        case 'project':
            return formatJsonArray(row.target_projects)
        default:
            return '-'
    }
}

onMounted(() => {
    loadData()
})
</script>