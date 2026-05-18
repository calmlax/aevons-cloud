import { defineStore } from 'pinia'
import { ref } from 'vue'

/**
 * 全局字典缓存 Store
 * 配合 useDict() 使用
 */
export const useDictStore = defineStore('dict', () => {
  // 字典缓存：key = dictType, value = 字典数组
  const dictCache = ref<Record<string, any[]>>({})

  /**
   * 设置字典缓存
   */
  function setDict(dictType: string, data: any[]) {
    dictCache.value[dictType] = data
  }

  /**
   * 获取字典缓存
   */
  function getDict(dictType: string): any[] | null {
    return dictCache.value[dictType] || null
  }

  /**
   * 清空指定字典
   */
  function clearDict(dictType: string) {
    delete dictCache.value[dictType]
  }

  /**
   * 清空所有字典（刷新缓存用）
   */
  function clearAllDict() {
    dictCache.value = {}
  }

  return {
    dictCache,
    setDict,
    getDict,
    clearDict,
    clearAllDict
  }
})