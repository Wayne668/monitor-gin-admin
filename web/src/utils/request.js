import { message } from 'ant-design-vue'
import jschardet from 'jschardet'
import XYHttp from 'xy-http'
import { config } from '@/config'

import { useUserStore } from '@/store'

const MSG_ERROR_KEY = Symbol('GLOBAL_ERROR')

const options = {
    enableAbortController: true,
    transformResponse: [
        (data) => {
            // 超过JS安全整数(2^53-1, 约16位)的大整数ID转为字符串，避免精度丢失(如 7675076316887220262 -> 7675076316887220000)
            if (typeof data === 'string') {
                data = data.replace(/(:\s*)(\d{16,})/g, '$1"$2"')
            }
            try {
                return JSON.parse(data)
            } catch {
                return data
            }
        },
    ],
    interceptorRequest: (request) => {
        const userStore = useUserStore()
        const isLogin = userStore.isLogin
        const token = userStore.token

        if (isLogin) {
            request.headers['Authorization'] = token
        }
    },
    interceptorRequestCatch: () => {},
    interceptorResponse: (response) => {
        // 错误处理
        const { success, msg = 'Network Error' } = response.data || {}
        if (![true].includes(success)) {
            message.error({
                content: msg,
                key: MSG_ERROR_KEY,
            })
        }
    },
    interceptorResponseCatch: (err) => {
        const { success, error } = err.response.data || {}
        if ([false].includes(success)) {
            if (error.code === 401) {
                return useUserStore().logout()
            }
            message.error({
                content: error.detail,
                key: MSG_ERROR_KEY,
            })
        }
    },
}

/**
 * 读取文件
 */
class ReadFile extends XYHttp {
    constructor() {
        super({
            baseURL: '',
            responseType: 'blob',
            transformResponse: [
                async (data) => {
                    const encoding = await this._encoding(data)
                    return new Promise((resolve) => {
                        let reader = new FileReader()
                        reader.readAsText(data, encoding)
                        reader.onload = function () {
                            resolve(reader.result)
                        }
                    })
                },
            ],
        })
    }

    /**
     * 文本编码
     * @param data
     * @returns {Promise<unknown>}
     * @private
     */
    _encoding(data) {
        return new Promise((resolve) => {
            let reader = new FileReader()
            reader.readAsBinaryString(data)
            reader.onload = function () {
                resolve(jschardet.detect(reader?.result).encoding)
            }
        })
    }
}

const basic = new XYHttp({
    ...options,
    baseURL: config('http.apiBasic'),
})

const readFile = new ReadFile()

export default {
    basic,
    readFile,
}
