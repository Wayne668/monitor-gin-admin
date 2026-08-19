<template>
    <div style="padding: 24px">
        <!-- 页头 -->
        <div style="display: flex; align-items: center; margin-bottom: 20px">
            <a-button
                type="text"
                @click="$router.back()"
                style="font-size: 20px; padding: 0 8px">
                ←
            </a-button>
            <h2 style="margin: 0 12px">素材删除</h2>
        </div>

        <!-- 上半部分：账户选择 -->
        <a-card title="选择账户" style="margin-bottom: 16px">
            <div style="display: flex; flex-direction: column; gap: 12px">
                <a-form-item
                    label="客户公司"
                    style="margin-bottom: 0">
                    <a-select
                        v-model:value="selectedAgentId"
                        placeholder="请选择客户公司"
                        style="max-width: 400px"
                        :options="agentTokenOptions"
                        allow-clear
                        @change="handleAgentChange" />
                </a-form-item>
                <AccountTransfer
                    v-model="selectedAccountIds"
                    target="account"
                    :agent-id="selectedAgentId"
                    :disabled="!selectedAgentId" />
                <div style="text-align: right">
                    <a-button
                        type="primary"
                        :loading="materialLoading"
                        :disabled="selectedAccountIds.length === 0"
                        @click="handleLoadMaterials">
                        确定
                    </a-button>
                </div>
            </div>
        </a-card>

        <!-- 下半部分：素材列表 -->
        <a-card v-if="!showFailedOnly">
            <template #title>
                <span>未审核素材</span>
            </template>
            <template #extra>
                <a-button
                    type="primary"
                    danger
                    :disabled="selectedRowKeys.length === 0"
                    :loading="deleting"
                    @click="handleBatchDelete">
                    删除选中 ({{ selectedRowKeys.length }})
                </a-button>
            </template>
            <a-table
                :columns="columns"
                :data-source="tableData"
                :row-key="(record) => record.materialId"
                :loading="materialLoading"
                :pagination="false"
                :row-selection="{
                    selectedRowKeys: selectedRowKeys,
                    onChange: onSelectChange,
                }"
                size="small"
                bordered>
                <template #bodyCell="{ column, record }">
                    <template v-if="column.key === 'materialId'">
                        {{ record.materialId }}
                    </template>
                    <template v-if="column.key === 'materialName'">
                        <a-tooltip :title="record.materialName">
                            <span>{{ record.materialName || '-' }}</span>
                        </a-tooltip>
                    </template>
                </template>
            </a-table>
            <a-empty
                v-if="!materialLoading && tableData.length === 0"
                description="请先选择账户并点击确定" />
        </a-card>

        <!-- 删除失败记录 -->
        <a-card v-else title="删除失败素材">
            <template #extra>
                <a-button @click="handleBackToMaterials">
                    返回素材列表
                </a-button>
            </template>
            <a-table
                :columns="failedColumns"
                :data-source="failedData"
                :row-key="(record) => record.materialId"
                :pagination="false"
                size="small"
                bordered>
                <template #bodyCell="{ column, record }">
                    <template v-if="column.key === 'materialId'">
                        {{ record.materialId }}
                    </template>
                    <template v-if="column.key === 'materialName'">
                        <a-tooltip :title="record.materialName">
                            <span>{{ record.materialName || '-' }}</span>
                        </a-tooltip>
                    </template>
                    <template v-if="column.key === 'errorMsg'">
                        <a-tooltip :title="record.errorMsg">
                            <span style="color: #ff4d4f">{{ record.errorMsg }}</span>
                        </a-tooltip>
                    </template>
                </template>
            </a-table>
            <a-empty
                v-if="failedData.length === 0"
                description="没有失败的记录" />
        </a-card>
    </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { Modal, message } from 'ant-design-vue'
import { getHostAccountList } from '@/apis/modules/hostAccount'
import { getUnauditedMaterial, deleteUnauditedMaterial } from '@/apis/modules/delUnitMaterial'
import AccountTransfer from './components/AccountTransfer.vue'

const columns = [
    { title: '素材ID', dataIndex: 'materialId', key: 'materialId', width: 250 },
    { title: '素材名称', dataIndex: 'materialName', key: 'materialName', minWidth: 200, ellipsis: true },
    { title: '营销ID', dataIndex: 'promotionId', key: 'promotionId', width: 250 },
    { title: '广告主ID', dataIndex: 'advertiserId', key: 'advertiserId', width: 250 },
]

const failedColumns = [
    { title: '素材ID', dataIndex: 'materialId', key: 'materialId', width: 250 },
    { title: '素材名称', dataIndex: 'materialName', key: 'materialName', minWidth: 200, ellipsis: true },
    { title: '营销ID', dataIndex: 'promotionId', key: 'promotionId', width: 250 },
    { title: '广告主ID', dataIndex: 'advertiserId', key: 'advertiserId', width: 250 },
    { title: '错误信息', dataIndex: 'errorMsg', key: 'errorMsg', minWidth: 200, ellipsis: true },
]

const selectedAgentId = ref(undefined)
const agentTokenOptions = ref([])
const selectedAccountIds = ref([])
const tableData = ref([])
const selectedRowKeys = ref([])
const materialLoading = ref(false)
const deleting = ref(false)
const failedData = ref([])
const showFailedOnly = ref(false)

const loadAgentTokens = async () => {
    try {
        const res = await getHostAccountList({ pageSize: 100 })
        const seen = new Set()
        agentTokenOptions.value = (res.data || [])
            .filter((item) => {
                if (seen.has(item.agentId)) return false
                seen.add(item.agentId)
                return true
            })
            .map((item) => ({
                label: item.advCompanyName || item.advertiserName,
                value: item.agentId,
            }))
    } catch (e) {
        message.error('加载客户公司失败')
    }
}

const handleAgentChange = () => {
    selectedAccountIds.value = []
    showFailedOnly.value = false
    failedData.value = []
}

const handleLoadMaterials = async () => {
    if (selectedAccountIds.value.length === 0) {
        message.warning('请至少选择一个账户')
        return
    }
    materialLoading.value = true
    selectedRowKeys.value = []
    showFailedOnly.value = false
    failedData.value = []
    try {
        const res = await getUnauditedMaterial({ accountIds: selectedAccountIds.value.join(',') })
        tableData.value = res.data || []
        if (tableData.value.length === 0) {
            message.info('没有未审核的素材')
        } else {
            message.success(`已加载 ${tableData.value.length} 条素材`)
        }
    } catch (e) {
        message.error('加载素材失败')
    } finally {
        materialLoading.value = false
    }
}

const onSelectChange = (keys) => {
    selectedRowKeys.value = keys
}

const handleBatchDelete = () => {
    if (selectedRowKeys.value.length === 0) return
    const selectedRows = tableData.value.filter((r) => selectedRowKeys.value.includes(r.materialId))
    Modal.confirm({
        title: '确认删除',
        content: `确定删除 ${selectedRowKeys.value.length} 个素材？此操作不可撤销。`,
        okType: 'danger',
        okText: '确定删除',
        cancelText: '取消',
        onOk: async () => {
            deleting.value = true
            try {
                const params = {
                    accountIds: [...new Set(selectedRows.map((r) => r.advertiserId))],
                    materials: selectedRows.map((r) => ({
                        materialId: r.materialId,
                        promotionId: r.promotionId,
                        advertiserId: r.advertiserId,
                        materialName: r.materialName,
                    })),
                }
                const res = await deleteUnauditedMaterial(params)
                const failed = res.data?.failed || []
                const successCount = selectedRowKeys.value.length - failed.length
                if (successCount > 0) {
                    message.success(`成功删除 ${successCount} 个素材`)
                }
                selectedRowKeys.value = []
                if (failed.length > 0) {
                    failedData.value = failed
                    showFailedOnly.value = true
                    message.warning(`${failed.length} 个素材删除失败，请查看失败详情`)
                } else {
                    await handleLoadMaterials()
                }
            } catch (e) {
                message.error('删除失败')
            } finally {
                deleting.value = false
            }
        },
    })
}

const handleBackToMaterials = () => {
    showFailedOnly.value = false
    failedData.value = []
    handleLoadMaterials()
}

onMounted(() => {
    loadAgentTokens()
})
</script>