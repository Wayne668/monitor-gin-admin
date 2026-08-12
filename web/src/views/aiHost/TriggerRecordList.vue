<template>
    <div style="padding: 24px">
        <div style="display: flex; align-items: center; margin-bottom: 20px">
            <a-button
                type="text"
                @click="goBack"
                style="font-size: 20px; padding: 0 8px">
                ←
            </a-button>
            <h2 style="margin: 0 12px">托管触发记录</h2>
        </div>

        <a-card>
            <a-table
                :columns="columns"
                :data-source="tableData"
                :pagination="pagination"
                :row-key="(record) => record.id"
                :loading="loading"
                :scroll="{ x: 1200 }">
                <template #bodyCell="{ column, record }">
                    <template v-if="'execute_status' === column.key">
                        <a-tag :color="statusColor(record.execute_status)">
                            {{ statusText(record.execute_status) }}
                        </a-tag>
                    </template>
                </template>
            </a-table>
        </a-card>
    </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { message } from 'ant-design-vue'
import { getHostTriggerRecordList } from '@/apis/modules/hostTriggerRecord'

const router = useRouter()
const route = useRoute()

const columns = [
    { title: '规则ID', dataIndex: 'rule_id', key: 'rule_id', width: 80 },
    { title: '广告主ID', dataIndex: 'advertiser_id', key: 'advertiser_id', width: 120 },
    { title: '推广ID', dataIndex: 'promotion_id', key: 'promotion_id', width: 120 },
    { title: '素材ID', dataIndex: 'material_id', key: 'material_id', width: 120 },
    { title: '目标', dataIndex: 'target', key: 'target', width: 80 },
    { title: '执行动作', dataIndex: 'execute_action', key: 'execute_action', width: 100 },
    { title: '执行状态', dataIndex: 'execute_status', key: 'execute_status', width: 100 },
    { title: '执行情况', dataIndex: 'execute_msg', key: 'execute_msg' },
    { title: '命中原因', dataIndex: 'trigger_reason', key: 'trigger_reason' },
    { title: '触发时间', dataIndex: 'created_at', key: 'created_at', width: 170 },
]

const loading = ref(false)
const tableData = ref([])

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

const statusColor = (status) => {
    const map = { pending: 'processing', succeed: 'success', failed: 'error' }
    return map[status] || 'default'
}

const statusText = (status) => {
    const map = { pending: '待执行', succeed: '成功', failed: '失败' }
    return map[status] || status
}

const loadData = async () => {
    loading.value = true
    try {
        const params = {
            current: pagination.current,
            pageSize: pagination.pageSize,
        }
        // 从路由参数获取 ruleId 筛选
        const ruleId = route.query.ruleId
        if (ruleId) {
            params.ruleId = Number(ruleId)
        }
        const res = await getHostTriggerRecordList(params)
        tableData.value = res.data || []
        pagination.total = res.total || 0
    } catch (e) {
        message.error('加载失败')
    } finally {
        loading.value = false
    }
}

const goBack = () => {
    if (window.history.length > 1) {
        router.back()
    } else {
        router.push({ name: 'aiHost' })
    }
}

onMounted(() => {
    loadData()
})
</script>