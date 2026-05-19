export function buildBasicAuth(clientId: string, clientSecret: string): string {
  return `Basic ${btoa(`${clientId}:${clientSecret}`)}`
}

export function extractCallbackParams(url: string): { code: string | null; state: string | null } {
  try {
    const parsed = new URL(url)
    return {
      code: parsed.searchParams.get('code'),
      state: parsed.searchParams.get('state'),
    }
  } catch {
    const search = url.includes('?') ? url.slice(url.indexOf('?')) : ''
    const params = new URLSearchParams(search)
    return { code: params.get('code'), state: params.get('state') }
  }
}
