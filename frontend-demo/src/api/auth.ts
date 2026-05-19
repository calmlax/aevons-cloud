import http from './http'
import type { TokenPair } from '../types/auth'
import { buildBasicAuth } from '../utils/auth'

/**
 * Exchange an authorization code for a TokenPair.
 * Uses Basic Auth (clientId:clientSecret) as required by the backend.
 */
export async function exchangeCode(
  clientId: string,
  clientSecret: string,
  code: string,
): Promise<TokenPair> {
  const response = await http.post<TokenPair>(
    '/auth/login',
    { grant_type: 'authorization_code', code },
    { headers: { Authorization: buildBasicAuth(clientId, clientSecret) } },
  )
  return response.data
}
