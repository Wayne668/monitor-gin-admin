import { ToolOutlined, RobotOutlined } from '@ant-design/icons-vue'

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
        ],
    },
]
