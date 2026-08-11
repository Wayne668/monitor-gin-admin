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
    { title: '规则ID', dataIndex: 'id', key: 'id', width: 100 },
    { title: '规则名称', dataIndex: 'rule_name', key: 'rule_name' },
    { title: '托管目标', dataIndex: 'target', key: 'target', width: 100 },
    { title: '状态', key: 'status', width: 100 },
    { title: '运行记录', key: 'action', width: 120, fixed: 'right' },
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
    router.push({ name: 'aiHostRuleLog', params: { id: row.id } })
}

onMounted(() => {
    loadData()
})
</script>