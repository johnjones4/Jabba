import { useParams } from "react-router-dom"
import Event from "./Event"

const EventWrapper = () => {
  const { id } = useParams()
  const idInt = parseInt(id as string)
  return (<Event id={idInt} /> )
}

export default EventWrapper
