import JabbaEvent, { iJabbaEvent } from "./Event"

interface iVendorType {
  eventVendorType: string
  eventVendorName: string
  status: string
  lastEvent: iJabbaEvent
}

interface iVendorTypes {
  eventVendorTypes: Array<string>
}

export default class VendorType {
  public readonly eventVendorType: string
  public readonly eventVendorName: string
  public readonly status: string
  public readonly lastEvent: JabbaEvent

  constructor(info: iVendorType) {
    this.eventVendorType = info.eventVendorType
    this.eventVendorName = info.eventVendorName
    this.status = info.status
    this.lastEvent = new JabbaEvent(info.lastEvent)
  }

  static async load(eventVendorType: string): Promise<VendorType> {
    const response = await fetch(`/api/event-vendor-type/${eventVendorType}`)
    const info = await response.json() as iVendorType
    return new VendorType(info)
  }

  static async loadVendorTypes(): Promise<Array<string>> {
    const response = await fetch('/api/event-vendor-type')
    const infos = await response.json() as iVendorTypes
    return infos.eventVendorTypes
  }

  static async loadAll(): Promise<Array<VendorType>> {
    const types = await VendorType.loadVendorTypes()
    const vendorTypes = new Array<VendorType>()
    for (let type of types) {
      vendorTypes.push(await VendorType.load(type))
    }
    return vendorTypes
  }
}
