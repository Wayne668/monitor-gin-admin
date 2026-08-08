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
            <!-- 授权账户 -->
            <a-divider orientation="left">授权账户</a-divider>

            <a-form-item
                label="选择代理商账户"
                name="selectedAccountId">
                <a-select
                    v-model:value="form.selectedAccountId"
                    placeholder="请选择代理商账户（用于获取 access_token）"
                    style="max-width: 400px"
                    :options="agentTokenOptions"
                    @change="handleAgentTokenChange" />
            </a-form-item>

            <!-- 选择托管对象 -->
            <a-divider orientation="left">选择托管对象</a-divider>

            <a-form-item
                label="托管目标"
                name="target">
                <a-radio-group
                    v-model:value="form.target"
                    button-style="solid">
                    <a-radio-button value="account">账户</a-radio-button>
                    <a-radio-button
                        value="project"
                        disabled>
                        项目
                    </a-radio-button>
                    <a-radio-button value="promotion">广告</a-radio-button>
                    <a-radio-button value="creative">创意</a-radio-button>
                </a-radio-group>
            </a-form-item>

            <a-form-item
                v-if="form.target === 'promotion' || form.target === 'creative'"
                label="范围类型"
                name="scopeType">
                <a-radio-group
                    v-model:value="form.scopeType"
                    button-style="solid">
                    <a-radio-button
                        v-for="opt in scopeTypeOptions"
                        :key="opt.value"
                        :value="opt.value">
                        {{ opt.label }}
                    </a-radio-button>
                </a-radio-group>
            </a-form-item>

            <a-form-item
                label="选择账户"
                name="selectedAccountIds">
                <AccountTransfer
                    v-model="form.selectedAccountIds"
                    :options="accountOptions"
                    :target="form.target"
                    :scope-type="form.scopeType"
                    @targets-loaded="handleTargetsLoaded" />
                <a-alert
                    v-if="showTargetTransfer && targetOptions.length === 0"
                    type="info"
                    show-icon
                    message="请先在上方选择账户并点击「确定」加载目标数据"
                    style="margin-top: 8px" />
            </a-form-item>

            <a-form-item
                v-if="showTargetTransfer && targetOptions.length > 0"
                :label="isExcludeScope ? '排除目标' : '选择目标'"
                name="selectedTargetIds">
                <TargetTransfer
                    v-model="form.selectedTargetIds"
                    :options="targetOptions"
                    :target="form.target" />
            </a-form-item>

            <!-- 条件和操作 -->
            <a-divider orientation="left">条件和操作</a-divider>

            <a-form-item :wrapper-col="{ span: 24 }">
                <div style="display: flex; gap: 20px; margin-left: 100px">
                    <div style="flex: 1 1 0%; min-width: 0; max-width: 60%; overflow: auto">
                        <ConditionBuilder v-model="form.conditionConfig" />
                    </div>

                    <div
                        style="
                            flex: 1 1 0%;
                            min-width: 0;
                            max-width: 25%;
                            border: 1px solid #d9d9d9;
                            border-radius: 6px;
                            padding: 10px;
                        ">
                        <div style="font-weight: bold; margin-bottom: 12px">就执行以下操作</div>
                        <a-select
                            v-model:value="form.action"
                            placeholder="请选择执行动作"
                            style="width: 100%"
                            :options="actionOptions" />
                    </div>
                </div>
            </a-form-item>

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
                    <a-checkbox value="sms" disabled>短信通知</a-checkbox>
                    <a-checkbox value="dingtalk">钉钉群机器人</a-checkbox>
                </a-checkbox-group>
            </a-form-item>

            <a-form-item
                v-if="form.notifyMethods.includes('sms')"
                label="手机号"
                name="notifyPhones">
                <a-input
                    v-model:value="form.notifyPhones"
                    placeholder="多个手机号用英文逗号分隔"
                    style="max-width: 400px" />
            </a-form-item>

            <a-form-item
                v-if="form.notifyMethods.includes('dingtalk')"
                label="Webhook URL"
                name="dingtalkWebhookUrl">
                <a-input
                    v-model:value="form.dingtalkWebhookUrl"
                    placeholder="请输入钉钉群机器人 Webhook URL"
                    style="max-width: 400px" />
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
import { ref, reactive, computed, watch, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { message } from 'ant-design-vue'
import ConditionBuilder from './components/ConditionBuilder.vue'
import AccountTransfer from './components/AccountTransfer.vue'
import TargetTransfer from './components/TargetTransfer.vue'
import { getEnabledAccountList } from '@/apis/modules/accountInfo'
import { getAgentTokenList } from '@/apis/modules/agentToken'
import { saveHostRule } from '@/apis/modules/hostRule'

const route = useRoute()
const router = useRouter()
const isEdit = computed(() => route.name === 'aiHostRuleEdit')
const formRef = ref()
const saving = ref(false)

const promotionActionOptions = [
    { label: '暂停广告+发送通知', value: 'pause' },
    { label: '启动广告+发送通知', value: 'restart' },
    { label: '发送通知', value: 'notify' },
]

const creativeActionOptions = [
    { label: '暂停素材+发送通知', value: 'pause' },
    { label: '启动素材+发送通知', value: 'restart' },
    { label: '删除素材+发送通知', value: 'delete' },
    { label: '发送通知', value: 'notify' },
]

// 根据 target 类型动态切换操作选项
const actionOptions = computed(() => {
    if (form.target === 'creative') return creativeActionOptions
    if (form.target === 'promotion') return promotionActionOptions
    return []
})

// 范围类型选项：根据 target 切换文案
const scopeTypeOptions = computed(() => {
    if (form.target === 'promotion') {
        return [
            { label: '指定账户的广告', value: 'account_promotion' },
            { label: '指定广告', value: 'promotion' },
            { label: '排除指定账户的广告', value: 'exclude_account_promotion' },
            { label: '排除指定广告', value: 'exclude_promotion' },
        ]
    }
    if (form.target === 'creative') {
        return [
            { label: '指定账户的创意', value: 'account_creative' },
            { label: '指定创意', value: 'creative' },
            { label: '排除指定账户的创意', value: 'exclude_account_creative' },
            { label: '排除指定创意', value: 'exclude_creative' },
        ]
    }
    return []
})

// 是否是排除模式
const isExcludeScope = computed(() => {
    return ['exclude_account_promotion', 'exclude_promotion', 'exclude_account_creative', 'exclude_creative'].includes(
        form.scopeType
    )
})

// 是否需要展示目标穿梭框（广告/创意选择）
// 仅"指定广告/创意"和"排除指定广告/创意"需要选择具体目标
// "指定账户的广告/创意"和"排除指定账户的广告/创意"针对账户下全部目标，无需穿梭框
const needTargetSelect = computed(() => {
    return ['promotion', 'creative', 'exclude_promotion', 'exclude_creative'].includes(form.scopeType)
})

const showTargetTransfer = computed(() => {
    return (form.target === 'promotion' || form.target === 'creative') && needTargetSelect.value
})

const accountOptions = ref([])
const agentTokenOptions = ref([])

const loadAgentTokens = async () => {
    try {
        const res = await getAgentTokenList({ pageSize: 999 })
        agentTokenOptions.value = (res.data || []).map((item) => ({
            label: `${item.accountName} (${item.accountId})`,
            value: item.accountId,
        }))
    } catch (e) {
        message.error('加载代理商账户失败')
    }
}

const loadAccounts = async () => {
    try {
        // 只取下拉所需字段，减少传输量；limit 使用后端默认（最新 100 条）
        const res = await getEnabledAccountList({ fields: 'advertiser_id,advertiser_name' })
        accountOptions.value = (res.data || []).map((item) => ({
            id: String(item.advertiserId),
            name: item.advertiserName,
        }))
    } catch (e) {
        message.error('加载账户列表失败')
    }
}

onMounted(() => {
    loadAccounts()
    loadAgentTokens()
})

const form = reactive({
    target: 'account',
    scopeType: undefined,
    selectedAccountId: undefined,
    selectedAccountIds: [],
    selectedTargetIds: [],
    targetPromotion: [],
    targetMaterial: [],
    conditionConfig: {
        logic: 'and',
        conditions: [{ time: '', metric: '', operator: '', value: '', unit: '' }],
    },
    action: undefined,
    checkFreq: '30',
    dateRange: undefined,
    notifyMethods: ['dingtalk'],
    notifyPhones: '',
    dingtalkWebhookUrl: '',
    name: '',
    agreeTerms: false,
})

// 目标穿梭框数据源：根据 target 选择对应的列表
const targetOptions = computed(() => {
    if (form.target === 'promotion') return form.targetPromotion
    if (form.target === 'creative') return form.targetMaterial
    return []
})

const handleTargetsLoaded = (items) => {
    if (form.target === 'promotion') {
        form.targetPromotion = items
    } else if (form.target === 'creative') {
        form.targetMaterial = items
    }
}

// 切换 target 时清空相关字段，避免数据混乱
watch(
    () => form.target,
    () => {
        form.action = undefined
        form.scopeType = undefined
        form.selectedTargetIds = []
        form.targetPromotion = []
        form.targetMaterial = []
    }
)

// 切换 scopeType 时清空已选目标和已加载的目标数据，避免与范围不匹配
watch(
    () => form.scopeType,
    (newVal, oldVal) => {
        // 仅在进出"需要选目标"的模式时清空已加载的目标数据
        const wasNeedSelect = ['promotion', 'creative', 'exclude_promotion', 'exclude_creative'].includes(oldVal)
        const nowNeedSelect = ['promotion', 'creative', 'exclude_promotion', 'exclude_creative'].includes(newVal)
        if (wasNeedSelect !== nowNeedSelect) {
            form.targetPromotion = []
            form.targetMaterial = []
        }
        form.selectedTargetIds = []
    }
)

const rules = reactive({
    target: [{ required: true, message: '请选择托管目标', trigger: 'change' }],
    scopeType: [{ required: true, message: '请选择范围类型', trigger: 'change' }],
    selectedAccountId: [{ required: true, message: '请选择代理商账户', trigger: 'change' }],
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
        const payload = {
            ruleName: form.name,
            target: form.target,
            scopeType: form.scopeType,
            selectedAccountId: form.selectedAccountId,
            selectedAccountIds: form.selectedAccountIds,
            selectedTargetIds: form.selectedTargetIds,
            conditionConfig: form.conditionConfig,
            action: form.action,
            checkFreq: parseInt(form.checkFreq === 'custom' ? '30' : form.checkFreq, 10),
            dateRange: form.dateRange || [],
            notifyMethods: form.notifyMethods,
            dingtalkWebhookUrl: form.dingtalkWebhookUrl,
        }
        await saveHostRule(payload)
        message.success(isEdit.value ? '规则已更新' : '规则已创建')
        setTimeout(() => router.push({ name: 'aiHost' }), 800)
    } catch (e) {
        message.error(e?.message || '保存失败')
    } finally {
        saving.value = false
    }
}

const handleAgentTokenChange = () => {
    // 代理商账户变更时，可重置相关状态
}
</script>
