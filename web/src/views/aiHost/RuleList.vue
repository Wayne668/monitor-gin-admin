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
                :data-source="tableData"
                :pagination="pagination"
                :row-key="(record) => record.id">
                <a-table-column
                    type="checkbox"
                    width="55" />
                <a-table-column
                    title="规则ID"
                    data-index="id"
                    :width="100" />
                <a-table-column
                    title="规则名称"
                    data-index="name" />
                <a-table-column
                    title="应用场景"
                    data-index="scene"
                    :width="160" />
                <a-table-column
                    title="媒体"
                    data-index="media"
                    :width="120" />
                <a-table-column
                    title="托管目标"
                    data-index="target"
                    :width="100" />
                <a-table-column
                    title="状态"
                    data-index="status"
                    :width="100">
                    <template #bodyCell="{ record }">
                        <a-tag :color="record.status === '启用' ? 'success' : 'default'">
                            {{ record.status }}
                        </a-tag>
                    </template>
                </a-table-column>
                <a-table-column
                    title="操作"
                    :width="200"
                    fixed="right">
                    <template #bodyCell="{ record }">
                        <a-button
                            type="link"
                            size="small"
                            @click="handleEdit(record)">
                            编辑
                        </a-button>
                        <a-button
                            type="link"
                            size="small"
                            @click="handleToggle(record)">
                            {{ record.status === '启用' ? '停用' : '启用' }}
                        </a-button>
                        <a-button
                            type="link"
                            danger
                            size="small"
                            @click="handleDelete(record)">
                            删除
                        </a-button>
                    </template>
                </a-table-column>
            </a-table>
        </a-card>
    </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { Modal, message } from 'ant-design-vue'

const router = useRouter()

const pagination = {
    current: 1,
    pageSize: 10,
    total: 0,
    showTotal: (total) => `共 ${total} 条`,
    showSizeChanger: true,
}

const tableData = ref([
    {
        id: '001',
        name: '广告止损规则-日常版',
        scene: '广告止损',
        media: '巨量广告2.0',
        target: '广告',
        status: '启用',
    },
    {
        id: '002',
        name: '预算优化-晚间投放',
        scene: '自动化优化预算',
        media: '巨量广告2.0',
        target: '项目',
        status: '停用',
    },
    {
        id: '003',
        name: '异常预警-成本超标',
        scene: '异常预警',
        media: '磁力智投',
        target: '账户',
        status: '启用',
    },
])
pagination.total = tableData.value.length

const handleAdd = () => {
    router.push({ name: 'aiHostRuleForm' })
}

const handleEdit = (row) => {
    router.push({ name: 'aiHostRuleEdit', params: { id: row.id } })
}

const handleToggle = (row) => {
    row.status = row.status === '启用' ? '停用' : '启用'
    message.success(`已${row.status}`)
}

const handleDelete = (row) => {
    Modal.confirm({
        title: '确认删除',
        content: `确定删除规则「${row.name}」？`,
        okType: 'danger',
        okText: '确定',
        cancelText: '取消',
        onOk() {
            tableData.value = tableData.value.filter((item) => item.id !== row.id)
            pagination.total = tableData.value.length
            message.success('删除成功')
        },
    })
}
</script>
