/**
 * 参数处理
 * @param params 参数
 */
export function tansParams(params: any): string {
  let result = '';
  for (const propName of Object.keys(params)) {
    const value = params[propName];
    const part = encodeURIComponent(propName) + '=';
    if (value !== null && value !== '' && typeof value !== 'undefined') {
      if (typeof value === 'object') {
        for (const key of Object.keys(value)) {
          if (value[key] !== null && value[key] !== '' && typeof value[key] !== 'undefined') {
            const temp = propName + '[' + key + ']';
            const subPart = encodeURIComponent(temp) + '=';
            result += subPart + encodeURIComponent(value[key]) + '&';
          }
        }
      } else {
        result += part + encodeURIComponent(value) + '&';
      }
    }
  }
  return result;
}

/**
 * 验证是否为 blob 格式
 * @param data
 */
export function blobValidate(data: any): boolean {
  return data.type !== 'application/json';
}
