<template>
    <div class="account-transfer">
        <div class="panel left-panel">
            <div class="panel-header">
                <a-form-item-rest>
                    <a-input
                        v-model:value="keyword"
                        placeholder="搜索账户名称/ID"
                        size="small"
                        allow-clear>
                        <template #prefix><SearchOutlined /></template>
                    </a-input>
                </a-form-item-rest>
            </div>
            <div class="panel-toolbar">
                <a-form-item-rest>
                    <a-checkbox
                        :checked="isAllSelected"
                        :indeterminate="isIndeterminate"
                        @change="handleSelectAll">
                        全选 ({{ filteredList.length }})
                    </a-checkbox>
                </a-form-item-rest>
            </div>
            <div class="panel-body">
                <a-checkbox-group
                    v-model:value="checkedKeys"
                    @change="onCheckedChange">
                    <div
                        v-for="item in filteredList"
                        :key="item.id"
                        class="account-item">
                        <a-checkbox :value="item.id">{{ item.name }} （{{ item.id }}）</a-checkbox>
                    </div>
                </a-checkbox-group>
                <a-empty
                    v-if="filteredList.length === 0"
                    description="无匹配账户" />
            </div>
        </div>

        <div class="panel right-panel">
            <div class="panel-header">
                <span>已选 {{ selectedItems.length }} 个账户</span>
                <a-button
                    type="link"
                    size="small"
                    :disabled="selectedItems.length === 0"
                    @click="clearAll">
                    清空
                </a-button>
            </div>
            <div class="panel-body tags-container">
                <a-tag
                    v-for="item in selectedItems"
                    :key="item.id"
                    closable
                    @close="removeItem(item.id)"
                    style="margin: 4px">
                    {{ item.name }}
                </a-tag>
                <a-empty
                    v-if="selectedItems.length === 0"
                    description="暂未选择" />
            </div>
            <div
                v-if="showConfirmBtn"
                class="panel-footer">
                <a-button
                    type="primary"
                    size="small"
                    :loading="loading"
                    :disabled="selectedItems.length === 0"
                    @click="handleConfirm">
                    确定
                </a-button>
            </div>
        </div>
    </div>
</template>

<script setup>
import { ref, computed, watch } from 'vue'
import { SearchOutlined } from '@ant-design/icons-vue'
import { message } from 'ant-design-vue'
import { getTargetByAccount } from '@/apis/modules/hostRule'

const props = defineProps({
    options: { type: Array, default: () => [] },
    modelValue: { type: Array, default: () => [] },
    target: { type: String, default: '' },
    scopeType: { type: String, default: '' },
    agentId: { type: [String, Number], default: undefined },
})

const emit = defineEmits(['update:modelValue', 'targets-loaded'])

const keyword = ref('')
const checkedKeys = ref([...props.modelValue])
const loading = ref(false)

// 仅"指定广告/创意"和"排除指定广告/创意"需要选具体目标，才显示「确定」按钮加载目标数据
const showConfirmBtn = computed(() => {
    if (props.target !== 'creative' && props.target !== 'promotion') return false
    return ['promotion', 'creative', 'exclude_promotion', 'exclude_creative'].includes(props.scopeType)
})

const filteredList = computed(() => {
    const kw = keyword.value.trim().toLowerCase()
    if (!kw) return props.options
    return props.options.filter((o) => o.name.toLowerCase().includes(kw) || String(o.id).toLowerCase().includes(kw))
})

const selectedItems = computed(() => props.options.filter((o) => checkedKeys.value.includes(o.id)))

const isAllSelected = computed(
    () => filteredList.value.length > 0 && filteredList.value.every((o) => checkedKeys.value.includes(o.id))
)
const isIndeterminate = computed(
    () => !isAllSelected.value && filteredList.value.some((o) => checkedKeys.value.includes(o.id))
)

const handleSelectAll = (e) => {
    const ids = filteredList.value.map((o) => o.id)
    if (e.target.checked) {
        checkedKeys.value = [...new Set([...checkedKeys.value, ...ids])]
    } else {
        checkedKeys.value = checkedKeys.value.filter((k) => !ids.includes(k))
    }
    // 同步 emit，确保父组件 v-model 立即更新，避免表单校验拿到旧值
    emit('update:modelValue', [...checkedKeys.value])
}

const removeItem = (id) => {
    checkedKeys.value = checkedKeys.value.filter((k) => k !== id)
    emit('update:modelValue', [...checkedKeys.value])
}
const clearAll = () => {
    checkedKeys.value = []
    emit('update:modelValue', [])
}

// a-checkbox-group change 事件：同步 emit
const onCheckedChange = (val) => {
    emit('update:modelValue', [...val])
}

const handleConfirm = async () => {
    if (checkedKeys.value.length === 0) {
        return
    }
    loading.value = true
    try {
        const res = await getTargetByAccount({
            target: props.target,
            accountIds: [...checkedKeys.value],
            agentId: props.agentId,
        })
        const items = res?.data || []
        message.success(`已加载 ${items.length} 条目标数据`)
        emit('targets-loaded', items)
    } catch (e) {
        // 错误信息由请求拦截器统一处理
    } finally {
        loading.value = false
    }
}

// 父组件 -> 子组件同步（如编辑回填）
watch(
    () => props.modelValue,
    (val) => {
        // 相等则跳过，避免与 emit 形成递归更新
        if (arraysEqual(checkedKeys.value, val)) return
        checkedKeys.value = [...val]
    },
    { deep: true }
)

function arraysEqual(a, b) {
    if (a === b) return true
    if (!a || !b || a.length !== b.length) return false
    const setB = new Set(b)
    for (const x of a) {
        if (!setB.has(x)) return false
    }
    return true
}
</script>

<style scoped>
.account-transfer {
    display: flex;
    gap: 16px;
}
.panel {
    border: 1px solid #d9d9d9;
    border-radius: 6px;
    display: flex;
    flex-direction: column;
    overflow: hidden;
}
.left-panel {
    flex: 1;
    min-width: 0;
    max-width: 50%;
}
.right-panel {
    flex: 1;
    min-width: 0;
    max-width: 40%;
}
.panel-header {
    padding: 10px 12px;
    border-bottom: 1px solid #f0f0f0;
    background: #fafafa;
}
.panel-toolbar {
    padding: 8px 12px;
    border-bottom: 1px solid #f0f0f0;
}
.panel-body {
    flex: 1;
    padding: 8px 12px;
    overflow-y: auto;
    max-height: 280px;
    min-height: 200px;
}
.panel-footer {
    padding: 8px 12px;
    border-top: 1px solid #f0f0f0;
    background: #fafafa;
    text-align: right;
}
.account-item {
    padding: 4px 0;
}
</style>
