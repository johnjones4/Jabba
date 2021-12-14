import { useSearchParams } from "react-router-dom"
import Events from "./Events"

const EventsWrapper = () => {
  const [params] = useSearchParams()
  const page = parseInt(params.get('page') || '0')
  return (<Events eventVendorType={params.get('eventVendorType')} page={isNaN(page) ? 0 : page} /> )
}

export default EventsWrapper
