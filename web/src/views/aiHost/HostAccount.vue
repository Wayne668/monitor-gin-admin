<template>
    <div style="padding: 24px">
        <a-card style="margin-bottom: 16px">
            <div style="display: flex; justify-content: space-between; align-items: center">
                <a-space>
                    <a-select
                        v-model:value="searchParams.agentId"
                        placeholder="代理商搜索"
                        style="width: 240px"
                        allow-clear
                        show-search
                        :filter-option="filterAgentOption"
                        :options="agentTokenOptions"
                        @change="handleSearch" />
                    <a-input-search
                        v-model:value="searchParams.advertiserId"
                        placeholder="广告主ID搜索"
                        style="width: 200px"
                        allow-clear
                        @search="handleSearch" />
                    <a-input-search
                        v-model:value="searchParams.advertiserName"
                        placeholder="账户名称搜索"
                        style="width: 200px"
                        allow-clear
                        @search="handleSearch" />
                </a-space>
                <a-button
                    type="primary"
                    size="large"
                    @click="handleAdd">
                    ＋ 新增托管账户
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
                    <template v-if="'budget' === column.key">
                        <a-button
                            type="link"
                            size="small"
                            @click="handleBudget(record)">
                            编辑
                        </a-button>
                    </template>
                    <template v-if="'budgetRecord' === column.key">
                        <a-button
                            type="link"
                            size="small"
                            @click="handleBudgetRecord(record)">
                            查看
                        </a-button>
                    </template>
                    <template v-if="'action' === column.key">
                        <a-button
                            type="link"
                            size="small"
                            @click="handleEdit(record)">
                            编辑
                        </a-button>
                        <a-button
                            type="link"
                            danger
                            size="small"
                            @click="handleDelete(record)">
                            删除
                        </a-button>
                    </template>
                </template>
            </a-table>
        </a-card>

        <a-modal
            v-model:open="modalVisible"
            :title="isEdit ? '编辑托管账户' : '新增托管账户'"
            :confirm-loading="submitting"
            width="520px"
            @ok="handleSubmit">
            <a-form
                ref="formRef"
                :model="form"
                :rules="rules"
                :label-col="{ span: 6 }"
                :wrapper-col="{ span: 16 }">
                <a-form-item
                    label="代理商ID"
                    name="agentId">
                    <a-input
                        v-model:value="form.agentId"
                        placeholder="请输入代理商ID" />
                </a-form-item>
                <a-form-item
                    label="广告主ID"
                    name="advertiserId">
                    <a-input
                        v-model:value="form.advertiserId"
                        placeholder="请输入广告主ID" />
                </a-form-item>
                <a-form-item
                    label="广告主名称"
                    name="advertiserName">
                    <a-input
                        v-model:value="form.advertiserName"
                        placeholder="请输入广告主名称" />
                </a-form-item>
            </a-form>
        </a-modal>

        <a-modal
            v-model:open="budgetModalVisible"
            title="账户预算"
            :confirm-loading="budgetSubmitting"
            width="520px"
            @ok="handleBudgetSubmit">
            <a-form
                ref="budgetFormRef"
                :model="budgetForm"
                :rules="budgetRules"
                :label-col="{ span: 6 }"
                :wrapper-col="{ span: 16 }">
                <a-form-item label="生效方式" name="effectType">
                    <a-radio-group v-model:value="budgetForm.effectType">
                        <a-radio value="immediate">立即生效</a-radio>
                        <a-radio value="nextDay">次日生效</a-radio>
                    </a-radio-group>
                </a-form-item>
                <template v-if="budgetForm.effectType === 'immediate'">
                    <a-form-item label="预算模式" name="budgetMode">
                        <a-select
                            v-model:value="budgetForm.budgetMode"
                            placeholder="请选择预算模式">
                            <a-select-option value="BUDGET_MODE_INFINITE">不限预算</a-select-option>
                            <a-select-option value="BUDGET_MODE_DAY">指定预算</a-select-option>
                        </a-select>
                    </a-form-item>
                    <a-form-item label="预算金额" name="budget">
                        <a-input
                            v-model:value="budgetForm.budget"
                            placeholder="请输入预算金额（元）"
                            :disabled="budgetForm.budgetMode === 'BUDGET_MODE_INFINITE'" />
                    </a-form-item>
                </template>
                <template v-if="budgetForm.effectType === 'nextDay'">
                    <a-form-item label="预算模式">
                        <a-input value="指定预算（BUDGET_MODE_DAY）" disabled />
                    </a-form-item>
                    <a-form-item label="预算金额">
                        <a-input value="1000" disabled />
                    </a-form-item>
                </template>
            </a-form>
        </a-modal>

        <a-modal
            v-model:open="recordModalVisible"
            title="预算修改记录"
            :footer="null"
            width="700px">
            <a-table
                :columns="recordColumns"
                :data-source="recordData"
                :loading="recordLoading"
                :pagination="false"
                row-key="id"
                size="small" />
        </a-modal>
    </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { Modal, message } from 'ant-design-vue'
import { getHostAccountList, getHostAccount, createHostAccount, updateHostAccount, delHostAccount } from '@/apis/modules/hostAccount'
import { getAgentTokenList } from '@/apis/modules/agentToken'
import { updateAdvertiserBudget, scheduleAdvertiserBudget, getBudgetRecords } from '@/apis/modules/advertiserBudget'

const columns = [
    { title: 'ID', dataIndex: 'id', key: 'id', width: 70 },
    { title: '代理商ID', dataIndex: 'agentId', key: 'agentId', width: 180 },
    { title: '广告主ID', dataIndex: 'advertiserId', key: 'advertiserId', width: 180 },
    { title: '广告主名称', dataIndex: 'advertiserName', key: 'advertiserName', ellipsis: true },
    // { title: '账户预算', key: 'budget', width: 100 },
    // { title: '预算记录', key: 'budgetRecord', width: 100 },
    { title: '操作', key: 'action', width: 140, fixed: 'right' },
]

const recordColumns = [
    { title: 'ID', dataIndex: 'id', key: 'id', width: 60 },
    { title: '预算金额', dataIndex: 'budget', key: 'budget', width: 100 },
    { title: '预算模式', dataIndex: 'budgetMod', key: 'budgetMod', width: 100 },
    { title: '是否生效', dataIndex: 'isSet', key: 'isSet', width: 80 },
    { title: '错误信息', dataIndex: 'errMsg', key: 'errMsg', ellipsis: true },
    { title: '创建时间', dataIndex: 'createdAt', key: 'createdAt', width: 170 },
]

const loading = ref(false)
const tableData = ref([])
const searchParams = reactive({ agentId: '', advertiserId: '', advertiserName: '' })

const agentTokenOptions = ref([])

const filterAgentOption = (input, option) => {
    return option.label.toLowerCase().indexOf(input.toLowerCase()) >= 0
}

const loadAgentTokens = async () => {
    try {
        const res = await getAgentTokenList({ pageSize: 100 })
        agentTokenOptions.value = (res.data || []).map((item) => ({
            label: `${item.accountName} (${item.accountId})`,
            value: item.accountId,
        }))
    } catch (e) {
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

const modalVisible = ref(false)
const submitting = ref(false)
const isEdit = ref(false)
const editId = ref(null)
const formRef = ref()

const defaultForm = {
    agentId: '',
    advertiserId: '',
    advertiserName: '',
}
const form = reactive({ ...defaultForm })

const rules = {
    agentId: [{ required: true, message: '请输入代理商ID', trigger: 'blur' }],
    advertiserId: [{ required: true, message: '请输入广告主ID', trigger: 'blur' }],
    advertiserName: [{ required: true, message: '请输入广告主名称', trigger: 'blur' }],
}

const loadData = async () => {
    loading.value = true
    try {
        const res = await getHostAccountList({
            current: pagination.current,
            pageSize: pagination.pageSize,
            agentId: searchParams.agentId || undefined,
            advertiserId: searchParams.advertiserId || undefined,
            advertiserName: searchParams.advertiserName || undefined,
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

const resetForm = () => {
    Object.assign(form, defaultForm)
}

const handleAdd = () => {
    isEdit.value = false
    editId.value = null
    resetForm()
    modalVisible.value = true
}

const handleEdit = async (row) => {
    isEdit.value = true
    editId.value = row.id
    try {
        const res = await getHostAccount(row.id)
        const data = res.data
        form.agentId = String(data.agentId || '')
        form.advertiserId = String(data.advertiserId || '')
        form.advertiserName = data.advertiserName || ''
        modalVisible.value = true
    } catch (e) {
        message.error('加载失败')
    }
}

const handleSubmit = async () => {
    try {
        await formRef.value.validate()
    } catch {
        return
    }
    submitting.value = true
    try {
        const payload = {
            agentId: Number(form.agentId),
            advertiserId: Number(form.advertiserId),
            advertiserName: form.advertiserName,
        }
        if (isEdit.value) {
            await updateHostAccount(editId.value, payload)
            message.success('更新成功')
        } else {
            await createHostAccount(payload)
            message.success('创建成功')
        }
        modalVisible.value = false
        loadData()
    } catch (e) {
        message.error(isEdit.value ? '更新失败' : '创建失败')
    } finally {
        submitting.value = false
    }
}

const handleDelete = (row) => {
    Modal.confirm({
        title: '确认删除',
        content: `确定删除托管账户「${row.advertiserName || row.advertiserId}」？`,
        okType: 'danger',
        okText: '确定',
        cancelText: '取消',
        onOk: async () => {
            try {
                await delHostAccount(row.id)
                message.success('删除成功')
                loadData()
            } catch (e) {
                message.error('删除失败')
            }
        },
    })
}

// 预算弹窗
const budgetModalVisible = ref(false)
const budgetSubmitting = ref(false)
const budgetFormRef = ref()
const currentBudgetRecord = ref(null)

const defaultBudgetForm = {
    effectType: 'immediate',
    budgetMode: '',
    budget: '',
}
const budgetForm = reactive({ ...defaultBudgetForm })

const budgetRules = {
    effectType: [{ required: true, message: '请选择生效方式', trigger: 'change' }],
    budgetMode: [{ required: true, message: '请选择预算模式', trigger: 'change' }],
}

const resetBudgetForm = () => {
    Object.assign(budgetForm, defaultBudgetForm)
}

const handleBudget = (record) => {
    currentBudgetRecord.value = record
    resetBudgetForm()
    budgetModalVisible.value = true
}

const handleBudgetSubmit = async () => {
    // 动态校验
    if (budgetForm.effectType === 'immediate') {
        if (!budgetForm.budgetMode) {
            message.error('请选择预算模式')
            return
        }
        if (budgetForm.budgetMode === 'BUDGET_MODE_DAY' && !budgetForm.budget) {
            message.error('请输入预算金额')
            return
        }
    }
    budgetSubmitting.value = true
    try {
        const record = currentBudgetRecord.value
        if (budgetForm.effectType === 'immediate') {
            await updateAdvertiserBudget({
                accountId: record.agentId,
                advertiserId: record.advertiserId,
                budgetMode: budgetForm.budgetMode,
                budget: budgetForm.budgetMode === 'BUDGET_MODE_DAY' ? Number(budgetForm.budget) : 0,
            })
            message.success('预算更新成功')
        } else {
            await scheduleAdvertiserBudget({
                advertiserId: record.advertiserId,
                budget: 1000,
                budgetMod: 'nextDay',
            })
            message.success('已提交次日生效预算')
        }
        budgetModalVisible.value = false
    } catch (e) {
        message.error(budgetForm.effectType === 'immediate' ? '预算更新失败' : '提交失败')
    } finally {
        budgetSubmitting.value = false
    }
}

// 预算记录弹窗
const recordModalVisible = ref(false)
const recordLoading = ref(false)
const recordData = ref([])

const handleBudgetRecord = async (record) => {
    recordModalVisible.value = true
    recordLoading.value = true
    try {
        const res = await getBudgetRecords({ advertiserId: record.advertiserId })
        recordData.value = res.data || []
    } catch (e) {
        message.error('加载预算记录失败')
    } finally {
        recordLoading.value = false
    }
}

onMounted(() => {
    loadAgentTokens()
    loadData()
})
</script>