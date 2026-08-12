import { ToolOutlined, RobotOutlined, FieldNumberOutlined, KeyOutlined, DeleteOutlined, SafetyOutlined } from '@ant-design/icons-vue'

export default [
    {
        path: 'tools',
        name: 'tools',
        component: 'RouteViewLayout',
        meta: {
            icon: ToolOutlined,
            title: '监测工具',
            isMenu: true,
            keepAlive: true,
            permission: '*',
        },
        children: [
            {
                path: 'aiHost',
                name: 'aiHost',
                component: 'aiHost/RuleList.vue',
                meta: {
                    icon: RobotOutlined,
                    title: 'AI托管',
                    isMenu: true,
                    keepAlive: true,
                    permission: '*',
                },
                children: [
                    {
                        path: 'rule/form',
                        name: 'aiHostRuleForm',
                        component: 'aiHost/RuleForm.vue',
                        meta: {
                            title: '新增托管规则',
                            isMenu: false,
                            keepAlive: false,
                            permission: '*',
                            active: 'aiHost',
                        },
                    },
                    {
                        path: 'rule/edit/:id',
                        name: 'aiHostRuleEdit',
                        component: 'aiHost/RuleForm.vue',
                        meta: {
                            title: '编辑托管规则',
                            isMenu: false,
                            keepAlive: false,
                            permission: '*',
                            active: 'aiHost',
                        },
                    },
                    {
                        path: 'rule/log/:id',
                        name: 'aiHostRuleLog',
                        component: 'aiHost/RuleLog.vue',
                        meta: {
                            title: '运行记录',
                            isMenu: false,
                            keepAlive: false,
                            permission: '*',
                            active: 'aiHost',
                        },
                    },
                    {
                        path: 'trigger-record',
                        name: 'aiHostTriggerRecord',
                        component: 'aiHost/TriggerRecordList.vue',
                        meta: {
                            title: '托管触发记录',
                            isMenu: false,
                            keepAlive: false,
                            permission: '*',
                            active: 'aiHost',
                        },
                    },
                ],
            },
            {
                path: 'hostField',
                name: 'hostField',
                component: 'aiHost/FieldList.vue',
                meta: {
                    icon: FieldNumberOutlined,
                    title: '字段管理',
                    isMenu: true,
                    keepAlive: true,
                    permission: '*',
                },
                children: [
                    {
                        path: 'form',
                        name: 'hostFieldForm',
                        component: 'aiHost/FieldForm.vue',
                        meta: {
                            title: '新增字段',
                            isMenu: false,
                            keepAlive: false,
                            permission: '*',
                            active: 'hostField',
                        },
                    },
                    {
                        path: 'edit/:id',
                        name: 'hostFieldEdit',
                        component: 'aiHost/FieldForm.vue',
                        meta: {
                            title: '编辑字段',
                            isMenu: false,
                            keepAlive: false,
                            permission: '*',
                            active: 'hostField',
                        },
                    },
                ],
            },
            {
                path: 'materialWarn',
                name: 'materialWarn',
                component: 'exception/404.vue',
                meta: {
                    title: '素材预警',
                    isMenu: true,
                    keepAlive: true,
                    permission: '*',
                },
            },
            {
                path: 'materialDel',
                name: 'materialDel',
                component: 'aiHost/MaterialDel.vue',
                meta: {
                    icon: DeleteOutlined,
                    title: '素材删除',
                    isMenu: true,
                    keepAlive: true,
                    permission: '*',
                },
                children: [
                    {
                        path: 'form',
                        name: 'materialDelForm',
                        component: 'aiHost/MaterialDelForm.vue',
                        meta: {
                            title: '素材删除操作',
                            isMenu: false,
                            keepAlive: false,
                            permission: '*',
                            active: 'materialDel',
                        },
                    },
                ],
            },
            {
                path: 'agentToken',
                name: 'agentToken',
                component: 'aiHost/TokenList.vue',
                meta: {
                    icon: KeyOutlined,
                    title: '账户token',
                    isMenu: true,
                    keepAlive: true,
                    permission: '*',
                },
            },
            {
                path: 'hostAccount',
                name: 'hostAccount',
                component: 'aiHost/HostAccount.vue',
                meta: {
                    icon: SafetyOutlined,
                    title: '托管账户',
                    isMenu: true,
                    keepAlive: true,
                    permission: '*',
                },
            },
        ],
    },
]
