import { ref, watch } from 'vue'
import {useDictStore} from '@/store/modules/dict'
import dictApi from '@/api/system/dict'
import { getAppLocale } from '@/locale';

/**
 * 批量获取字典数据（支持响应式 + 异步自动更新 + 缓存 + 无重复请求）
 * @param args 字典类型数组
 */
export function useDict(...args: string[]) {
  const dictStore = useDictStore()
  const dictData = ref<Record<string, any[]>>({})

  function getLocaleDict(dictData: any[]) {
    if (!Array.isArray(dictData) || dictData.length === 0) {
      return []
    }
    const locale = getAppLocale()
    // 筛选当前语言的字典项
    return dictData.filter(dict => dict?.langCode === locale)
  }

  // 初始化：自动加载所有字典
  async function loadDict(type: string) {
    // 已有缓存直接返回
    const cache = dictStore.getDict(type)
    if (cache) {
      dictData.value[type] = getLocaleDict(cache)
      return
    }

    // 无缓存 → 请求 + 存缓存
    try {
      const data = await dictApi.getByDictType(type)
      const list = (data) as any  || []
      dictData.value[type] = getLocaleDict(list)
      dictStore.setDict(type, list as any[])
    } catch (err) {
      console.error(`字典【${type}】加载失败：`, err)
      dictData.value[type] = []
    }
  }

  // 并行加载所有字典
  async function loadAll() {
    await Promise.all(args.map(type => loadDict(type)))
  }

  // 自动加载
  loadAll()

  // 监听仓库变化（外部刷新后自动更新）
  watch(
    () => args.map(type => dictStore.getDict(type)),
    () => {
      args.forEach(type => {
        const d = dictStore.getDict(type)
        if (d) dictData.value[type] = getLocaleDict(d)
      })
    },
    { deep: true, immediate: true }
  )

  // 返回响应式字典对象
  return dictData
}