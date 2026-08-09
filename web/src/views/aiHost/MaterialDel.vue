<template>
    <div style="padding: 24px">
        <a-card style="margin-bottom: 16px">
            <div style="display: flex; justify-content: space-between; align-items: center">
                <a-space>
                    <a-input-search
                        v-model:value="searchParams.materialName"
                        placeholder="素材名称搜索"
                        style="width: 200px"
                        allow-clear
                        @search="handleSearch" />
                    <a-select
                        v-model:value="searchParams.isDeleted"
                        placeholder="删除状态"
                        style="width: 140px"
                        allow-clear
                        @change="handleSearch">
                        <a-select-option value="pending">待处理</a-select-option>
                        <a-select-option value="deleted">已删除</a-select-option>
                        <a-select-option value="failed">失败</a-select-option>
                    </a-select>
                    <a-select
                        v-model:value="retryAccountId"
                        placeholder="选择代理商账户"
                        style="width: 200px"
                        :options="agentTokenOptions" />
                    <a-button
                        type="primary"
                        :loading="retryLoading"
                        :disabled="!retryAccountId"
                        @click="handleRetryFailed">
                        重试失败删除
                    </a-button>
                </a-space>
                <a-button
                    type="primary"
                    size="large"
                    @click="handleOpenForm">
                    ＋ 素材删除
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
                    <template v-if="column.key === 'isDeleted'">
                        <a-tag :color="isDeletedColor(record.isDeleted)">
                            {{ isDeletedText(record.isDeleted) }}
                        </a-tag>
                    </template>
                    <template v-if="column.key === 'errorMsg'">
                        <a-tooltip :title="record.errorMsg">
                            <span>{{ record.errorMsg || '-' }}</span>
                        </a-tooltip>
                    </template>
                    <template v-if="column.key === 'createdAt'">
                        {{ record.createdAt || '-' }}
                    </template>
                    <template v-if="column.key === 'updatedAt'">
                        {{ record.updatedAt || '-' }}
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
import { getDeleteUnauditedMaterialList, retryFailedDelete } from '@/apis/modules/deleteUnauditedMaterial'
import { getAgentTokenList } from '@/apis/modules/agentToken'

const router = useRouter()

const columns = [
    { title: 'ID', dataIndex: 'id', key: 'id', width: 70 },
    { title: '素材ID', dataIndex: 'materialId', key: 'materialId', width: 120 },
    { title: '素材名称', dataIndex: 'materialName', key: 'materialName', minWidth: 200, ellipsis: true },
    { title: '广告主ID', dataIndex: 'advertiserId', key: 'advertiserId', width: 120 },
    { title: '营销ID', dataIndex: 'promotionId', key: 'promotionId', width: 120 },
    { title: '删除状态', dataIndex: 'isDeleted', key: 'isDeleted', width: 100 },
    { title: '错误信息', dataIndex: 'errorMsg', key: 'errorMsg', minWidth: 200, ellipsis: true },
    { title: '重试次数', dataIndex: 'retryTimes', key: 'retryTimes', width: 90 },
    { title: '创建时间', dataIndex: 'createdAt', key: 'createdAt', width: 170 },
    { title: '更新时间', dataIndex: 'updatedAt', key: 'updatedAt', width: 170 },
]

const loading = ref(false)
const tableData = ref([])
const searchParams = reactive({ materialName: '', isDeleted: undefined })

const retryAccountId = ref(undefined)
const retryLoading = ref(false)
const agentTokenOptions = ref([])

const loadAgentTokenOptions = async () => {
    try {
        const res = await getAgentTokenList({ current: 1, pageSize: 100 })
        agentTokenOptions.value = (res.data || []).map((t) => ({
            label: `${t.accountName} (${t.accountId})`,
            value: t.accountId,
        }))
    } catch {
        // ignore
    }
}

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
        const res = await getDeleteUnauditedMaterialList({
            current: pagination.current,
            pageSize: pagination.pageSize,
            materialName: searchParams.materialName || undefined,
            isDeleted: searchParams.isDeleted || undefined,
        })
        tableData.value = res.data || []
        pagination.total = res.total || 0
    } catch (e) {
        message.error('加载失败')
    } finally {
        loading.value = false
    }
}

const handleSearch = () => {
    pagination.current = 1
    loadData()
}

const handleOpenForm = () => {
    router.push('/tools/materialDel/form')
}

const handleRetryFailed = async () => {
    if (!retryAccountId.value) {
        message.warning('请先选择代理商账户')
        return
    }
    retryLoading.value = true
    try {
        await retryFailedDelete(retryAccountId.value)
        message.success('重试删除完成')
        loadData()
    } catch (e) {
        message.error(e?.message || '重试失败')
    } finally {
        retryLoading.value = false
    }
}

const isDeletedColor = (status) => {
    const map = { pending: 'orange', deleted: 'green', failed: 'red' }
    return map[status] || 'default'
}

const isDeletedText = (status) => {
    const map = { pending: '待处理', deleted: '已删除', failed: '失败' }
    return map[status] || status || '-'
}

onMounted(() => {
    loadData()
    loadAgentTokenOptions()
})
</script>