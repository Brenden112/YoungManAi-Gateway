import assert from 'node:assert/strict'
import { describe, test } from 'node:test'
import { CHANNEL_STATUS } from '../constants'
import {
  filterVisibleChannelsForRole,
  isDisabledChannel,
} from './channel-visibility'

const channels = [
  {
    id: 1,
    name: 'official-fixture',
    provider_type: 'official_cloud',
    status: CHANNEL_STATUS.ENABLED,
  },
  {
    id: 2,
    name: 'experimental-enabled-fixture',
    provider_type: 'experimental_proxy',
    status: CHANNEL_STATUS.ENABLED,
  },
  {
    id: 3,
    name: 'experimental-disabled-fixture',
    provider_type: 'experimental_proxy',
    status: CHANNEL_STATUS.MANUAL_DISABLED,
  },
]

describe('channel experimental_proxy visibility', () => {
  test('normal users cannot see experimental_proxy channels', () => {
    const visible = filterVisibleChannelsForRole(channels, 1)

    assert.deepEqual(
      visible.map((channel) => channel.name),
      ['official-fixture']
    )
  })

  test('admins can see experimental_proxy channels', () => {
    const visible = filterVisibleChannelsForRole(channels, 10)

    assert.deepEqual(
      visible.map((channel) => channel.name),
      [
        'official-fixture',
        'experimental-enabled-fixture',
        'experimental-disabled-fixture',
      ]
    )
  })

  test('disabled experimental_proxy channels remain marked disabled', () => {
    const visible = filterVisibleChannelsForRole(channels, 10)
    const disabledExperimental = visible.find(
      (channel) => channel.name === 'experimental-disabled-fixture'
    )

    assert.equal(Boolean(disabledExperimental), true)
    assert.equal(isDisabledChannel(disabledExperimental!), true)
  })
})
