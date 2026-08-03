<template>
    <div style="padding: 24px">
        <!-- 页头 -->
        <div style="display: flex; align-items: center; margin-bottom: 20px">
            <a-button
                type="text"
                @click="$router.back()"
                style="font-size: 20px; padding: 0 8px"
                >←</a-button
            >
            <h2 style="margin: 0 12px">{{ isEdit ? '编辑托管规则' : '新增托管规则' }}</h2>
            <a-typography-link type="primary">使用指引</a-typography-link>
        </div>

        <!-- 表单主体 -->
        <a-form
            ref="formRef"
            :model="form"
            :rules="rules"
            :label-col="{ span: 4 }"
            :wrapper-col="{ span: 18 }"
            style="background: #fff; padding: 24px; border-radius: 8px">
            <!-- 选择托管对象 -->
            <a-divider orientation="left">选择托管对象</a-divider>

            <a-form-item
                label="托管目标"
                name="target">
                <a-radio-group
                    v-model:value="form.target"
                    button-style="solid">
                    <a-radio-button value="account">账户</a-radio-button>
                    <a-radio-button value="project">项目</a-radio-button>
                    <a-radio-button value="ad">广告</a-radio-button>
                </a-radio-group>
            </a-form-item>

            <a-form-item
                label="选择账户"
                name="selectedAccountIds">
                <AccountTransfer
                    v-model="form.selectedAccountIds"
                    :options="accountOptions" />
            </a-form-item>

            <!-- 条件和操作 -->
            <a-divider orientation="left">条件和操作</a-divider>

            <div style="display: flex; gap: 20px">
                <div style="flex: 1">
                    <ConditionBuilder v-model="form.conditionConfig" />
                </div>

                <div style="flex: 1; border: 1px solid #d9d9d9; border-radius: 6px; padding: 16px">
                    <div style="font-weight: bold; margin-bottom: 12px">就执行以下操作</div>
                    <a-select
                        v-model:value="form.action"
                        placeholder="请选择执行动作"
                        style="width: 100%"
                        :options="actionOptions" />
                </div>
            </div>

            <!-- 执行和通知频率 -->
            <a-divider orientation="left">执行和通知频率</a-divider>

            <a-form-item
                label="检查频率"
                name="checkFreq"
                style="margin-top: 16px">
                <a-radio-group
                    v-model:value="form.checkFreq"
                    button-style="solid">
                    <a-radio-button value="15">15分钟</a-radio-button>
                    <a-radio-button value="30">30分钟</a-radio-button>
                    <a-radio-button value="60">60分钟</a-radio-button>
                    <a-radio-button value="custom">自定义</a-radio-button>
                </a-radio-group>
            </a-form-item>

            <a-form-item
                label="生效日期"
                name="dateRange">
                <a-range-picker
                    v-model:value="form.dateRange"
                    value-format="YYYY-MM-DD" />
            </a-form-item>

            <a-form-item
                label="通知方式"
                name="notifyMethods">
                <a-checkbox-group v-model:value="form.notifyMethods">
                    <a-checkbox value="sms">短信通知</a-checkbox>
                    <a-checkbox value="email">邮件通知</a-checkbox>
                    <a-checkbox value="feishu">飞书群机器人</a-checkbox>
                    <a-checkbox value="dingtalk">钉钉群机器人</a-checkbox>
                    <a-checkbox value="wecom">企业微信群机器人</a-checkbox>
                </a-checkbox-group>
            </a-form-item>

            <a-form-item
                label="规则名称"
                name="name">
                <a-input
                    v-model:value="form.name"
                    placeholder="请输入规则名称"
                    style="max-width: 400px" />
            </a-form-item>

            <a-form-item
                label=" "
                name="agreeTerms"
                :colon="false">
                <a-checkbox v-model:checked="form.agreeTerms">
                    我同意授权深圳市零一聚合根据我提交的规则信息管理我名下广告账户以及广告
                </a-checkbox>
            </a-form-item>

            <a-form-item
                label=" "
                :colon="false">
                <a-button @click="$router.back()">取消</a-button>
                <a-button
                    type="primary"
                    :loading="saving"
                    @click="handleSave"
                    style="margin-left: 8px"
                    >保存</a-button
                >
            </a-form-item>
        </a-form>
    </div>
</template>

<script setup>
import { ref, reactive, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { message } from 'ant-design-vue'
import ConditionBuilder from './components/ConditionBuilder.vue'
import AccountTransfer from './components/AccountTransfer.vue'

const route = useRoute()
const router = useRouter()
const isEdit = computed(() => route.name === 'aiHostRuleEdit')
const formRef = ref()
const saving = ref(false)

const actionOptions = [
    { label: '暂停广告', value: 'pause' },
    { label: '调整预算', value: 'adjustBudget' },
    { label: '发送通知', value: 'notify' },
    { label: '一键起量', value: 'boost' },
]

const accountOptions = ref([
    { id: '18364286876574', name: '相约派对-曝光-30-0030-79过渡' },
    { id: '1833434044638763', name: 'Ta星-nw-素材' },
    { id: '184686325214410', name: '设备-SW-素材' },
    { id: '185978646929227', name: '颇为-闪耀-iOS-头条-SC-01-首付-4516' },
    { id: '1675308379', name: '等爱-颇为-每件-7h' },
    { id: '1836428598491001', name: '好甜-nw-素材' },
    { id: '1836428588251546', name: '相约派对-nw-素材' },
    { id: '1717', name: '头杀-颇为-暖心绿洲-安卓-7H-4' },
])

const form = reactive({
    target: 'ad',
    selectedAccountIds: [],
    conditionConfig: {
        logic: 'and',
        conditions: [{ time: 'today', metric: 'cost', operator: '>=', value: 100, unit: 'yuan' }],
    },
    action: undefined,
    checkFreq: '30',
    dateRange: undefined,
    notifyMethods: ['sms'],
    name: '',
    agreeTerms: false,
})

const rules = reactive({
    target: [{ required: true, message: '请选择托管目标', trigger: 'change' }],
    selectedAccountIds: [{ type: 'array', required: true, min: 1, message: '请至少选择一个账户', trigger: 'change' }],
    action: [{ required: true, message: '请选择执行操作', trigger: 'change' }],
    checkFreq: [{ required: true, message: '请选择检查频率', trigger: 'change' }],
    dateRange: [{ required: true, message: '请选择生效日期', trigger: 'change' }],
    notifyMethods: [{ type: 'array', required: true, min: 1, message: '请至少选择一种通知方式', trigger: 'change' }],
    name: [{ required: true, message: '请输入规则名称', trigger: 'blur' }],
    agreeTerms: [
        {
            validator: (rule, value, callback) => {
                if (!value) callback(new Error('请勾选服务条款'))
                else callback()
            },
            trigger: 'change',
        },
    ],
})

const handleSave = async () => {
    if (!formRef.value) return

    try {
        await formRef.value.validate()
    } catch {
        message.warning('请完善必填信息')
        return
    }

    saving.value = true
    try {
        console.log('提交数据:', JSON.parse(JSON.stringify(form)))
        message.success(isEdit.value ? '规则已更新' : '规则已创建')
        setTimeout(() => router.push({ name: 'aiHost' }), 800)
    } finally {
        saving.value = false
    }
}
</script>
