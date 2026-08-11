<template>
    <div style="padding: 24px">
        <div style="display: flex; align-items: center; margin-bottom: 20px">
            <a-button
                type="text"
                @click="$router.back()"
                style="font-size: 20px; padding: 0 8px">
                ←
            </a-button>
            <h2 style="margin: 0 12px">运行记录</h2>
        </div>

        <a-card>
            <a-empty description="运行记录功能开发中..." />
        </a-card>
    </div>
</template>

<script setup>
// 后期实现
</script>