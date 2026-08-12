<template>
    <div style="padding: 24px">
        <a-card>
            <template #title>
                {{ isEdit ? '编辑字段' : '新增字段' }}
            </template>
            <template #extra>
                <a-button @click="handleBack">返回</a-button>
            </template>

            <a-form
                ref="formRef"
                :model="form"
                :rules="rules"
                layout="vertical"
                style="max-width: 600px">
                <a-form-item
                    label="字段"
                    name="field">
                    <a-input
                        v-model:value="form.field"
                        placeholder="请输入字段" />
                </a-form-item>
                <a-form-item
                    label="字段名称"
                    name="name">
                    <a-input
                        v-model:value="form.name"
                        placeholder="请输入字段名称" />
                </a-form-item>
                <a-form-item
                    label="分类"
                    name="cate">
                    <a-select
                        v-model:value="form.cate"
                        placeholder="请选择分类">
                        <a-select-option value="dimension">维度</a-select-option>
                        <a-select-option value="metric">指标</a-select-option>
                    </a-select>
                </a-form-item>
                <a-form-item
                    label="可恢复"
                    name="stash">
                    <a-select
                        v-model:value="form.stash"
                        placeholder="请选择是否可恢复">
                        <a-select-option :value="1">可恢复</a-select-option>
                        <a-select-option :value="2">不可恢复</a-select-option>
                    </a-select>
                </a-form-item>
                <a-form-item
                    label="单位"
                    name="unit">
                    <a-input
                        v-model:value="form.unit"
                        placeholder="请输入单位" />
                </a-form-item>
                <a-form-item
                    label="公式"
                    name="formula">
                    <a-textarea
                        v-model:value="form.formula"
                        placeholder="请输入公式"
                        :rows="3" />
                </a-form-item>
                <a-form-item
                    label="开启"
                    name="status">
                    <a-switch
                        :checked="form.status === 1"
                        @change="(val) => (form.status = val ? 1 : 0)" />
                </a-form-item>
                <a-form-item>
                    <a-button
                        type="primary"
                        :loading="submitting"
                        @click="handleSubmit">
                        保存
                    </a-button>
                    <a-button
                        style="margin-left: 8px"
                        @click="handleBack">
                        取消
                    </a-button>
                </a-form-item>
            </a-form>
        </a-card>
    </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { message } from 'ant-design-vue'
import { getHostField, createHostField, updateHostField } from '@/apis/modules/hostField'

const router = useRouter()
const route = useRoute()

const formRef = ref()
const submitting = ref(false)
const isEdit = ref(false)

const form = reactive({
    field: '',
    name: '',
    cate: undefined,
    stash: undefined,
    unit: '',
    formula: '',
    status: 1,
})

const rules = {
    field: [{ required: true, message: '请输入字段', trigger: 'blur' }],
    name: [{ required: true, message: '请输入字段名称', trigger: 'blur' }],
    cate: [{ required: true, message: '请选择分类', trigger: 'change' }],
    stash: [{ required: true, message: '请选择是否可恢复', trigger: 'change' }],
}

const loadData = async () => {
    const id = route.params.id
    if (!id) return
    isEdit.value = true
    try {
        const res = await getHostField(id)
        Object.assign(form, res.data)
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
            await updateHostField(route.params.id, form)
            message.success('更新成功')
        } else {
            await createHostField(form)
            message.success('创建成功')
        }
        handleBack()
    } catch (e) {
        message.error(isEdit.value ? '更新失败' : '创建失败')
    } finally {
        submitting.value = false
    }
}

const handleBack = () => {
    router.push({ name: 'hostField' })
}

onMounted(() => {
    loadData()
})
</script>
