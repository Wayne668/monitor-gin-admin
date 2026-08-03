<template>
    <div class="condition-builder">
        <div class="builder-header">
            <span class="title">当以下情况发生时：</span>
            <a-radio-group
                v-model:value="logic"
                size="small">
                <a-radio-button value="and">满足全部 (AND)</a-radio-button>
                <a-radio-button value="or">满足任意 (OR)</a-radio-button>
            </a-radio-group>
        </div>

        <div class="condition-list">
            <div
                v-for="(item, index) in conditions"
                :key="item._id"
                class="condition-row">
                <a-select
                    v-model:value="item.time"
                    placeholder="时间范围"
                    style="width: 120px"
                    :options="timeOptions" />
                <a-select
                    v-model:value="item.metric"
                    placeholder="监控指标"
                    style="width: 140px"
                    :options="metricOptions" />
                <a-select
                    v-model:value="item.operator"
                    placeholder="比较"
                    style="width: 110px"
                    :options="operatorOptions" />
                <a-input-number
                    v-model:value="item.value"
                    :min="0"
                    :precision="2"
                    placeholder="阈值"
                    style="width: 140px" />
                <a-select
                    v-model:value="item.unit"
                    placeholder="单位"
                    style="width: 90px"
                    :options="unitOptions" />
                <a-button
                    type="link"
                    danger
                    :disabled="conditions.length <= 1"
                    @click="removeCondition(index)">
                    <template #icon><DeleteOutlined /></template>
                </a-button>
            </div>
        </div>

        <a-button
            type="dashed"
            block
            size="small"
            @click="addCondition"
            style="margin-top: 12px">
            <template #icon><PlusOutlined /></template>
            添加条件
        </a-button>
    </div>
</template>

<script setup>
import { ref, watch } from 'vue'
import { DeleteOutlined, PlusOutlined } from '@ant-design/icons-vue'

const props = defineProps({
    modelValue: { type: Object, default: () => ({ logic: 'and', conditions: [] }) },
})

const emit = defineEmits(['update:modelValue'])

const timeOptions = [
    { label: '当天', value: 'today' },
    { label: '近7天', value: 'last_7d' },
    { label: '近30天', value: 'last_30d' },
]
const metricOptions = [
    { label: '消耗', value: 'cost' },
    { label: '展现量', value: 'impression' },
    { label: '点击量', value: 'click' },
    { label: '转化数', value: 'conversion' },
    { label: '平均转化成本', value: 'cpa' },
    { label: 'ROI', value: 'roi' },
]
const operatorOptions = [
    { label: '大于等于', value: '>=' },
    { label: '小于等于', value: '<=' },
    { label: '大于', value: '>' },
    { label: '小于', value: '<' },
    { label: '等于', value: '==' },
]
const unitOptions = [
    { label: '元', value: 'yuan' },
    { label: '个', value: 'count' },
    { label: '%', value: 'percent' },
]

const logic = ref(props.modelValue.logic || 'and')

let _seq = 0
const genId = () => `cond_${Date.now()}_${_seq++}`

const conditions = ref((props.modelValue.conditions || []).map((c) => ({ ...c, _id: genId() })))

watch(
    () => props.modelValue,
    (val) => {
        logic.value = val.logic
        if (val.conditions.length !== conditions.value.length) {
            conditions.value = val.conditions.map((c) => ({ ...c, _id: genId() }))
        }
    },
    { deep: true }
)

const syncToParent = () => {
    // eslint-disable-next-line no-unused-vars
    emit('update:modelValue', { logic: logic.value, conditions: conditions.value.map(({ _id, ...rest }) => rest) })
}
watch([logic, conditions], syncToParent, { deep: true })

const addCondition = () => {
    conditions.value.push({ _id: genId(), time: 'today', metric: '', operator: '>=', value: undefined, unit: 'yuan' })
}
const removeCondition = (index) => {
    if (conditions.value.length > 1) conditions.value.splice(index, 1)
}
</script>

<style scoped>
.condition-builder {
    background: #fafafa;
    border: 1px solid #f0f0f0;
    border-radius: 6px;
    padding: 16px;
}
.builder-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 14px;
}
.title {
    font-weight: 600;
    color: #262626;
    font-size: 13px;
}
.condition-row {
    display: flex;
    gap: 8px;
    align-items: center;
    margin-bottom: 10px;
}
</style>
