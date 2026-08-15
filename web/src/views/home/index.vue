<template>
    <div class="welcome-content">
        <div class="welcome-icon-wrapper">
            <svg class="welcome-icon" viewBox="0 0 1024 1024" xmlns="http://www.w3.org/2000/svg">
                <path
                    d="M512 64C264.6 64 64 264.6 64 512s200.6 448 448 448 448-200.6 448-448S759.4 64 512 64zm0 820c-205.4 0-372-166.6-372-372s166.6-372 372-372 372 166.6 372 372-166.6 372-372 372z"
                    fill="#E8F3FF" />
                <path
                    d="M512 256c-141.4 0-256 114.6-256 256s114.6 256 256 256 256-114.6 256-256-114.6-256-256-256zm0 448c-106 0-192-86-192-192s86-192 192-192 192 86 192 192-86 192-192 192z"
                    fill="#5B9BFF" />
                <circle cx="512" cy="480" r="80" fill="#5B9BFF" />
                <path
                    d="M512 320c-88.4 0-160 71.6-160 160 0 44.2 17.9 84.2 47 113l56.6-56.6c-14.6-14.6-23.6-34.8-23.6-57.4 0-44.2 35.8-80 80-80s80 35.8 80 80c0 22.6-9 42.8-23.6 57.4L625 593c29.1-28.8 47-68.8 47-113 0-88.4-71.6-160-160-160z"
                    fill="#3D7DD8" />
            </svg>
        </div>
        <h1 class="welcome-title">欢迎回来</h1>
        <p class="welcome-subtitle">Hi, {{ userInfo?.name || userInfo?.username || 'Admin' }}</p>
    </div>
</template>

<script setup>
import { storeToRefs } from 'pinia'
import { useUserStore } from '@/store'

defineOptions({
    name: 'home',
})

const userStore = useUserStore()
const { userInfo } = storeToRefs(userStore)
</script>

<style lang="less" scoped>
.welcome-content {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    min-height: 60vh;
    text-align: center;
}

.welcome-icon-wrapper {
    margin-bottom: 24px;
}

.welcome-icon {
    width: 96px;
    height: 96px;
    animation: floatUpDown 3s ease-in-out infinite;
}

@keyframes floatUpDown {
    0%,
    100% {
        transform: translateY(0);
    }
    50% {
        transform: translateY(-12px);
    }
}

.welcome-title {
    font-size: 28px;
    font-weight: 600;
    color: rgba(0, 0, 0, 0.88);
    margin: 0 0 12px;
}

.welcome-subtitle {
    font-size: 16px;
    font-weight: 600;
    color: #7c6ee6;
    margin: 0;
}
</style>
