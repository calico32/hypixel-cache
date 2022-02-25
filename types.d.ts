import { Player } from 'hypixel-types'

export type HypixelCacheResponse = SuccessPlayerFound | SuccessNotFound | Failure

export interface SuccessPlayerFound {
  success: true
  player: Player
  fetchedAt: string
  cached: boolean
  username: string
  uuid: string
}

export interface SuccessNotFound {
  success: true
  player: undefined
  fetchedAt: string
  cached: boolean
}

export interface Failure {
  success: false
  error: string
}
