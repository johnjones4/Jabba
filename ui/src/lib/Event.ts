import { ItemsPerPage } from "./consts"
import { stringify } from 'querystring'

export interface iJabbaEvent {
  id: number
  eventVendorType: string
  eventVendorID: string
  created: string
  vendorInfo: any
  isNormal: boolean
}

export interface iJabbaEvents {
  items: Array<iJabbaEvent>
}

export default class JabbaEvent {
  public readonly id: number
  public readonly eventVendorType: string
  public readonly eventVendorID: string
  public readonly created: Date
  public readonly vendorInfo: any
  public readonly isNormal: boolean

  constructor(info: iJabbaEvent) {
    this.id = info.id
    this.eventVendorType = info.eventVendorType
    this.eventVendorID = info.eventVendorID
    this.created = new Date(Date.parse(info.created))
    this.vendorInfo = info.vendorInfo
    this.isNormal = info.isNormal
  }

  static async load(id: number): Promise<JabbaEvent> {
    const url = `/api/event/${id}`
    const response = await fetch(url)
    return new JabbaEvent(await response.json() as iJabbaEvent)
  }

  static async loadEvents(type: string | null, page: number): Promise<Array<JabbaEvent>> {
    const params : any = {
      limit: ItemsPerPage,
      offset: ItemsPerPage * page
    }
    if (type !== null) {
      params.eventVendorType = type
    }
    const url = '/api/event?' + stringify(params)
    const response = await fetch(url)
    const eventsResponse = await response.json() as iJabbaEvents
    return eventsResponse.items.map(e => new JabbaEvent(e))
  }
}
