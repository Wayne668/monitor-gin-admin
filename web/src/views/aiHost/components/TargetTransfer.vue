<template>
    <div class="target-transfer">
        <div class="panel left-panel">
            <div class="panel-header">
                <a-form-item-rest>
                    <a-input
                        v-model:value="keyword"
                        placeholder="搜索名称/ID"
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
                        class="target-item">
                        <a-checkbox :value="item.id">
                            {{ item.id }} - {{ item.name }}
                            <span class="account-tag">[{{ item.advertiserId }}]</span>
                        </a-checkbox>
                    </div>
                </a-checkbox-group>
                <a-empty
                    v-if="filteredList.length === 0"
                    :description="emptyText" />
            </div>
        </div>

        <div class="panel right-panel">
            <div class="panel-header">
                <span>已选 {{ selectedItems.length }} 个{{ targetLabel }}</span>
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
        </div>
    </div>
</template>

<script setup>
import { ref, computed, watch } from 'vue'
import { SearchOutlined } from '@ant-design/icons-vue'

const props = defineProps({
    // 目标列表（由父组件从接口获取后传入）
    options: { type: Array, default: () => [] },
    // 选中的目标ID列表（v-model）
    modelValue: { type: Array, default: () => [] },
    // 目标类型，用于文案展示
    target: { type: String, default: '' },
})

const emit = defineEmits(['update:modelValue'])

const keyword = ref('')

const targetLabel = computed(() => {
    if (props.target === 'promotion') return '广告'
    if (props.target === 'creative') return '创意'
    return '目标'
})

const emptyText = computed(() => {
    if (!props.options || props.options.length === 0) return '请先选择账户并加载数据'
    return '无匹配数据'
})

const checkedKeys = ref([...props.modelValue])

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

// 外部 options 变化时，清理已不在列表中的选中项
watch(
    () => props.options,
    (val) => {
        const validIds = new Set(val.map((o) => o.id))
        const filtered = checkedKeys.value.filter((k) => validIds.has(k))
        if (filtered.length !== checkedKeys.value.length) {
            checkedKeys.value = filtered
            emit('update:modelValue', [...filtered])
        }
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
.target-transfer {
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
.target-item {
    padding: 4px 0;
}
.account-tag {
    color: #999;
    font-size: 12px;
    margin-left: 4px;
}
</style>
