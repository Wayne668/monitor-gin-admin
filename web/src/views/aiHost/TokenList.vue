<template>
    <div style="padding: 24px">
        <a-card style="margin-bottom: 16px">
            <div style="display: flex; justify-content: space-between; align-items: center">
                <a-space>
                    <a-input-search
                        v-model:value="searchParams.accountName"
                        placeholder="账号名搜索"
                        style="width: 200px"
                        allow-clear
                        @search="handleSearch" />
                    <a-input-search
                        v-model:value="searchParams.accountId"
                        placeholder="账号ID搜索"
                        style="width: 200px"
                        allow-clear
                        @search="handleSearch" />
                </a-space>
                <a-button
                    type="primary"
                    size="large"
                    @click="handleAdd">
                    ＋ 新增Token
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
                    <template v-if="'authStatus' === column.key">
                        <a-tag :color="record.authStatus === 'success' ? 'green' : 'orange'">
                            {{ record.authStatus || '-' }}
                        </a-tag>
                    </template>
                    <template v-if="'tokenTime' === column.key">
                        {{ record.tokenTime ? formatTime(record.tokenTime) : '-' }}
                    </template>
                    <template v-if="'accessToken' === column.key">
                        <a-tooltip :title="record.accessToken">
                            <span>{{ maskToken(record.accessToken) }}</span>
                        </a-tooltip>
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
            :title="isEdit ? '编辑Token' : '新增Token'"
            :confirm-loading="submitting"
            width="640px"
            @ok="handleSubmit">
            <a-form
                ref="formRef"
                :model="form"
                :rules="rules"
                :label-col="{ span: 6 }"
                :wrapper-col="{ span: 16 }">
                <a-form-item
                    label="账号名"
                    name="accountName">
                    <a-input
                        v-model:value="form.accountName"
                        placeholder="请输入账号名" />
                </a-form-item>
                <a-form-item
                    label="账号ID"
                    name="accountId">
                    <a-input
                        v-model:value="form.accountId"
                        placeholder="请输入账号ID" />
                </a-form-item>
                <a-form-item
                    label="授权状态"
                    name="authStatus">
                    <a-input
                        v-model:value="form.authStatus"
                        placeholder="请输入授权状态" />
                </a-form-item>
                <a-form-item
                    label="AccessToken"
                    name="accessToken">
                    <a-textarea
                        v-model:value="form.accessToken"
                        placeholder="请输入AccessToken"
                        :rows="2" />
                </a-form-item>
                <a-form-item
                    label="RefreshToken"
                    name="refreshToken">
                    <a-textarea
                        v-model:value="form.refreshToken"
                        placeholder="请输入RefreshToken"
                        :rows="2" />
                </a-form-item>
                <a-form-item
                    label="App ID"
                    name="appId">
                    <a-input
                        v-model:value="form.appId"
                        placeholder="请输入App ID" />
                </a-form-item>
                <a-form-item
                    label="App Secret"
                    name="appSecret">
                    <a-input
                        v-model:value="form.appSecret"
                        placeholder="请输入App Secret" />
                </a-form-item>
                <a-form-item
                    label="App名称"
                    name="appName">
                    <a-input
                        v-model:value="form.appName"
                        placeholder="请输入App名称" />
                </a-form-item>
                <a-form-item
                    label="备注"
                    name="remarks">
                    <a-textarea
                        v-model:value="form.remarks"
                        placeholder="请输入备注"
                        :rows="2" />
                </a-form-item>
            </a-form>
        </a-modal>
    </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { Modal, message } from 'ant-design-vue'
import { getAgentTokenList, getAgentToken, createAgentToken, updateAgentToken, delAgentToken } from '@/apis/modules/agentToken'

const columns = [
    { title: 'ID', dataIndex: 'id', key: 'id', width: 70 },
    { title: '账号名', dataIndex: 'accountName', key: 'accountName', width: 150 },
    { title: '账号ID', dataIndex: 'accountId', key: 'accountId', width: 120 },
    { title: '授权状态', dataIndex: 'authStatus', key: 'authStatus', width: 100 },
    { title: 'AccessToken', dataIndex: 'accessToken', key: 'accessToken', width: 200, ellipsis: true },
    { title: 'App名称', dataIndex: 'appName', key: 'appName', width: 120 },
    { title: '更新时间', dataIndex: 'tokenTime', key: 'tokenTime', width: 160 },
    { title: '备注', dataIndex: 'remarks', key: 'remarks', ellipsis: true },
    { title: '操作', key: 'action', width: 140, fixed: 'right' },
]

const loading = ref(false)
const tableData = ref([])
const searchParams = reactive({ accountName: '', accountId: '' })

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
    accountName: '',
    accountId: '',
    authStatus: '',
    accessToken: '',
    refreshToken: '',
    tokenTime: 0,
    remarks: '',
    appId: '',
    appSecret: '',
    appName: '',
}
const form = reactive({ ...defaultForm })

const rules = {
    accountId: [{ required: true, message: '请输入账号ID', trigger: 'blur' }],
}

const loadData = async () => {
    loading.value = true
    try {
        const res = await getAgentTokenList({
            current: pagination.current,
            pageSize: pagination.pageSize,
            accountName: searchParams.accountName || undefined,
            accountId: searchParams.accountId || undefined,
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
        const res = await getAgentToken(row.id)
        Object.assign(form, res.data)
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
        if (isEdit.value) {
            await updateAgentToken(editId.value, form)
            message.success('更新成功')
        } else {
            await createAgentToken(form)
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
        content: `确定删除账户「${row.accountName || row.accountId}」的token？`,
        okType: 'danger',
        okText: '确定',
        cancelText: '取消',
        onOk: async () => {
            try {
                await delAgentToken(row.id)
                message.success('删除成功')
                loadData()
            } catch (e) {
                message.error('删除失败')
            }
        },
    })
}

const maskToken = (token) => {
    if (!token) return '-'
    if (token.length <= 12) return token
    return token.slice(0, 8) + '****' + token.slice(-4)
}

const formatTime = (ts) => {
    if (!ts) return '-'
    const d = new Date(ts * 1000)
    const pad = (n) => String(n).padStart(2, '0')
    return `${d.getFullYear()}-${pad(d.getMonth() + 1)}-${pad(d.getDate())} ${pad(d.getHours())}:${pad(d.getMinutes())}:${pad(d.getSeconds())}`
}

onMounted(() => {
    loadData()
})
</script>
