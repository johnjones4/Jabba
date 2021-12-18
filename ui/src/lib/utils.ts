import JabbaEvent from "./Event"

export const eventInfoString = (event?: JabbaEvent) => {
  if (event?.vendorInfo['log']) {
    return event?.vendorInfo['log']
  }
  if (event?.vendorInfo['statusCode'] && event?.vendorInfo['body']) {
    return [event?.vendorInfo['statusCode'], event?.vendorInfo['body']].join('\n\n')
  }
  return JSON.stringify(event?.vendorInfo, null, '  ')
}