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
                <AccountTransfer
                    :options="accountOptions"
                    v-model="selectedAccountIds"
                    target="account" />
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
        <a-card title="未审核素材">
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
    </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { Modal, message } from 'ant-design-vue'
import { getEnabledAccountList } from '@/apis/modules/accountInfo'
import { getUnauditedMaterial, deleteUnauditedMaterial } from '@/apis/modules/hostRule'
import AccountTransfer from './components/AccountTransfer.vue'

const columns = [
    { title: '素材ID', dataIndex: 'materialId', key: 'materialId', width: 120 },
    { title: '素材名称', dataIndex: 'materialName', key: 'materialName', minWidth: 200, ellipsis: true },
    { title: '营销ID', dataIndex: 'promotionId', key: 'promotionId', width: 120 },
    { title: '广告主ID', dataIndex: 'advertiserId', key: 'advertiserId', width: 120 },
]

const accountOptions = ref([])
const selectedAccountIds = ref([])
const tableData = ref([])
const selectedRowKeys = ref([])
const materialLoading = ref(false)
const deleting = ref(false)

const loadAccounts = async () => {
    try {
        const res = await getEnabledAccountList({ fields: 'advertiser_id,advertiser_name' })
        accountOptions.value = (res.data || []).map((item) => ({
            id: item.advertiserId,
            name: item.advertiserName || String(item.advertiserId),
        }))
    } catch (e) {
        message.error('加载账户列表失败')
    }
}

const handleLoadMaterials = async () => {
    if (selectedAccountIds.value.length === 0) {
        message.warning('请至少选择一个账户')
        return
    }
    materialLoading.value = true
    selectedRowKeys.value = []
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
                    })),
                }
                await deleteUnauditedMaterial(params)
                message.success(`成功删除 ${selectedRowKeys.value.length} 个素材`)
                selectedRowKeys.value = []
                await handleLoadMaterials()
            } catch (e) {
                message.error('删除失败')
            } finally {
                deleting.value = false
            }
        },
    })
}

onMounted(() => {
    loadAccounts()
})
</script>