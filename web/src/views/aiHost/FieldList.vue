<template>
    <div style="padding: 24px">
        <a-card style="margin-bottom: 16px">
            <div style="display: flex; justify-content: space-between; align-items: center">
                <h2 style="margin: 0">字段管理</h2>
                <a-button
                    type="primary"
                    size="large"
                    @click="handleAdd">
                    ＋ 新增字段
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
                    <template v-if="'cate' === column.key">
                        <a-tag :color="record.cate === 'dimension' ? 'blue' : 'green'">
                            {{ record.cate === 'dimension' ? '维度' : '指标' }}
                        </a-tag>
                    </template>
                    <template v-if="'stash' === column.key">
                        <a-tag :color="record.stash === 1 ? 'success' : 'default'">
                            {{ record.stash === 1 ? '可恢复' : '不可恢复' }}
                        </a-tag>
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
    </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { Modal, message } from 'ant-design-vue'
import { getHostFieldList, delHostField } from '@/apis/modules/hostField'

const router = useRouter()

const columns = [
    { title: 'ID', dataIndex: 'id', key: 'id', width: 80 },
    { title: '字段', dataIndex: 'field', key: 'field', width: 150 },
    { title: '字段名称', dataIndex: 'name', key: 'name', width: 150 },
    { title: '分类', dataIndex: 'cate', key: 'cate', width: 100 },
    { title: '可恢复', dataIndex: 'stash', key: 'stash', width: 100 },
    { title: '单位', dataIndex: 'unit', key: 'unit', width: 80 },
    { title: '公式', dataIndex: 'formula', key: 'formula' },
    { title: '操作', key: 'action', width: 200, fixed: 'right' },
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

const loadData = async () => {
    loading.value = true
    try {
        const res = await getHostFieldList({
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
    router.push({ name: 'hostFieldForm' })
}

const handleEdit = (row) => {
    router.push({ name: 'hostFieldEdit', params: { id: row.id } })
}

const handleDelete = (row) => {
    Modal.confirm({
        title: '确认删除',
        content: `确定删除字段「${row.name}」？`,
        okType: 'danger',
        okText: '确定',
        cancelText: '取消',
        onOk: async () => {
            try {
                await delHostField(row.id)
                message.success('删除成功')
                loadData()
            } catch (e) {
                message.error('删除失败')
            }
        },
    })
}

onMounted(() => {
    loadData()
})
</script>
